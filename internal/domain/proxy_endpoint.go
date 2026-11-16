package domain

import "fmt"

// ProxyEndpoint is a global proxy catalog entry with a Keychain ref for the password.
type ProxyEndpoint struct {
	ID          ProxyID
	Name        string
	Kind        ProxyKind
	Host        string
	Port        int
	KeychainRef KeychainRef
	SendDomain  bool
}

// NewProxyEndpoint returns a ProxyEndpoint with SendDomain true.
// Empty host returns ErrEmptyHost. Port outside 1–65535 returns ErrInvalidPort.
func NewProxyEndpoint(id ProxyID, name string, kind ProxyKind, host string, port int, ref KeychainRef) (ProxyEndpoint, error) {
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
		KeychainRef: ref,
		SendDomain:  true,
	}, nil
}
