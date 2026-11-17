package domain_test

import (
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func quietFacts() domain.ProtectionFacts {
	return domain.ProtectionFacts{
		Bridge:           domain.BridgeStateReady,
		VersionKnown:     true,
		Desired:          "desired",
		DesiredSet:       true,
		Applied:          "applied",
		AppliedConfirmed: true,
		ConfirmFailed:    false,
		RollbackFailed:   false,
		RecheckPending:   false,
		DatabaseDamaged:  false,
		DomainRules:      false,
		DNSIncomplete:    false,
		Proxies:          nil,
		Apps:             nil,
	}
}

func TestDeriveProtection(t *testing.T) {
	proxyAction, err := domain.ProxyAction("proxy-1")
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}

	covered := domain.AppCoverage{
		ID:          "app-1",
		Action:      proxyAction,
		TCPCovered:  true,
		UDPCovered:  true,
		UDPRequired: true,
	}
	open := domain.AppCoverage{
		ID:          "app-2",
		Action:      proxyAction,
		TCPCovered:  false,
		UDPCovered:  false,
		UDPRequired: false,
	}
	available := domain.ProxyFact{
		ID:     "proxy-1",
		Health: domain.ProxyHealthAvailable,
	}
	blockCovered := domain.AppCoverage{
		ID:          "app-1",
		Action:      domain.BlockAction(),
		TCPCovered:  true,
		UDPCovered:  true,
		UDPRequired: true,
	}
	directOnly := domain.AppCoverage{
		ID:          "app-1",
		Action:      domain.DirectAction(),
		TCPCovered:  false,
		UDPCovered:  false,
		UDPRequired: false,
	}

	udpNotRequired := covered
	udpNotRequired.UDPRequired = false
	udpNotRequired.UDPCovered = false

	udpRequiredGap := covered
	udpRequiredGap.UDPCovered = false

	tcpUDPUnavailableApp := covered
	tcpUDPUnavailableApp.UDPCovered = false

	tests := []struct {
		name string
		in   domain.ProtectionFacts
		want domain.ProtectionState
	}{
		{
			name: "step1",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.RecheckPending = true
				f.Apps = []domain.AppCoverage{covered}
				f.Proxies = []domain.ProxyFact{available}
				return f
			}(),
			want: domain.ProtectionStateUnknown,
		},
		{
			name: "step2",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Bridge = domain.BridgeStateUnknown
				return f
			}(),
			want: domain.ProtectionStateUnknown,
		},
		{
			name: "step3",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Bridge = domain.BridgeStateStarting
				return f
			}(),
			want: domain.ProtectionStateInitializing,
		},
		{
			name: "step4",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Bridge = domain.BridgeStateFailed
				f.ConfirmFailed = true
				f.RollbackFailed = true
				return f
			}(),
			want: domain.ProtectionStateUnprotected,
		},
		{
			name: "step5",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Bridge = domain.BridgeStateFailed
				return f
			}(),
			want: domain.ProtectionStateFailed,
		},
		{
			name: "step6",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Bridge = domain.BridgeStateConflict
				return f
			}(),
			want: domain.ProtectionStateUnprotected,
		},
		{
			name: "step7",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Bridge = domain.BridgeStateNotInstalled
				return f
			}(),
			want: domain.ProtectionStateInactive,
		},
		{
			name: "step8",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Bridge = domain.BridgeStateStopped
				return f
			}(),
			want: domain.ProtectionStateUnprotected,
		},
		{
			name: "step9",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Bridge = domain.BridgeStateStopped
				f.DesiredSet = false
				return f
			}(),
			want: domain.ProtectionStateInactive,
		},
		{
			name: "step10",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.AppliedConfirmed = false
				return f
			}(),
			want: domain.ProtectionStateUnprotected,
		},
		{
			name: "step11",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{directOnly}
				return f
			}(),
			want: domain.ProtectionStateInactive,
		},
		{
			name: "step12",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{covered}
				f.Proxies = []domain.ProxyFact{available}
				return f
			}(),
			want: domain.ProtectionStateActive,
		},
		{
			name: "step12_block",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{blockCovered}
				return f
			}(),
			want: domain.ProtectionStateActive,
		},
		{
			name: "step12_udp_not_required",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{udpNotRequired}
				f.Proxies = []domain.ProxyFact{available}
				return f
			}(),
			want: domain.ProtectionStateActive,
		},
		{
			name: "step13",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{covered, open}
				f.Proxies = []domain.ProxyFact{available}
				return f
			}(),
			want: domain.ProtectionStateDegraded,
		},
		{
			name: "step13_udp_required",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{udpRequiredGap}
				f.Proxies = []domain.ProxyFact{available}
				return f
			}(),
			want: domain.ProtectionStateDegraded,
		},
		{
			name: "step13_proxy_unavailable",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{covered}
				f.Proxies = []domain.ProxyFact{{
					ID:     "proxy-1",
					Health: domain.ProxyHealthUnavailable,
				}}
				return f
			}(),
			want: domain.ProtectionStateDegraded,
		},
		{
			name: "step13_auth_failed",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{covered}
				f.Proxies = []domain.ProxyFact{{
					ID:     "proxy-1",
					Health: domain.ProxyHealthAuthenticationFailed,
				}}
				return f
			}(),
			want: domain.ProtectionStateDegraded,
		},
		{
			name: "step13_tcp_available_udp_unavailable",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{tcpUDPUnavailableApp}
				f.Proxies = []domain.ProxyFact{{
					ID:     "proxy-1",
					Health: domain.ProxyHealthTCPAvailableUDPUnavailable,
				}}
				return f
			}(),
			want: domain.ProtectionStateDegraded,
		},
		{
			name: "step13_dns_incomplete",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{covered}
				f.Proxies = []domain.ProxyFact{available}
				f.DomainRules = true
				f.DNSIncomplete = true
				return f
			}(),
			want: domain.ProtectionStateDegraded,
		},
		{
			name: "step14",
			in: func() domain.ProtectionFacts {
				f := quietFacts()
				f.Apps = []domain.AppCoverage{open}
				f.Proxies = []domain.ProxyFact{available}
				return f
			}(),
			want: domain.ProtectionStateUnprotected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.DeriveProtection(tt.in)
			if got != tt.want {
				t.Fatalf("DeriveProtection() = %q, want %q", got, tt.want)
			}
		})
	}
}
