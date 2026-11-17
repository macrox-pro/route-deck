package service

import (
	"context"
	"fmt"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProxyHealthRow is one catalog proxy with its health for the snapshot DTO layer.
type ProxyHealthRow struct {
	ID     domain.ProxyID
	Name   string
	Health domain.ProxyHealth
}

// SnapshotResult is the application-layer snapshot before Wails DTO mapping.
type SnapshotResult struct {
	Bridge           domain.BridgeState
	Protection       domain.ProtectionState
	VersionKnown     bool
	DesiredSet       bool
	DesiredID        domain.ProfileID
	AppliedID        domain.ProfileID
	AppliedConfirmed bool
	DatabaseError    string
	Proxies          []ProxyHealthRow
}

// Snapshotter builds honest protection snapshots without a ProxyBridge CLI.
type Snapshotter struct {
	Proxies port.ProxyStore
	// DatabaseError is set when the DB failed to open at startup.
	DatabaseError string
}

// Snapshot derives ProtectionState. Without a CLI, VersionKnown is false so the word is Unknown.
func (s *Snapshotter) Snapshot(ctx context.Context) (SnapshotResult, error) {
	if err := ctx.Err(); err != nil {
		return SnapshotResult{}, fmt.Errorf("snapshot: %w", err)
	}

	result := SnapshotResult{
		Bridge:       domain.BridgeStateUnknown,
		VersionKnown: false,
		DesiredSet:   false,
	}

	facts := domain.ProtectionFacts{
		Bridge:       domain.BridgeStateUnknown,
		VersionKnown: false,
		DesiredSet:   false,
	}

	if s.DatabaseError != "" {
		facts.DatabaseDamaged = true
		result.DatabaseError = s.DatabaseError
		result.Protection = domain.DeriveProtection(facts)
		return result, nil
	}

	proxies, err := s.Proxies.List(ctx)
	if err != nil {
		return SnapshotResult{}, fmt.Errorf("snapshot list proxies: %w", err)
	}
	rows := make([]ProxyHealthRow, 0, len(proxies))
	for _, proxy := range proxies {
		rows = append(rows, ProxyHealthRow{
			ID:     proxy.ID,
			Name:   proxy.Name,
			Health: domain.ProxyHealthUnknown,
		})
	}
	result.Proxies = rows
	result.Protection = domain.DeriveProtection(facts)
	return result, nil
}
