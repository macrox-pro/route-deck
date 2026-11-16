package domain_test

import (
	"errors"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestParseBridgeState(t *testing.T) {
	known := []domain.BridgeState{
		domain.BridgeStateNotInstalled,
		domain.BridgeStateDetected,
		domain.BridgeStateStarting,
		domain.BridgeStateReady,
		domain.BridgeStateStopping,
		domain.BridgeStateStopped,
		domain.BridgeStateFailed,
		domain.BridgeStateConflict,
		domain.BridgeStateUnknown,
	}
	for _, want := range known {
		t.Run("known_"+string(want), func(t *testing.T) {
			got, err := domain.ParseBridgeState(string(want))
			if err != nil {
				t.Fatalf("ParseBridgeState(%q) error = %v", want, err)
			}
			if got != want {
				t.Fatalf("ParseBridgeState(%q) = %q, want %q", want, got, want)
			}
		})
	}

	tests := []struct {
		name string
		in   string
	}{
		{name: "unknown_token", in: "Bogus"},
		{name: "empty", in: ""},
		{name: "case_variant", in: "ready"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ParseBridgeState(tt.in)
			if !errors.Is(err, domain.ErrUnknownState) {
				t.Fatalf("ParseBridgeState(%q) error = %v, want %v", tt.in, err, domain.ErrUnknownState)
			}
			if got != "" {
				t.Fatalf("ParseBridgeState(%q) = %q, want empty", tt.in, got)
			}
		})
	}
}

func TestParseProxyHealth(t *testing.T) {
	known := []domain.ProxyHealth{
		domain.ProxyHealthUnknown,
		domain.ProxyHealthChecking,
		domain.ProxyHealthAvailable,
		domain.ProxyHealthUnavailable,
		domain.ProxyHealthAuthenticationFailed,
		domain.ProxyHealthTCPAvailableUDPUnavailable,
	}
	for _, want := range known {
		t.Run("known_"+string(want), func(t *testing.T) {
			got, err := domain.ParseProxyHealth(string(want))
			if err != nil {
				t.Fatalf("ParseProxyHealth(%q) error = %v", want, err)
			}
			if got != want {
				t.Fatalf("ParseProxyHealth(%q) = %q, want %q", want, got, want)
			}
		})
	}

	got, err := domain.ParseProxyHealth("NotAHealth")
	if !errors.Is(err, domain.ErrUnknownState) {
		t.Fatalf("ParseProxyHealth(unknown) error = %v, want %v", err, domain.ErrUnknownState)
	}
	if got != "" {
		t.Fatalf("ParseProxyHealth(unknown) = %q, want empty", got)
	}
}

func TestParseProtectionState(t *testing.T) {
	known := []domain.ProtectionState{
		domain.ProtectionStateInactive,
		domain.ProtectionStateInitializing,
		domain.ProtectionStateActive,
		domain.ProtectionStateDegraded,
		domain.ProtectionStateUnprotected,
		domain.ProtectionStateFailed,
		domain.ProtectionStateUnknown,
	}
	for _, want := range known {
		t.Run("known_"+string(want), func(t *testing.T) {
			got, err := domain.ParseProtectionState(string(want))
			if err != nil {
				t.Fatalf("ParseProtectionState(%q) error = %v", want, err)
			}
			if got != want {
				t.Fatalf("ParseProtectionState(%q) = %q, want %q", want, got, want)
			}
		})
	}

	tests := []struct {
		name string
		in   string
	}{
		{name: "unknown_token", in: "MaybeProtected"},
		{name: "empty", in: ""},
		{name: "case_variant_active", in: "active"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ParseProtectionState(tt.in)
			if !errors.Is(err, domain.ErrUnknownState) {
				t.Fatalf("ParseProtectionState(%q) error = %v, want %v", tt.in, err, domain.ErrUnknownState)
			}
			if got == domain.ProtectionStateActive {
				t.Fatalf("ParseProtectionState(%q) returned Active on error", tt.in)
			}
			if got != "" {
				t.Fatalf("ParseProtectionState(%q) = %q, want empty", tt.in, got)
			}
		})
	}
}
