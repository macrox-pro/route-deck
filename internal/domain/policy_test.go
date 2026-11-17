package domain_test

import (
	"errors"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestParseProxyKind(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    domain.ProxyKind
		wantErr error
	}{
		{name: "socks5", in: "socks5", want: domain.ProxyKindSOCKS5},
		{name: "http", in: "http", want: domain.ProxyKindHTTP},
		{name: "rejected_https", in: "https", wantErr: domain.ErrUnknownProxyKind},
		{name: "rejected_empty", in: "", wantErr: domain.ErrUnknownProxyKind},
		{name: "rejected_case", in: "SOCKS5", wantErr: domain.ErrUnknownProxyKind},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ParseProxyKind(tt.in)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("ParseProxyKind(%q) error = %v", tt.in, err)
				}
				if got != tt.want {
					t.Fatalf("ParseProxyKind(%q) = %q, want %q", tt.in, got, tt.want)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ParseProxyKind(%q) error = %v, want %v", tt.in, err, tt.wantErr)
			}
			if got != "" {
				t.Fatalf("ParseProxyKind(%q) = %q, want empty", tt.in, got)
			}
		})
	}
}

func TestRoutingAction(t *testing.T) {
	t.Run("proxy_action_empty_id", func(t *testing.T) {
		got, err := domain.ProxyAction("")
		if !errors.Is(err, domain.ErrEmptyID) {
			t.Fatalf("ProxyAction(\"\") error = %v, want %v", err, domain.ErrEmptyID)
		}
		if got.IsDirect() {
			t.Fatal("ProxyAction(\"\") on error must not be Direct")
		}
		if got.IsBlock() {
			t.Fatal("ProxyAction(\"\") on error must not be Block")
		}
		if _, ok := got.ProxyID(); ok {
			t.Fatal("ProxyAction(\"\") on error must not report a proxy id")
		}
	})

	t.Run("zero_value_not_direct", func(t *testing.T) {
		var zero domain.RoutingAction
		if zero.IsDirect() {
			t.Fatal("zero RoutingAction must not be Direct")
		}
		if zero.IsBlock() {
			t.Fatal("zero RoutingAction must not be Block")
		}
		if id, ok := zero.ProxyID(); ok {
			t.Fatalf("zero RoutingAction reported proxy id %q", id)
		}
		if zero.IsSet() {
			t.Fatal("zero RoutingAction must not be set")
		}
	})

	t.Run("is_set", func(t *testing.T) {
		proxy, err := domain.ProxyAction("proxy-1")
		if err != nil {
			t.Fatalf("ProxyAction: %v", err)
		}
		tests := []struct {
			name string
			a    domain.RoutingAction
			want bool
		}{
			{name: "proxy", a: proxy, want: true},
			{name: "direct", a: domain.DirectAction(), want: true},
			{name: "block", a: domain.BlockAction(), want: true},
			{name: "zero", a: domain.RoutingAction{}, want: false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.a.IsSet(); got != tt.want {
					t.Fatalf("IsSet() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("proxy_action_valid", func(t *testing.T) {
		got, err := domain.ProxyAction("proxy-1")
		if err != nil {
			t.Fatalf("ProxyAction(proxy-1) error = %v", err)
		}
		if got.IsDirect() || got.IsBlock() {
			t.Fatal("ProxyAction must not be Direct or Block")
		}
		id, ok := got.ProxyID()
		if !ok || id != "proxy-1" {
			t.Fatalf("ProxyAction.ProxyID() = (%q, %v), want (proxy-1, true)", id, ok)
		}
	})

	t.Run("direct_action", func(t *testing.T) {
		got := domain.DirectAction()
		if !got.IsDirect() {
			t.Fatal("DirectAction().IsDirect() = false")
		}
		if got.IsBlock() {
			t.Fatal("DirectAction must not be Block")
		}
		if _, ok := got.ProxyID(); ok {
			t.Fatal("DirectAction must not report a proxy id")
		}
	})

	t.Run("block_action", func(t *testing.T) {
		got := domain.BlockAction()
		if !got.IsBlock() {
			t.Fatal("BlockAction().IsBlock() = false")
		}
		if got.IsDirect() {
			t.Fatal("BlockAction must not be Direct")
		}
		if _, ok := got.ProxyID(); ok {
			t.Fatal("BlockAction must not report a proxy id")
		}
	})
}

func TestProtocolScopeAndUDPLimitationConstants(t *testing.T) {
	if domain.ProtocolScopeTCP != "TCP" || domain.ProtocolScopeUDP != "UDP" || domain.ProtocolScopeBOTH != "BOTH" {
		t.Fatalf("ProtocolScope constants drifted: %q %q %q", domain.ProtocolScopeTCP, domain.ProtocolScopeUDP, domain.ProtocolScopeBOTH)
	}
	if domain.UDPLimitationAbsent != "Absent" || domain.UDPLimitationDirectUDP != "DirectUDP" {
		t.Fatalf("UDPLimitation constants drifted: %q %q", domain.UDPLimitationAbsent, domain.UDPLimitationDirectUDP)
	}
}
