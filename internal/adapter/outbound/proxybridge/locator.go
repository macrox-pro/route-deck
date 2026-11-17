// Package proxybridge implements ProxyBridge CLI discovery, process control, and profile serialization.
package proxybridge

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// CommandRunner runs a short command and returns combined output.
type CommandRunner func(ctx context.Context, name string, args ...string) ([]byte, error)

// Locator discovers ProxyBridge and probes --version.
type Locator struct {
	// Paths overrides platform documented locations when non-nil (including empty).
	Paths  []string
	Runner CommandRunner
	Logger *slog.Logger
	Stat   func(string) (os.FileInfo, error)
}

var _ port.BridgeLocator = (*Locator)(nil)

func (l *Locator) logger() *slog.Logger {
	if l != nil && l.Logger != nil {
		return l.Logger
	}
	return slog.Default()
}

func (l *Locator) runner() CommandRunner {
	if l != nil && l.Runner != nil {
		return l.Runner
	}
	return defaultRunner
}

func (l *Locator) stat() func(string) (os.FileInfo, error) {
	if l != nil && l.Stat != nil {
		return l.Stat
	}
	return os.Stat
}

// PlatformPaths returns documented install locations for this OS.
func (l *Locator) PlatformPaths() []string {
	if l != nil && l.Paths != nil {
		return append([]string(nil), l.Paths...)
	}
	return DefaultDocumentedPaths()
}

func defaultRunner(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.Bytes(), err
}

// Locate checks customPath first, then platform documented paths.
func (l *Locator) Locate(ctx context.Context, customPath string) (port.BridgeDiscovery, error) {
	if err := ctx.Err(); err != nil {
		return port.BridgeDiscovery{}, fmt.Errorf("bridge locate: %w", err)
	}
	platform := l.PlatformPaths()
	candidates := make([]string, 0, 1+len(platform))
	if customPath != "" {
		candidates = append(candidates, filepath.Clean(customPath))
	}
	candidates = append(candidates, platform...)

	stat := l.stat()
	for _, path := range candidates {
		info, err := stat(path)
		if err != nil || info.IsDir() {
			continue
		}
		probeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		out, runErr := l.runner()(probeCtx, path, "--version")
		cancel()
		version, known, warning := ParseVersion(string(out))
		disc := port.BridgeDiscovery{
			State:                domain.BridgeStateDetected,
			Path:                 path,
			Version:              version,
			VersionKnown:         known,
			CompatibilityWarning: warning,
		}
		if runErr != nil && !known {
			l.logger().Info("bridge version probe failed", "component", "proxybridge", "path", path, "err", runErr)
		} else {
			l.logger().Info("bridge discovered",
				"component", "proxybridge",
				"path", path,
				"version", version,
				"bridge_state", string(disc.State),
			)
		}
		return disc, nil
	}
	l.logger().Info("bridge not installed", "component", "proxybridge", "bridge_state", string(domain.BridgeStateNotInstalled))
	return port.BridgeDiscovery{State: domain.BridgeStateNotInstalled}, nil
}
