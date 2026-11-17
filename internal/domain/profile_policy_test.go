package domain_test

import (
	"errors"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestParseProtocolScope(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    domain.ProtocolScope
		wantErr error
	}{
		{name: "tcp", in: "TCP", want: domain.ProtocolScopeTCP},
		{name: "udp", in: "UDP", want: domain.ProtocolScopeUDP},
		{name: "both", in: "BOTH", want: domain.ProtocolScopeBOTH},
		{name: "rejected", in: "tcp", wantErr: domain.ErrUnknownProtocolScope},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ParseProtocolScope(tt.in)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("got (%q, %v), want %q", got, err, tt.want)
			}
		})
	}
}

func TestParseUDPLimitation(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    domain.UDPLimitation
		wantErr error
	}{
		{name: "absent", in: "Absent", want: domain.UDPLimitationAbsent},
		{name: "direct", in: "DirectUDP", want: domain.UDPLimitationDirectUDP},
		{name: "rejected", in: "", wantErr: domain.ErrUnknownUDPLimitation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ParseUDPLimitation(tt.in)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("got (%q, %v), want %q", got, err, tt.want)
			}
		})
	}
}

func TestNewApplicationPolicy(t *testing.T) {
	proxy, err := domain.ProxyAction("proxy-1")
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}

	tests := []struct {
		name       string
		action     domain.RoutingAction
		scope      domain.ProtocolScope
		limitation domain.UDPLimitation
		wantErr    error
		wantScope  domain.ProtocolScope
		wantLimit  domain.UDPLimitation
	}{
		{
			name:    "proxy_requires_scope",
			action:  proxy,
			scope:   "",
			wantErr: domain.ErrUnsetScope,
		},
		{
			name:       "proxy_both",
			action:     proxy,
			scope:      domain.ProtocolScopeBOTH,
			limitation: "",
			wantScope:  domain.ProtocolScopeBOTH,
			wantLimit:  domain.UDPLimitationAbsent,
		},
		{
			name:      "direct_empty_scope_ok",
			action:    domain.DirectAction(),
			scope:     "",
			wantScope: "",
			wantLimit: domain.UDPLimitationAbsent,
		},
		{
			name:    "unset_action",
			action:  domain.RoutingAction{},
			scope:   domain.ProtocolScopeTCP,
			wantErr: domain.ErrUnsetAction,
		},
		{
			name:    "bad_scope",
			action:  proxy,
			scope:   domain.ProtocolScope("tcp"),
			wantErr: domain.ErrUnknownProtocolScope,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewApplicationPolicy("proc-1", tt.action, tt.scope, tt.limitation)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("error = %v", err)
			}
			if got.ProtocolScope != tt.wantScope || got.UDPLimitation != tt.wantLimit {
				t.Fatalf("got scope=%q limit=%q, want %q %q", got.ProtocolScope, got.UDPLimitation, tt.wantScope, tt.wantLimit)
			}
		})
	}
}

func TestNewProxyEndpointWithUsername(t *testing.T) {
	got, err := domain.NewProxyEndpointWithUsername("p1", "n", domain.ProxyKindSOCKS5, "127.0.0.1", 1080, "alice", "kc")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if got.Username != "alice" {
		t.Fatalf("Username = %q, want alice", got.Username)
	}
}
