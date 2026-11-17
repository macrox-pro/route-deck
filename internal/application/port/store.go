// Package port holds application boundary interfaces.
package port

import (
	"context"

	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProxyStore loads and saves proxy endpoints by id.
type ProxyStore interface {
	Get(ctx context.Context, id domain.ProxyID) (domain.ProxyEndpoint, error)
	List(ctx context.Context) ([]domain.ProxyEndpoint, error)
	Save(ctx context.Context, endpoint domain.ProxyEndpoint) error
}

// ApplicationStore loads and saves user-facing applications.
type ApplicationStore interface {
	Get(ctx context.Context, id domain.ApplicationID) (domain.Application, error)
	List(ctx context.Context) ([]domain.Application, error)
	Save(ctx context.Context, app domain.Application) error
}

// ProcessIdentityStore loads and saves process identities.
type ProcessIdentityStore interface {
	Get(ctx context.Context, id domain.ProcessIdentityID) (domain.ProcessIdentity, error)
	List(ctx context.Context) ([]domain.ProcessIdentity, error)
	Save(ctx context.Context, identity domain.ProcessIdentity) error
}

// ProfileStore loads and saves profiles with their policies.
type ProfileStore interface {
	Get(ctx context.Context, id domain.ProfileID) (domain.Profile, error)
	List(ctx context.Context) ([]domain.Profile, error)
	Save(ctx context.Context, profile domain.Profile) error
}

// SecretStore loads secret bytes for a Keychain reference.
type SecretStore interface {
	Load(ctx context.Context, ref domain.KeychainRef) ([]byte, error)
}
