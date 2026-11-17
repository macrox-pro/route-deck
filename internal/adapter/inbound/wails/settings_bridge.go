package wails

import (
	"context"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// PortSettingsReader adapts port.SettingsStore for ProtectionService.
type PortSettingsReader struct {
	Store port.SettingsStore
}

// Get returns the desired profile id and whether it is set.
func (r PortSettingsReader) Get(ctx context.Context) (domain.ProfileID, bool, error) {
	settings, err := r.Store.Get(ctx)
	if err != nil {
		return "", false, err
	}
	return settings.DesiredProfileID, settings.DesiredSet, nil
}
