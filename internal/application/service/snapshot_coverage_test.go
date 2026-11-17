package service_test

import (
	"context"
	"testing"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type memProcessStore struct {
	rows []domain.ProcessIdentity
}

func (m *memProcessStore) Get(context.Context, domain.ProcessIdentityID) (domain.ProcessIdentity, error) {
	return domain.ProcessIdentity{}, nil
}

func (m *memProcessStore) List(context.Context) ([]domain.ProcessIdentity, error) {
	return append([]domain.ProcessIdentity(nil), m.rows...), nil
}

func (m *memProcessStore) Save(context.Context, domain.ProcessIdentity) error { return nil }

func (m *memProcessStore) Delete(context.Context, domain.ProcessIdentityID) error { return nil }

func readyTrue() (domain.BridgeState, bool) {
	return domain.BridgeStateReady, true
}

func TestSnapshotProtectionCoverage(t *testing.T) {
	ready, versionKnown := readyTrue()
	proc := domain.ProcessIdentity{
		ID:             "proc-1",
		ApplicationID:  "app-1",
		BundleID:       "com.example.app",
		ExecutableName: "Example",
		Path:           "/Applications/Example.app/Contents/MacOS/Example",
	}
	endpoint, err := domain.NewProxyEndpoint("proxy-1", "lab", domain.ProxyKindSOCKS5, "127.0.0.1", 1080, "kc")
	if err != nil {
		t.Fatalf("NewProxyEndpoint: %v", err)
	}

	blockPolicy, err := domain.NewApplicationPolicy(proc.ID, domain.BlockAction(), domain.ProtocolScopeBOTH, domain.UDPLimitationAbsent)
	if err != nil {
		t.Fatalf("NewApplicationPolicy block: %v", err)
	}
	proxyAction, err := domain.ProxyAction(endpoint.ID)
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}
	proxyPolicy, err := domain.NewApplicationPolicy(proc.ID, proxyAction, domain.ProtocolScopeBOTH, domain.UDPLimitationAbsent)
	if err != nil {
		t.Fatalf("NewApplicationPolicy proxy: %v", err)
	}

	tests := []struct {
		name             string
		policies         []domain.ApplicationPolicy
		appliedConfirmed bool
		health           domain.ProxyHealth
		want             domain.ProtectionState
	}{
		{
			name:             "block_active",
			policies:         []domain.ApplicationPolicy{blockPolicy},
			appliedConfirmed: true,
			want:             domain.ProtectionStateActive,
		},
		{
			name:             "empty_inactive",
			policies:         nil,
			appliedConfirmed: true,
			want:             domain.ProtectionStateInactive,
		},
		{
			name:             "proxy_udp_gap_degraded",
			policies:         []domain.ApplicationPolicy{proxyPolicy},
			appliedConfirmed: true,
			health:           domain.ProxyHealthTCPAvailableUDPUnavailable,
			want:             domain.ProtectionStateDegraded,
		},
		{
			name:             "not_confirmed_unprotected",
			policies:         []domain.ApplicationPolicy{blockPolicy},
			appliedConfirmed: false,
			want:             domain.ProtectionStateUnprotected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := domain.NewProfile("prof-1", "AI", tt.policies)
			if err != nil {
				t.Fatalf("NewProfile: %v", err)
			}
			proxies := &memProxyStore{byID: map[domain.ProxyID]domain.ProxyEndpoint{
				endpoint.ID: endpoint,
			}}
			profiles := &memProfileStore{byID: map[domain.ProfileID]domain.Profile{
				profile.ID: profile,
			}}
			processes := &memProcessStore{rows: []domain.ProcessIdentity{proc}}
			healthByProxy := map[domain.ProxyID]domain.ProxyHealth{}
			if tt.health != "" {
				healthByProxy[endpoint.ID] = tt.health
			} else {
				healthByProxy[endpoint.ID] = domain.ProxyHealthAvailable
			}
			snap := &service.Snapshotter{
				Proxies:              proxies,
				Profiles:             profiles,
				Processes:            processes,
				BridgeOverride:       &ready,
				VersionKnownOverride: &versionKnown,
				AppliedID:            profile.ID,
				AppliedConfirmed:     tt.appliedConfirmed,
				HealthByProxy:        healthByProxy,
			}
			got, err := snap.Snapshot(context.Background())
			if err != nil {
				t.Fatalf("Snapshot: %v", err)
			}
			if got.Protection != tt.want {
				t.Fatalf("Protection = %q, want %q", got.Protection, tt.want)
			}
		})
	}
}
