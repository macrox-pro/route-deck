package wails

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/logfile"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/buildmeta"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// SettingsService exposes non-secret settings and log helpers.
type SettingsService struct {
	Settings *service.Settings
	Logger   *slog.Logger
}

// SettingsDTO is the settings payload.
type SettingsDTO struct {
	ProxyBridgePath string `json:"proxyBridgePath"`
	// DiscoveredProxyBridgePath comes from BridgeLocator.Locate with an empty
	// custom path (platform documented locations). Settings UI uses it as placeholder.
	DiscoveredProxyBridgePath string `json:"discoveredProxyBridgePath"`
	DebugLogging              bool   `json:"debugLogging"`
	// LaunchAtLogin reflects OS Autostart registration (ADR-0010). Not stored in SQLite.
	LaunchAtLogin bool   `json:"launchAtLogin"`
	DesiredID     string `json:"desiredId"`
	DesiredSet    bool   `json:"desiredSet"`
	BridgeState   string `json:"bridgeState"`
	BridgeVersion string `json:"bridgeVersion"`
	VersionKnown  bool   `json:"versionKnown"`
	CompatWarning string `json:"compatibilityWarning"`
	// AppVersion is the link-time marketing version ("dev" or "1.2.3").
	AppVersion string `json:"appVersion"`
	// AppVersionLabel is the short chrome label ("dev" or "v1.2.3").
	AppVersionLabel string `json:"appVersionLabel"`
	// AppCommit is the short git SHA when injected at link time; empty otherwise.
	AppCommit string `json:"appCommit"`
}

func (s *SettingsService) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

func settingsDTO(view service.SettingsView) SettingsDTO {
	return SettingsDTO{
		ProxyBridgePath:           view.ProxyBridgePath,
		DiscoveredProxyBridgePath: view.DiscoveredProxyBridgePath,
		DebugLogging:              view.DebugLogging,
		LaunchAtLogin:             view.LaunchAtLogin,
		DesiredID:                 string(view.DesiredID),
		DesiredSet:                view.DesiredSet,
		BridgeState:               string(view.BridgeState),
		BridgeVersion:             view.BridgeVersion,
		VersionKnown:              view.VersionKnown,
		CompatWarning:             view.CompatibilityWarning,
		AppVersion:                buildmeta.Version,
		AppVersionLabel:           buildmeta.DisplayVersion(),
		AppCommit:                 buildmeta.Commit,
	}
}

// Get returns current settings plus bridge discovery for the configured path.
func (s *SettingsService) Get(ctx context.Context) (SettingsDTO, error) {
	if s == nil || s.Settings == nil {
		return SettingsDTO{}, fmt.Errorf("settings service: not configured")
	}
	view, err := s.Settings.Get(ctx)
	if err != nil {
		return SettingsDTO{}, err
	}
	return settingsDTO(view), nil
}

// SetProxyBridgePath stores a custom path and rediscovers.
func (s *SettingsService) SetProxyBridgePath(ctx context.Context, path string) (SettingsDTO, error) {
	if s == nil || s.Settings == nil {
		return SettingsDTO{}, fmt.Errorf("settings service: not configured")
	}
	view, err := s.Settings.SetProxyBridgePath(ctx, path)
	if err != nil {
		return SettingsDTO{}, err
	}
	return settingsDTO(view), nil
}

// SetDebugLogging toggles slog LevelVar and persists the flag.
func (s *SettingsService) SetDebugLogging(ctx context.Context, enabled bool) (SettingsDTO, error) {
	if s == nil || s.Settings == nil {
		return SettingsDTO{}, fmt.Errorf("settings service: not configured")
	}
	view, err := s.Settings.SetDebugLogging(ctx, enabled)
	if err != nil {
		return SettingsDTO{}, err
	}
	return settingsDTO(view), nil
}

// SetLaunchAtLogin toggles OS Launch at Login via Autostart. It does not Apply.
func (s *SettingsService) SetLaunchAtLogin(ctx context.Context, enabled bool) (SettingsDTO, error) {
	if s == nil || s.Settings == nil {
		return SettingsDTO{}, fmt.Errorf("settings service: not configured")
	}
	view, err := s.Settings.SetLaunchAtLogin(ctx, enabled)
	if err != nil {
		return SettingsDTO{}, err
	}
	return settingsDTO(view), nil
}

// OpenLogsFolder reveals the log directory in Finder.
func (s *SettingsService) OpenLogsFolder(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	app := application.Get()
	if app == nil {
		return fmt.Errorf("open logs folder: wails application unavailable")
	}
	logPath, err := logfile.DefaultPath()
	if err != nil {
		return fmt.Errorf("open logs folder: %w", err)
	}
	dir := filepath.Dir(logPath)
	if err := app.Env.OpenFileManager(dir, false); err != nil {
		s.logger().Warn("open logs folder failed", "component", "wails.settings", "err", err)
		return fmt.Errorf("open logs folder: %w", err)
	}
	return nil
}
