package domain

import "fmt"

// BridgeState is the ProxyBridge process dimension.
type BridgeState string

// BridgeState values from CONCEPTION. Unknown is never treated as protected.
const (
	BridgeStateNotInstalled BridgeState = "NotInstalled"
	BridgeStateDetected     BridgeState = "Detected"
	BridgeStateStarting     BridgeState = "Starting"
	BridgeStateReady        BridgeState = "Ready"
	BridgeStateStopping     BridgeState = "Stopping"
	BridgeStateStopped      BridgeState = "Stopped"
	BridgeStateFailed       BridgeState = "Failed"
	BridgeStateConflict     BridgeState = "Conflict"
	BridgeStateUnknown      BridgeState = "Unknown"
)

// ProxyHealth is the per-proxy health dimension.
type ProxyHealth string

// ProxyHealth values from CONCEPTION.
const (
	ProxyHealthUnknown                    ProxyHealth = "Unknown"
	ProxyHealthChecking                   ProxyHealth = "Checking"
	ProxyHealthAvailable                  ProxyHealth = "Available"
	ProxyHealthUnavailable                ProxyHealth = "Unavailable"
	ProxyHealthAuthenticationFailed       ProxyHealth = "AuthenticationFailed"
	ProxyHealthTCPAvailableUDPUnavailable ProxyHealth = "TCPAvailableUDPUnavailable"
)

// ProtectionState is the aggregate protection dimension.
type ProtectionState string

// ProtectionState values from CONCEPTION. Active requires full coverage.
const (
	ProtectionStateInactive     ProtectionState = "Inactive"
	ProtectionStateInitializing ProtectionState = "Initializing"
	ProtectionStateActive       ProtectionState = "Active"
	ProtectionStateDegraded     ProtectionState = "Degraded"
	ProtectionStateUnprotected  ProtectionState = "Unprotected"
	ProtectionStateFailed       ProtectionState = "Failed"
	ProtectionStateUnknown      ProtectionState = "Unknown"
)

// ParseBridgeState returns a BridgeState for an exact CONCEPTION token.
func ParseBridgeState(s string) (BridgeState, error) {
	state := BridgeState(s)
	switch state {
	case BridgeStateNotInstalled,
		BridgeStateDetected,
		BridgeStateStarting,
		BridgeStateReady,
		BridgeStateStopping,
		BridgeStateStopped,
		BridgeStateFailed,
		BridgeStateConflict,
		BridgeStateUnknown:
		return state, nil
	default:
		return "", fmt.Errorf("parse bridge state %q: %w", s, ErrUnknownState)
	}
}

// ParseProxyHealth returns a ProxyHealth for an exact CONCEPTION token.
func ParseProxyHealth(s string) (ProxyHealth, error) {
	health := ProxyHealth(s)
	switch health {
	case ProxyHealthUnknown,
		ProxyHealthChecking,
		ProxyHealthAvailable,
		ProxyHealthUnavailable,
		ProxyHealthAuthenticationFailed,
		ProxyHealthTCPAvailableUDPUnavailable:
		return health, nil
	default:
		return "", fmt.Errorf("parse proxy health %q: %w", s, ErrUnknownState)
	}
}

// ParseProtectionState returns a ProtectionState for an exact CONCEPTION token.
// An unknown token never returns ProtectionStateActive.
func ParseProtectionState(s string) (ProtectionState, error) {
	state := ProtectionState(s)
	switch state {
	case ProtectionStateInactive,
		ProtectionStateInitializing,
		ProtectionStateActive,
		ProtectionStateDegraded,
		ProtectionStateUnprotected,
		ProtectionStateFailed,
		ProtectionStateUnknown:
		return state, nil
	default:
		return "", fmt.Errorf("parse protection state %q: %w", s, ErrUnknownState)
	}
}
