package domain

import "fmt"

// ProxyKind is the proxy protocol type used in the catalog.
type ProxyKind string

// ProxyKind values accepted by ProxyBridge profiles.
const (
	ProxyKindSOCKS5 ProxyKind = "socks5"
	ProxyKindHTTP   ProxyKind = "http"
)

// ProtocolScope is the protocol coverage for a rule.
type ProtocolScope string

// ProtocolScope values match the ProxyBridge profile contract.
const (
	ProtocolScopeTCP  ProtocolScope = "TCP"
	ProtocolScopeUDP  ProtocolScope = "UDP"
	ProtocolScopeBOTH ProtocolScope = "BOTH"
)

// UDPLimitation records whether the user accepted Direct-UDP for a path that cannot carry UDP.
type UDPLimitation string

// UDPLimitation values. Absent means the default compile (UDP BLOCK when the proxy cannot carry UDP).
const (
	UDPLimitationAbsent    UDPLimitation = "Absent"
	UDPLimitationDirectUDP UDPLimitation = "DirectUDP"
)

// RoutingAction is a policy action: Proxy with an id, Direct, or Block.
// The zero value is none of those; it is not a silent Direct fallback.
type RoutingAction struct {
	proxyID ProxyID
	direct  bool
	block   bool
}

// ProxyAction returns a Proxy RoutingAction for a non-empty ProxyID.
func ProxyAction(id ProxyID) (RoutingAction, error) {
	if id == "" {
		return RoutingAction{}, fmt.Errorf("proxy action: %w", ErrEmptyID)
	}
	return RoutingAction{proxyID: id}, nil
}

// DirectAction returns an explicit Direct RoutingAction.
func DirectAction() RoutingAction {
	return RoutingAction{direct: true}
}

// BlockAction returns a Block RoutingAction.
func BlockAction() RoutingAction {
	return RoutingAction{block: true}
}

// IsDirect reports whether the action is explicit Direct.
func (a RoutingAction) IsDirect() bool {
	return a.direct
}

// IsBlock reports whether the action is Block.
func (a RoutingAction) IsBlock() bool {
	return a.block
}

// ProxyID returns the proxy id when the action is Proxy.
func (a RoutingAction) ProxyID() (ProxyID, bool) {
	if a.proxyID == "" || a.direct || a.block {
		return "", false
	}
	return a.proxyID, true
}

// IsSet reports whether the action is Proxy, Direct, or Block.
// The zero value is unset and must not be treated as Direct.
func (a RoutingAction) IsSet() bool {
	return a.proxyID != "" || a.direct || a.block
}

// ParseProxyKind returns a ProxyKind for the exact tokens socks5 or http.
func ParseProxyKind(s string) (ProxyKind, error) {
	kind := ProxyKind(s)
	switch kind {
	case ProxyKindSOCKS5, ProxyKindHTTP:
		return kind, nil
	default:
		return "", fmt.Errorf("parse proxy kind %q: %w", s, ErrUnknownProxyKind)
	}
}

// ParseProtocolScope returns a ProtocolScope for TCP, UDP, or BOTH.
func ParseProtocolScope(s string) (ProtocolScope, error) {
	scope := ProtocolScope(s)
	switch scope {
	case ProtocolScopeTCP, ProtocolScopeUDP, ProtocolScopeBOTH:
		return scope, nil
	default:
		return "", fmt.Errorf("parse protocol scope %q: %w", s, ErrUnknownProtocolScope)
	}
}

// ParseUDPLimitation returns a UDPLimitation for Absent or DirectUDP.
func ParseUDPLimitation(s string) (UDPLimitation, error) {
	limitation := UDPLimitation(s)
	switch limitation {
	case UDPLimitationAbsent, UDPLimitationDirectUDP:
		return limitation, nil
	default:
		return "", fmt.Errorf("parse udp limitation %q: %w", s, ErrUnknownUDPLimitation)
	}
}
