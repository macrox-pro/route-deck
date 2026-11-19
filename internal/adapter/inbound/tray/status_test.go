package tray_test

import (
	"strings"
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/tray"
	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestCLIOwned(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		bridge domain.BridgeState
		want   bool
	}{
		{name: "ready", bridge: domain.BridgeStateReady, want: true},
		{name: "starting", bridge: domain.BridgeStateStarting, want: true},
		{name: "stopping", bridge: domain.BridgeStateStopping, want: true},
		{name: "conflict", bridge: domain.BridgeStateConflict, want: false},
		{name: "stopped", bridge: domain.BridgeStateStopped, want: false},
		{name: "failed", bridge: domain.BridgeStateFailed, want: false},
		{name: "detected", bridge: domain.BridgeStateDetected, want: false},
		{name: "unknown", bridge: domain.BridgeStateUnknown, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tray.CLIOwned(tt.bridge); got != tt.want {
				t.Fatalf("CLIOwned(%s) = %v, want %v", tt.bridge, got, tt.want)
			}
		})
	}
}

func TestStatusLabel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		protection domain.ProtectionState
		bridge     domain.BridgeState
		want       string
	}{
		{
			name:       "ready omits bridge",
			protection: domain.ProtectionStateActive,
			bridge:     domain.BridgeStateReady,
			want:       "Active",
		},
		{
			name:       "appends bridge when not ready",
			protection: domain.ProtectionStateUnprotected,
			bridge:     domain.BridgeStateConflict,
			want:       "Unprotected · Conflict",
		},
		{
			name:       "empty defaults to Unknown",
			protection: "",
			bridge:     "",
			want:       "Unknown · Unknown",
		},
		{
			name:       "starting with initializing",
			protection: domain.ProtectionStateInitializing,
			bridge:     domain.BridgeStateStarting,
			want:       "Initializing · Starting",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tray.StatusLabel(tt.protection, tt.bridge); got != tt.want {
				t.Fatalf("StatusLabel = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestQuitDirtyMessage(t *testing.T) {
	got := tray.QuitDirtyMessage(false)
	if !strings.Contains(got, "unsaved edits") {
		t.Fatalf("got %q", got)
	}
	owned := tray.QuitDirtyMessage(true)
	if !strings.Contains(owned, "Protection will stop") {
		t.Fatalf("owned = %q", owned)
	}
}

func TestQuitMessage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		cliOwned bool
		wantSub  string
	}{
		{name: "owned mentions stop", cliOwned: true, wantSub: "Protection will stop"},
		{name: "not owned short", cliOwned: false, wantSub: "Quit RouteDeck?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tray.QuitMessage(tt.cliOwned)
			if !strings.Contains(got, tt.wantSub) {
				t.Fatalf("QuitMessage(%v) = %q, want substring %q", tt.cliOwned, got, tt.wantSub)
			}
		})
	}
}
