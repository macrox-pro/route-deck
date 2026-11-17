// Package dock keeps the macOS Dock badge in sync with protection snapshot facts (ADR-0014).
package dock

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
	waildock "github.com/wailsapp/wails/v3/pkg/services/dock"
)

// BadgeSetter is the Wails dock badge surface used by Attacher.
type BadgeSetter interface {
	SetBadge(label string) error
	RemoveBadge() error
}

// Snapshotter loads protection facts for badge updates.
type Snapshotter interface {
	Snapshot(ctx context.Context) (service.SnapshotResult, error)
}

// Attacher keeps the Dock badge in sync with protection snapshot facts.
type Attacher struct {
	Dock     BadgeSetter
	Snapshot Snapshotter
	Logger   *slog.Logger

	mu          sync.Mutex
	unsupported atomic.Bool
	loggedOnce  atomic.Bool
	lastBadge   string
}

func (a *Attacher) logger() *slog.Logger {
	if a != nil && a.Logger != nil {
		return a.Logger
	}
	return slog.Default()
}

// Refresh recomputes and applies the Dock badge from the current snapshot.
func (a *Attacher) Refresh(ctx context.Context) {
	if a == nil || a.Dock == nil {
		return
	}
	if a.unsupported.Load() {
		return
	}
	facts := SnapshotFacts{}
	if a.Snapshot != nil {
		snap, err := a.Snapshot.Snapshot(ctx)
		if err != nil {
			a.logger().Debug("dock badge snapshot failed", "component", "dock", "err", err)
			return
		}
		facts.AppliedConfirmed = snap.AppliedConfirmed
		facts.Protection = snap.Protection
	}
	a.apply(BadgeFor(facts))
}

// ApplyState updates the badge from known protection words without a snapshot round-trip.
func (a *Attacher) ApplyState(appliedConfirmed bool, protection domain.ProtectionState) {
	if a == nil || a.Dock == nil {
		return
	}
	if a.unsupported.Load() {
		return
	}
	a.apply(BadgeFor(SnapshotFacts{
		AppliedConfirmed: appliedConfirmed,
		Protection:       protection,
	}))
}

func (a *Attacher) apply(badge string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if badge == a.lastBadge {
		return
	}
	var err error
	if badge == "" {
		err = a.Dock.RemoveBadge()
	} else {
		err = a.Dock.SetBadge(badge)
	}
	if err != nil {
		a.unsupported.Store(true)
		if a.loggedOnce.CompareAndSwap(false, true) {
			a.logger().Debug("dock badge unsupported", "component", "dock", "err", err)
		}
		return
	}
	a.lastBadge = badge
}

// NewWailsService constructs the official Wails dock service for registration.
func NewWailsService() *waildock.DockService {
	return waildock.New()
}
