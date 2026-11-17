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
	Delete(ctx context.Context, id domain.ProxyID) error
}

// ApplicationStore loads and saves user-facing applications.
type ApplicationStore interface {
	Get(ctx context.Context, id domain.ApplicationID) (domain.Application, error)
	List(ctx context.Context) ([]domain.Application, error)
	Save(ctx context.Context, app domain.Application) error
	Delete(ctx context.Context, id domain.ApplicationID) error
}

// ProcessIdentityStore loads and saves process identities.
type ProcessIdentityStore interface {
	Get(ctx context.Context, id domain.ProcessIdentityID) (domain.ProcessIdentity, error)
	List(ctx context.Context) ([]domain.ProcessIdentity, error)
	Save(ctx context.Context, identity domain.ProcessIdentity) error
	Delete(ctx context.Context, id domain.ProcessIdentityID) error
}

// ProfileStore loads and saves profiles with their policies.
type ProfileStore interface {
	Get(ctx context.Context, id domain.ProfileID) (domain.Profile, error)
	List(ctx context.Context) ([]domain.Profile, error)
	Save(ctx context.Context, profile domain.Profile) error
	Delete(ctx context.Context, id domain.ProfileID) error
}

// Settings holds persisted application settings (non-secret).
type Settings struct {
	DesiredProfileID domain.ProfileID
	DesiredSet       bool
	ProxyBridgePath  string
	DebugLogging     bool
}

// SettingsStore loads and saves Settings.
type SettingsStore interface {
	Get(ctx context.Context) (Settings, error)
	Save(ctx context.Context, settings Settings) error
}

// SecretStore loads and stores secret bytes for a Keychain reference.
type SecretStore interface {
	Load(ctx context.Context, ref domain.KeychainRef) ([]byte, error)
	Save(ctx context.Context, ref domain.KeychainRef, secret []byte) error
	Delete(ctx context.Context, ref domain.KeychainRef) error
}

// InspectedProcess is one process candidate from an application bundle.
type InspectedProcess struct {
	BundleID       string
	ExecutableName string
	Path           string
	DisplayName    string
	IsHelper       bool
	// IsUpdater is set by the platform adapter (for example Sparkle on macOS).
	IsUpdater bool
	// RealPath is the EvalSymlinks path used for dedupe. Empty when unused.
	RealPath string
}

// InspectedApplication is the result of inspecting an application path.
type InspectedApplication struct {
	DisplayName string
	BundleID    string
	Main        InspectedProcess
	Helpers     []InspectedProcess
}

// ApplicationInspector reads platform-specific application metadata.
type ApplicationInspector interface {
	Inspect(ctx context.Context, path string) (InspectedApplication, error)
}
