package port

import "context"

// LaunchAtLogin controls OS login registration. The OS is the source of truth.
type LaunchAtLogin interface {
	IsEnabled(ctx context.Context) (bool, error)
	SetEnabled(ctx context.Context, enabled bool) error
}
