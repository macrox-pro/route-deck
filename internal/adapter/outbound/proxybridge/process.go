package proxybridge

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/pkg/retry"
)

// Process manages one ProxyBridge CLI instance.
type Process struct {
	Logger *slog.Logger
	// LineNotify is optional. Called for non-fatal user-action lines (once per Start for approval).
	LineNotify func(class port.BridgeLineClass)

	mu               sync.Mutex
	cmd              *exec.Cmd
	cancel           context.CancelFunc
	ready            chan struct{}
	failed           chan error // classification / early exit for WaitConfirm
	exited           chan error // single cmd.Wait result; Stop must not call Wait again
	running          bool
	approvalNotified atomic.Bool
}

var _ port.BridgeProcess = (*Process)(nil)

func (p *Process) logger() *slog.Logger {
	if p != nil && p.Logger != nil {
		return p.Logger
	}
	return slog.Default()
}

// Start launches binaryPath with --profile - and profileJSON on stdin.
func (p *Process) Start(ctx context.Context, binaryPath string, profileJSON []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.running {
		return fmt.Errorf("bridge process: already running")
	}
	runCtx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(runCtx, binaryPath, "--profile", "-", "--verbose", "1")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("bridge process stdin: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("bridge process stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("bridge process stderr: %w", err)
	}
	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("bridge process start: %w", err)
	}
	if _, err := stdin.Write(profileJSON); err != nil {
		_ = cmd.Process.Kill()
		cancel()
		return fmt.Errorf("bridge process write profile: %w", err)
	}
	_ = stdin.Close()

	p.cmd = cmd
	p.cancel = cancel
	p.ready = make(chan struct{})
	p.failed = make(chan error, 1)
	p.exited = make(chan error, 1)
	p.running = true
	p.approvalNotified.Store(false)

	go p.readLines(stdout)
	go p.readLines(stderr)
	go func() {
		err := cmd.Wait()
		p.mu.Lock()
		p.running = false
		p.mu.Unlock()
		select {
		case p.exited <- err:
		default:
		}
		select {
		case p.failed <- err:
		default:
		}
	}()
	_ = ctx
	return nil
}

func (p *Process) readLines(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		class := ClassifyLine(line)
		truncated := line
		if len(truncated) > 500 {
			truncated = truncated[:500]
		}
		p.logger().Info("bridge line", "component", "proxybridge", "line_class", string(class), "msg", truncated)
		switch class {
		case port.BridgeLineRedirectReady:
			select {
			case <-p.ready:
			default:
				close(p.ready)
			}
		case port.BridgeLineGUIConflict, port.BridgeLineAnotherCLI, port.BridgeLineRejectedRules, port.BridgeLineExtensionTimeout:
			select {
			case p.failed <- fmt.Errorf("%s", class):
			default:
			}
		case port.BridgeLineExtensionApproval:
			if p.approvalNotified.CompareAndSwap(false, true) && p.LineNotify != nil {
				p.LineNotify(class)
			}
		case port.BridgeLineUnclassified, port.BridgeLineProfileLoaded, port.BridgeLineTunnelStopped:
			// logged only
		}
	}
}

// WaitConfirm waits until redirect-ready or a failure classification.
func (p *Process) WaitConfirm(ctx context.Context) error {
	p.mu.Lock()
	ready := p.ready
	failed := p.failed
	p.mu.Unlock()
	if ready == nil {
		return fmt.Errorf("bridge process: not started")
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("bridge process wait: %w", ctx.Err())
	case err := <-failed:
		if err == nil {
			return fmt.Errorf("bridge process: exited before confirm")
		}
		return fmt.Errorf("bridge process: %w", err)
	case <-ready:
		// still check rejected-rules raced in
		select {
		case err := <-failed:
			if err != nil {
				return fmt.Errorf("bridge process: %w", err)
			}
		default:
		}
		return nil
	}
}

// Stop grace: CLI may take ~5s for extension cleanup; keep headroom beyond that.
const (
	stopGracePeriod  = 20 * time.Second
	stopKillWait     = 5 * time.Second
	stopKillAttempts = 3
	stopKillDelay    = time.Second
)

// Stop signals SIGTERM, waits for exit, then SIGKILL with retries if needed.
// Exit after our signal (including non-zero Wait) is success: the process is gone.
// Only fails when the process cannot be signaled or does not die after kill retries.
// MUST NOT call cmd.Wait here — Start owns the single Wait and publishes on exited.
func (p *Process) Stop(ctx context.Context) error {
	p.mu.Lock()
	cmd := p.cmd
	exited := p.exited
	running := p.running
	p.mu.Unlock()
	if !running || cmd == nil || cmd.Process == nil {
		return nil
	}
	if exited == nil {
		return fmt.Errorf("bridge process stop: missing exit channel")
	}

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return p.reapExited(exited, stopKillWait)
		}
		return fmt.Errorf("bridge process signal: %w", err)
	}

	select {
	case <-ctx.Done():
		p.logger().Warn("bridge stop interrupted; killing process", "component", "proxybridge", "err", ctx.Err())
		return p.killUntilExited(cmd, exited)
	case <-time.After(stopGracePeriod):
		p.logger().Warn("bridge stop grace elapsed; killing process", "component", "proxybridge", "grace", stopGracePeriod.String())
		return p.killUntilExited(cmd, exited)
	case err := <-exited:
		return p.stoppedAfterWait(err)
	}
}

func (p *Process) killUntilExited(cmd *exec.Cmd, exited <-chan error) error {
	// Detached from the caller's cancel: stop must still reap after SIGKILL.
	err := retry.Do(context.Background(), retry.Policy{Attempts: stopKillAttempts, Delay: stopKillDelay}, func(context.Context) error {
		if cmd.Process != nil {
			if killErr := cmd.Process.Kill(); killErr != nil && !errors.Is(killErr, os.ErrProcessDone) {
				p.logger().Warn("bridge kill", "component", "proxybridge", "err", killErr)
			}
		}
		ok, waitErr := recvExited(exited, stopKillWait)
		if !ok {
			return fmt.Errorf("exit wait timed out after %s", stopKillWait)
		}
		if waitErr != nil {
			p.logger().Debug("bridge wait after kill", "component", "proxybridge", "err", waitErr)
		}
		return nil
	})
	if err != nil {
		p.markStopped()
		return fmt.Errorf("bridge process stop: process did not exit after kill: %w", err)
	}
	p.markStopped()
	return nil
}

func (p *Process) reapExited(exited <-chan error, timeout time.Duration) error {
	ok, waitErr := recvExited(exited, timeout)
	if !ok {
		p.markStopped()
		return nil
	}
	return p.stoppedAfterWait(waitErr)
}

func recvExited(exited <-chan error, timeout time.Duration) (ok bool, waitErr error) {
	select {
	case err := <-exited:
		return true, err
	case <-time.After(timeout):
		return false, nil
	}
}

func (p *Process) stoppedAfterWait(err error) error {
	p.markStopped()
	if err != nil {
		// Terminated-by-signal / non-zero exit after SIGTERM is still a successful stop.
		p.logger().Debug("bridge wait after stop", "component", "proxybridge", "err", err)
	}
	return nil
}

func (p *Process) markStopped() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.running = false
	if p.cancel != nil {
		p.cancel()
	}
}

// Running reports whether the process is considered running.
func (p *Process) Running() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}
