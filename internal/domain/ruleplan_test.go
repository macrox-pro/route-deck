package domain_test

import (
	"errors"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestBuildRulePlan(t *testing.T) {
	socks, err := domain.NewProxyEndpoint("p1", "socks", domain.ProxyKindSOCKS5, "127.0.0.1", 1080, "kc")
	if err != nil {
		t.Fatal(err)
	}
	httpProxy, err := domain.NewProxyEndpoint("p2", "http", domain.ProxyKindHTTP, "127.0.0.1", 8080, "kc")
	if err != nil {
		t.Fatal(err)
	}
	identity := domain.ProcessIdentity{
		ID: "proc-1", ApplicationID: "app-1", BundleID: "com.example.app", ExecutableName: "App",
	}
	proxyAction, err := domain.ProxyAction("p1")
	if err != nil {
		t.Fatal(err)
	}
	httpAction, err := domain.ProxyAction("p2")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		policy    domain.ApplicationPolicy
		proxies   []domain.ProxyEndpoint
		wantRules int
		wantTCP   bool
		wantUDP   bool
		wantErr   error
	}{
		{
			name:      "socks_both",
			policy:    mustPolicy(t, "proc-1", proxyAction, domain.ProtocolScopeBOTH, ""),
			proxies:   []domain.ProxyEndpoint{socks},
			wantRules: 1, wantTCP: true, wantUDP: true,
		},
		{
			name:      "http_both_block_udp",
			policy:    mustPolicy(t, "proc-1", httpAction, domain.ProtocolScopeBOTH, domain.UDPLimitationAbsent),
			proxies:   []domain.ProxyEndpoint{httpProxy},
			wantRules: 2, wantTCP: true, wantUDP: false,
		},
		{
			name:      "http_both_direct_udp",
			policy:    mustPolicy(t, "proc-1", httpAction, domain.ProtocolScopeBOTH, domain.UDPLimitationDirectUDP),
			proxies:   []domain.ProxyEndpoint{httpProxy},
			wantRules: 2, wantTCP: true, wantUDP: false,
		},
		{
			name:      "tcp_only_blocks_udp",
			policy:    mustPolicy(t, "proc-1", proxyAction, domain.ProtocolScopeTCP, ""),
			proxies:   []domain.ProxyEndpoint{socks},
			wantRules: 2, wantTCP: true, wantUDP: true,
		},
		{
			name:    "unresolved",
			policy:  mustPolicy(t, "proc-1", proxyAction, domain.ProtocolScopeBOTH, ""),
			proxies: nil,
			wantErr: domain.ErrUnresolvedProxy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := domain.BuildRulePlan(
				domain.Profile{ID: "prof", Name: "Work", Policies: []domain.ApplicationPolicy{tt.policy}},
				tt.proxies,
				[]domain.ProcessIdentity{identity},
			)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("error = %v", err)
			}
			if len(plan.Rules) != tt.wantRules {
				t.Fatalf("rules = %d, want %d", len(plan.Rules), tt.wantRules)
			}
			if len(plan.Apps) != 1 {
				t.Fatalf("apps = %d", len(plan.Apps))
			}
			if plan.Apps[0].TCPCovered != tt.wantTCP || plan.Apps[0].UDPCovered != tt.wantUDP {
				t.Fatalf("coverage tcp=%v udp=%v, want %v %v", plan.Apps[0].TCPCovered, plan.Apps[0].UDPCovered, tt.wantTCP, tt.wantUDP)
			}
		})
	}
}

func mustPolicy(t *testing.T, proc string, action domain.RoutingAction, scope domain.ProtocolScope, lim domain.UDPLimitation) domain.ApplicationPolicy {
	t.Helper()
	p, err := domain.NewApplicationPolicy(domain.ProcessIdentityID(proc), action, scope, lim)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestBuildRulePlanCollapsesSharedBundleID(t *testing.T) {
	identities := []domain.ProcessIdentity{
		{ID: "h1", ApplicationID: "app", BundleID: "com.openai.codex.helper", ExecutableName: "Codex (GPU)"},
		{ID: "h2", ApplicationID: "app", BundleID: "com.openai.codex.helper", ExecutableName: "Codex (Service)"},
		{ID: "h3", ApplicationID: "app", BundleID: "com.openai.codex.helper", ExecutableName: "Codex (Aperitif)"},
		{ID: "h4", ApplicationID: "app", BundleID: "com.openai.codex.helper", ExecutableName: "Codex (Aperitif GPU)"},
	}
	policies := make([]domain.ApplicationPolicy, 0, 4)
	for _, id := range []string{"h1", "h2", "h3", "h4"} {
		policies = append(policies, mustPolicy(t, id, domain.BlockAction(), domain.ProtocolScopeBOTH, ""))
	}
	plan, err := domain.BuildRulePlan(
		domain.Profile{ID: "p", Name: "Work", Policies: policies},
		nil,
		identities,
	)
	if err != nil {
		t.Fatalf("BuildRulePlan: %v", err)
	}
	if len(plan.Rules) != 1 {
		t.Fatalf("rules = %d, want 1; %+v", len(plan.Rules), plan.Rules)
	}
	if len(plan.Rules[0].ProcessNames) != 1 || plan.Rules[0].ProcessNames[0] != "com.openai.codex.helper" {
		t.Fatalf("processNames = %+v, want [com.openai.codex.helper]", plan.Rules[0].ProcessNames)
	}
}

func TestBuildRulePlanConflictingMatch(t *testing.T) {
	identities := []domain.ProcessIdentity{
		{ID: "h1", ApplicationID: "app", BundleID: "com.openai.codex.helper", ExecutableName: "Codex (GPU)"},
		{ID: "h2", ApplicationID: "app", BundleID: "com.openai.codex.helper", ExecutableName: "Codex (Service)"},
	}
	_, err := domain.BuildRulePlan(
		domain.Profile{ID: "p", Name: "Work", Policies: []domain.ApplicationPolicy{
			mustPolicy(t, "h1", domain.BlockAction(), domain.ProtocolScopeBOTH, ""),
			mustPolicy(t, "h2", domain.DirectAction(), domain.ProtocolScopeBOTH, ""),
		}},
		nil,
		identities,
	)
	if !errors.Is(err, domain.ErrConflictingMatch) {
		t.Fatalf("error = %v, want ErrConflictingMatch", err)
	}
}
