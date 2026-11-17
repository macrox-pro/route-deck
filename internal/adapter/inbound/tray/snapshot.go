package tray

import (
	"context"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// Snapshotter adapts a SnapshotReader-like Snapshot call for tray status.
type Snapshotter struct {
	Snapshot func(ctx context.Context) (service.SnapshotResult, error)
}

// ProtectionAndBridge returns protection and bridge words from a snapshot.
func (s Snapshotter) ProtectionAndBridge(ctx context.Context) (domain.ProtectionState, domain.BridgeState, error) {
	if s.Snapshot == nil {
		return domain.ProtectionStateUnknown, domain.BridgeStateUnknown, nil
	}
	result, err := s.Snapshot(ctx)
	if err != nil {
		return domain.ProtectionStateUnknown, domain.BridgeStateUnknown, err
	}
	return result.Protection, result.Bridge, nil
}
