// Package port holds application boundary interfaces.
package port

import (
	"context"

	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProxyStore loads one proxy endpoint by id.
type ProxyStore interface {
	Get(ctx context.Context, id domain.ProxyID) (domain.ProxyEndpoint, error)
}

// SecretStore loads secret bytes for a Keychain reference.
type SecretStore interface {
	Load(ctx context.Context, ref domain.KeychainRef) ([]byte, error)
}
