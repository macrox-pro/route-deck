// Package tray is the composition-root adapter for the macOS menu-bar tray (ADR-0010).
package tray

import "github.com/macrox-pro/route-deck/internal/domain"

// CLIOwned reports whether RouteDeck currently owns the managed ProxyBridge CLI.
// Conflict is never owned.
func CLIOwned(bridge domain.BridgeState) bool {
	switch bridge {
	case domain.BridgeStateReady, domain.BridgeStateStarting, domain.BridgeStateStopping:
		return true
	case domain.BridgeStateNotInstalled,
		domain.BridgeStateDetected,
		domain.BridgeStateStopped,
		domain.BridgeStateFailed,
		domain.BridgeStateConflict,
		domain.BridgeStateUnknown:
		return false
	default:
		return false
	}
}

// StatusLabel builds the disabled Status menu text from domain enum words only.
func StatusLabel(protection domain.ProtectionState, bridge domain.BridgeState) string {
	if protection == "" {
		protection = domain.ProtectionStateUnknown
	}
	if bridge == "" {
		bridge = domain.BridgeStateUnknown
	}
	if bridge == domain.BridgeStateReady {
		return string(protection)
	}
	return string(protection) + " · " + string(bridge)
}

// QuitMessage is the Question dialog body for Quit.
func QuitMessage(cliOwned bool) string {
	if cliOwned {
		return "Quit RouteDeck? Protection will stop and the managed ProxyBridge process will exit."
	}
	return "Quit RouteDeck?"
}

// QuitDirtyMessage is the Question body when an editor has unsaved edits (ADR-0018).
func QuitDirtyMessage(cliOwned bool) string {
	base := "You have unsaved edits. Quitting discards them."
	if cliOwned {
		return base + " Protection will stop and the managed ProxyBridge process will exit."
	}
	return base
}
