// Package domain holds RouteDeck value types and parsers with no aggregates.
package domain

import "errors"

// ErrEmptyID is returned when a typed id string is empty.
var ErrEmptyID = errors.New("empty id")

// ErrUnknownState is returned when a state token does not match a known constant.
var ErrUnknownState = errors.New("unknown state")

// ErrUnknownProxyKind is returned when a proxy kind token is not socks5 or http.
var ErrUnknownProxyKind = errors.New("unknown proxy kind")

// ErrInvalidPort is returned when a proxy port is outside 1–65535.
var ErrInvalidPort = errors.New("invalid port")

// ErrEmptyHost is returned when a proxy host is empty.
var ErrEmptyHost = errors.New("empty host")

// ErrUnresolvedProxy is returned when a Proxy action references a proxy id that is not known.
var ErrUnresolvedProxy = errors.New("unresolved proxy")

// ErrEmptyName is returned when a profile name is empty.
var ErrEmptyName = errors.New("empty name")

// ErrUnsetAction is returned when a RoutingAction is the zero value.
var ErrUnsetAction = errors.New("unset routing action")

// ErrUnsetScope is returned when a Proxy policy has no protocol scope.
var ErrUnsetScope = errors.New("unset protocol scope")

// ErrUnknownProtocolScope is returned when a protocol scope token is unknown.
var ErrUnknownProtocolScope = errors.New("unknown protocol scope")

// ErrUnknownUDPLimitation is returned when a UDP limitation token is unknown.
var ErrUnknownUDPLimitation = errors.New("unknown udp limitation")

// ErrReferenced is returned when deleting an entity that a profile still references.
var ErrReferenced = errors.New("entity is referenced")

// ErrDuplicateBundleID is returned when an application bundle id already exists.
var ErrDuplicateBundleID = errors.New("duplicate bundle id")

// ErrDuplicateProfileName is returned when a profile name already exists.
var ErrDuplicateProfileName = errors.New("duplicate profile name")
