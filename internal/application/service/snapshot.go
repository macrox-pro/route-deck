package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProxyHealthRow is one catalog proxy with its health for the snapshot DTO layer.
type ProxyHealthRow struct {
	ID     domain.ProxyID
	Name   string
	Kind   domain.ProxyKind
	Host   string
	Port   int
	Health domain.ProxyHealth
}

// SnapshotResult is the application-layer snapshot before Wails DTO mapping.
type SnapshotResult struct {
	Bridge                    domain.BridgeState
	Protection                domain.ProtectionState
	VersionKnown              bool
	Version                   string
	BridgePath                string
	CompatibilityWarning      string
	DesiredSet                bool
	DesiredID                 domain.ProfileID
	AppliedID                 domain.ProfileID
	AppliedConfirmed          bool
	DesiredDiffersFromApplied bool
	RecheckPending            bool
	DatabaseError             string
	Proxies                   []ProxyHealthRow
}

// Snapshotter builds honest protection snapshots.
type Snapshotter struct {
	Proxies       port.ProxyStore
	Profiles      port.ProfileStore
	Processes     port.ProcessIdentityStore
	Settings      port.SettingsStore
	Locator       port.BridgeLocator
	DatabaseError string
	Logger        *slog.Logger
	// Optional overrides used by the coordinator once wired.
	BridgeOverride            *domain.BridgeState
	VersionKnownOverride      *bool
	DesiredOverride           *port.Settings
	AppliedID                 domain.ProfileID
	AppliedConfirmed          bool
	DesiredDiffersFromApplied bool
	RecheckPending            bool
	HealthByProxy             map[domain.ProxyID]domain.ProxyHealth
}

func (s *Snapshotter) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// Snapshot derives ProtectionState from available facts.
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
		s.logger().Debug("snapshot database damaged",
			"component", "snapshotter",
			"protection_state", string(result.Protection),
		)
		return result, nil
	}

	customPath := ""
	if s.DesiredOverride != nil {
		result.DesiredSet = s.DesiredOverride.DesiredSet
		result.DesiredID = s.DesiredOverride.DesiredProfileID
		facts.DesiredSet = s.DesiredOverride.DesiredSet
		facts.Desired = s.DesiredOverride.DesiredProfileID
		customPath = s.DesiredOverride.ProxyBridgePath
	} else if s.Settings != nil {
		settings, err := s.Settings.Get(ctx)
		if err != nil {
			return SnapshotResult{}, fmt.Errorf("snapshot settings: %w", err)
		}
		result.DesiredSet = settings.DesiredSet
		result.DesiredID = settings.DesiredProfileID
		facts.DesiredSet = settings.DesiredSet
		facts.Desired = settings.DesiredProfileID
		customPath = settings.ProxyBridgePath
	}

	if s.Locator != nil {
		disc, err := s.Locator.Locate(ctx, customPath)
		if err != nil {
			return SnapshotResult{}, fmt.Errorf("snapshot locate: %w", err)
		}
		result.Bridge = disc.State
		result.VersionKnown = disc.VersionKnown
		result.Version = disc.Version
		result.BridgePath = disc.Path
		result.CompatibilityWarning = disc.CompatibilityWarning
		facts.Bridge = disc.State
		facts.VersionKnown = disc.VersionKnown
	}

	if s.BridgeOverride != nil {
		result.Bridge = *s.BridgeOverride
		facts.Bridge = *s.BridgeOverride
	}
	if s.VersionKnownOverride != nil {
		result.VersionKnown = *s.VersionKnownOverride
		facts.VersionKnown = *s.VersionKnownOverride
	}

	result.AppliedID = s.AppliedID
	result.AppliedConfirmed = s.AppliedConfirmed
	result.DesiredDiffersFromApplied = s.DesiredDiffersFromApplied
	if result.DesiredSet && result.AppliedConfirmed && result.DesiredID != "" && result.DesiredID != result.AppliedID {
		result.DesiredDiffersFromApplied = true
	}
	result.RecheckPending = s.RecheckPending
	facts.Applied = s.AppliedID
	facts.AppliedConfirmed = s.AppliedConfirmed
	facts.RecheckPending = s.RecheckPending

	proxies, err := s.Proxies.List(ctx)
	if err != nil {
		return SnapshotResult{}, fmt.Errorf("snapshot list proxies: %w", err)
	}
	rows := make([]ProxyHealthRow, 0, len(proxies))
	proxyFacts := make([]domain.ProxyFact, 0, len(proxies))
	for _, proxy := range proxies {
		health := domain.ProxyHealthUnknown
		if s.HealthByProxy != nil {
			if h, ok := s.HealthByProxy[proxy.ID]; ok {
				health = h
			}
		}
		rows = append(rows, ProxyHealthRow{
			ID:     proxy.ID,
			Name:   proxy.Name,
			Kind:   proxy.Kind,
			Host:   proxy.Host,
			Port:   proxy.Port,
			Health: health,
		})
		proxyFacts = append(proxyFacts, domain.ProxyFact{ID: proxy.ID, Health: health})
	}
	result.Proxies = rows
	facts.Proxies = proxyFacts

	if s.AppliedConfirmed && s.AppliedID != "" {
		if s.Profiles == nil || s.Processes == nil {
			return SnapshotResult{}, fmt.Errorf("snapshot: profiles or processes store not configured")
		}
		profile, err := s.Profiles.Get(ctx, s.AppliedID)
		if err != nil {
			return SnapshotResult{}, fmt.Errorf("snapshot applied profile: %w", err)
		}
		identities, err := s.Processes.List(ctx)
		if err != nil {
			return SnapshotResult{}, fmt.Errorf("snapshot process identities: %w", err)
		}
		plan, err := domain.BuildRulePlan(profile, proxies, identities)
		if err != nil {
			return SnapshotResult{}, fmt.Errorf("snapshot rule plan: %w", err)
		}
		facts.Apps = plan.Apps
	}

	result.Protection = domain.DeriveProtection(facts)
	s.logger().Debug("snapshot built",
		"component", "snapshotter",
		"protection_state", string(result.Protection),
		"bridge_state", string(result.Bridge),
	)
	return result, nil
}
