package wails

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// LiveSnapshotter syncs coordinator and health facts into Snapshotter before each read.
type LiveSnapshotter struct {
	Base           *service.Snapshotter
	Protection     *service.Protection
	Health         *service.HealthMonitor
	RecheckPending *atomic.Bool
}

// Snapshot implements SnapshotReader.
func (l *LiveSnapshotter) Snapshot(ctx context.Context) (service.SnapshotResult, error) {
	if l == nil || l.Base == nil {
		return service.SnapshotResult{}, fmt.Errorf("live snapshot: not configured")
	}
	if l.Protection != nil {
		appliedID, confirmed := l.Protection.Applied()
		l.Base.AppliedID = appliedID
		l.Base.AppliedConfirmed = confirmed
		bridge := l.Protection.BridgeState()
		if bridge != domain.BridgeStateUnknown {
			l.Base.BridgeOverride = &bridge
		}
		l.Base.DesiredDiffersFromApplied = l.Protection.DesiredDiffers()
	}
	if l.Health != nil {
		l.Base.HealthByProxy = l.Health.Snapshot()
	}
	if l.RecheckPending != nil {
		l.Base.RecheckPending = l.RecheckPending.Load()
	}
	return l.Base.Snapshot(ctx)
}

// ClearRecheckPending clears the wake/network pending flag after a successful recheck.
func (l *LiveSnapshotter) ClearRecheckPending() {
	if l != nil && l.RecheckPending != nil {
		l.RecheckPending.Store(false)
	}
}

// MarkRecheckPending sets Unknown until health/bridge recheck finishes.
func (l *LiveSnapshotter) MarkRecheckPending() {
	if l != nil && l.RecheckPending != nil {
		l.RecheckPending.Store(true)
	}
}
