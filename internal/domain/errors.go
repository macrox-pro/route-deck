// Package domain holds RouteDeck value types and parsers with no aggregates.
package domain

import "errors"

// ErrEmptyID is returned when a typed id string is empty.
var ErrEmptyID = errors.New("empty id")

// ErrUnknownState is returned when a state token does not match a known constant.
var ErrUnknownState = errors.New("unknown state")

// ErrUnknownProxyKind is returned when a proxy kind token is not socks5 or http.
var ErrUnknownProxyKind = errors.New("unknown proxy kind")
