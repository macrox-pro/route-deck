package port

import (
	"context"

	"github.com/macrox-pro/route-deck/internal/domain"
)

// BridgeDiscovery is the result of locating ProxyBridge.
type BridgeDiscovery struct {
	State                domain.BridgeState
	Path                 string
	Version              string
	VersionKnown         bool
	CompatibilityWarning string
}

// BridgeLocator finds the ProxyBridge binary and probes --version.
// Platform install locations live in the adapter (PlatformPaths); callers MUST
// NOT hard-code OS paths. Locate with an empty customPath searches only
// PlatformPaths.
type BridgeLocator interface {
	Locate(ctx context.Context, customPath string) (BridgeDiscovery, error)
	PlatformPaths() []string
}

// BridgeLineClass is a classified CLI stdout/stderr line.
type BridgeLineClass string

// Bridge line classifications from CLI stdout/stderr.
const (
	BridgeLineUnclassified      BridgeLineClass = "unclassified"
	BridgeLineRedirectReady     BridgeLineClass = "redirect_ready"
	BridgeLineProfileLoaded     BridgeLineClass = "profile_loaded"
	BridgeLineGUIConflict       BridgeLineClass = "gui_conflict"
	BridgeLineAnotherCLI        BridgeLineClass = "another_cli"
	BridgeLineTunnelStopped     BridgeLineClass = "tunnel_stopped"
	BridgeLineExtensionApproval BridgeLineClass = "extension_approval"
	BridgeLineExtensionTimeout  BridgeLineClass = "extension_timeout"
	BridgeLineRejectedRules     BridgeLineClass = "rejected_rules"
)

// BridgeProcess manages one ProxyBridge CLI process.
type BridgeProcess interface {
	Start(ctx context.Context, binaryPath string, profileJSON []byte) error
	Stop(ctx context.Context) error
	WaitConfirm(ctx context.Context) error
	Running() bool
}

// CompiledProfile is adapter-owned stdin JSON plus counts for logging.
type CompiledProfile struct {
	JSON       []byte
	ProxyCount int
	RuleCount  int
}

// ProfileCompiler turns domain inputs into ProxyBridge stdin JSON.
type ProfileCompiler interface {
	Compile(ctx context.Context, profile domain.Profile, proxies []domain.ProxyEndpoint, identities []domain.ProcessIdentity) (CompiledProfile, error)
}

// ProxyHealthChecker probes one proxy endpoint.
type ProxyHealthChecker interface {
	Check(ctx context.Context, endpoint domain.ProxyEndpoint, password []byte) (domain.ProxyHealth, error)
}
