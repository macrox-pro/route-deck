package proxybridge

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

// Process manages one ProxyBridge CLI instance.
type Process struct {
	Logger *slog.Logger

	mu      sync.Mutex
	cmd     *exec.Cmd
	cancel  context.CancelFunc
	ready   chan struct{}
	failed  chan error
	running bool
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
	p.running = true

	go p.readLines(stdout)
	go p.readLines(stderr)
	go func() {
		err := cmd.Wait()
		p.mu.Lock()
		p.running = false
		p.mu.Unlock()
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

// Stop signals SIGTERM and waits with timeout.
func (p *Process) Stop(ctx context.Context) error {
	p.mu.Lock()
	cmd := p.cmd
	running := p.running
	p.mu.Unlock()
	if !running || cmd == nil || cmd.Process == nil {
		return nil
	}
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("bridge process signal: %w", err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		return fmt.Errorf("bridge process stop: %w", ctx.Err())
	case <-time.After(6 * time.Second):
		_ = cmd.Process.Kill()
		return fmt.Errorf("bridge process stop: timeout")
	case err := <-done:
		p.mu.Lock()
		p.running = false
		if p.cancel != nil {
			p.cancel()
		}
		p.mu.Unlock()
		if err != nil {
			return fmt.Errorf("bridge process stop: %w", err)
		}
		return nil
	}
}

// Running reports whether the process is considered running.
func (p *Process) Running() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}
