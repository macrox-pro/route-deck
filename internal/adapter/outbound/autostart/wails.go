// Package autostart adapts Wails AutostartManager to port.LaunchAtLogin (ADR-0010).
package autostart

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Manager wraps application.AutostartManager.
type Manager struct {
	Autostart *application.AutostartManager
}

// IsEnabled reports whether the OS has a login registration.
func (m *Manager) IsEnabled(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if m == nil || m.Autostart == nil {
		return false, fmt.Errorf("autostart: not configured")
	}
	return m.Autostart.IsEnabled()
}

// SetEnabled registers or removes Launch at Login.
func (m *Manager) SetEnabled(ctx context.Context, enabled bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m == nil || m.Autostart == nil {
		return fmt.Errorf("autostart: not configured")
	}
	if enabled {
		if err := m.Autostart.Enable(); err != nil {
			return fmt.Errorf("autostart enable: %w", err)
		}
		return nil
	}
	if err := m.Autostart.Disable(); err != nil {
		return fmt.Errorf("autostart disable: %w", err)
	}
	return nil
}
