package domain

import "fmt"

// ProxyEndpoint is a global proxy catalog entry with a Keychain ref for the password.
type ProxyEndpoint struct {
	ID          ProxyID
	Name        string
	Kind        ProxyKind
	Host        string
	Port        int
	Username    string
	KeychainRef KeychainRef
	SendDomain  bool
}

// NewProxyEndpoint returns a ProxyEndpoint with SendDomain true and an empty username.
// Empty host returns ErrEmptyHost. Port outside 1–65535 returns ErrInvalidPort.
func NewProxyEndpoint(id ProxyID, name string, kind ProxyKind, host string, port int, ref KeychainRef) (ProxyEndpoint, error) {
	return NewProxyEndpointWithUsername(id, name, kind, host, port, "", ref)
}

// NewProxyEndpointWithUsername returns a ProxyEndpoint with SendDomain true.
// Empty host returns ErrEmptyHost. Port outside 1–65535 returns ErrInvalidPort.
func NewProxyEndpointWithUsername(id ProxyID, name string, kind ProxyKind, host string, port int, username string, ref KeychainRef) (ProxyEndpoint, error) {
	if host == "" {
		return ProxyEndpoint{}, fmt.Errorf("new proxy endpoint: %w", ErrEmptyHost)
	}
	if port < 1 || port > 65535 {
		return ProxyEndpoint{}, fmt.Errorf("new proxy endpoint: %w", ErrInvalidPort)
	}
	return ProxyEndpoint{
		ID:          id,
		Name:        name,
		Kind:        kind,
		Host:        host,
		Port:        port,
		Username:    username,
		KeychainRef: ref,
		SendDomain:  true,
	}, nil
}
