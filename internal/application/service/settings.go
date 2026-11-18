package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// Appearance preference values persisted in SQLite (UI chrome).
const (
	AppearanceSystem = "system"
	AppearanceLight  = "light"
	AppearanceDark   = "dark"
)

// NormalizeAppearance maps empty to system; returns ok=false for unknown values.
func NormalizeAppearance(raw string) (string, bool) {
	switch raw {
	case "", AppearanceSystem:
		return AppearanceSystem, true
	case AppearanceLight:
		return AppearanceLight, true
	case AppearanceDark:
		return AppearanceDark, true
	default:
		return "", false
	}
}

// CoerceAppearance always returns system, light, or dark. Unknown values become system.
func CoerceAppearance(raw string) string {
	if normalized, ok := NormalizeAppearance(raw); ok {
		return normalized
	}
	return AppearanceSystem
}

// Settings is the application use case for non-secret settings and bridge discovery.
type Settings struct {
	Store         port.SettingsStore
	Locator       port.BridgeLocator
	LaunchAtLogin port.LaunchAtLogin
	Level         *slog.LevelVar
	Logger        *slog.Logger
}

// SettingsView is store state plus locator results for the UI.
type SettingsView struct {
	ProxyBridgePath           string
	DiscoveredProxyBridgePath string
	DebugLogging              bool
	LaunchAtLogin             bool
	Appearance                string
	DesiredID                 domain.ProfileID
	DesiredSet                bool
	BridgeState               domain.BridgeState
	BridgeVersion             string
	VersionKnown              bool
	CompatibilityWarning      string
}

func (s *Settings) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// Get loads settings and resolves bridge paths through BridgeLocator only.
// DiscoveredProxyBridgePath is Locate with an empty custom path (platform paths).
func (s *Settings) Get(ctx context.Context) (SettingsView, error) {
	if s == nil || s.Store == nil {
		return SettingsView{}, fmt.Errorf("settings: not configured")
	}
	stored, err := s.loadSanitized(ctx)
	if err != nil {
		return SettingsView{}, fmt.Errorf("settings get: %w", err)
	}
	view := SettingsView{
		ProxyBridgePath: stored.ProxyBridgePath,
		DebugLogging:    stored.DebugLogging,
		Appearance:      stored.Appearance,
		DesiredSet:      stored.DesiredSet,
		DesiredID:       stored.DesiredProfileID,
	}
	if s.LaunchAtLogin != nil {
		enabled, launchErr := s.LaunchAtLogin.IsEnabled(ctx)
		if launchErr != nil {
			s.logger().Debug("launch at login status failed", "component", "settings", "err", launchErr)
		} else {
			view.LaunchAtLogin = enabled
		}
	}
	if s.Locator == nil {
		return view, nil
	}
	if discovered, discErr := s.Locator.Locate(ctx, ""); discErr != nil {
		s.logger().Debug("platform bridge locate failed", "component", "settings", "err", discErr)
	} else {
		view.DiscoveredProxyBridgePath = discovered.Path
	}
	disc, locErr := s.Locator.Locate(ctx, stored.ProxyBridgePath)
	if locErr != nil {
		s.logger().Debug("configured bridge locate failed", "component", "settings", "err", locErr)
		return view, nil
	}
	view.BridgeState = disc.State
	view.BridgeVersion = disc.Version
	view.VersionKnown = disc.VersionKnown
	view.CompatibilityWarning = disc.CompatibilityWarning
	return view, nil
}

// loadSanitized loads settings and rewrites the SQLite row when values are invalid.
func (s *Settings) loadSanitized(ctx context.Context) (port.Settings, error) {
	stored, err := s.Store.Get(ctx)
	if err != nil {
		return port.Settings{}, err
	}
	cleaned, changed := sanitizeSettings(stored)
	if changed {
		if saveErr := s.Store.Save(ctx, cleaned); saveErr != nil {
			s.logger().Warn("settings sanitize save failed", "component", "settings", "err", saveErr)
			// Still return cleaned values so the app keeps running.
		}
	}
	return cleaned, nil
}

func sanitizeSettings(in port.Settings) (port.Settings, bool) {
	out := in
	changed := false
	appearance := CoerceAppearance(in.Appearance)
	if in.Appearance != appearance {
		out.Appearance = appearance
		changed = true
	}
	if out.DesiredSet && out.DesiredProfileID == "" {
		out.DesiredSet = false
		changed = true
	}
	return out, changed
}

// SetProxyBridgePath stores a custom path and returns a refreshed view.
func (s *Settings) SetProxyBridgePath(ctx context.Context, path string) (SettingsView, error) {
	if s == nil || s.Store == nil {
		return SettingsView{}, fmt.Errorf("settings: not configured")
	}
	stored, err := s.loadSanitized(ctx)
	if err != nil {
		return SettingsView{}, fmt.Errorf("settings path: %w", err)
	}
	stored.ProxyBridgePath = path
	if err := s.Store.Save(ctx, stored); err != nil {
		return SettingsView{}, fmt.Errorf("settings path: %w", err)
	}
	s.logger().Info("proxybridge path updated", "component", "settings")
	return s.Get(ctx)
}

// SetDebugLogging toggles slog LevelVar (when set) and persists the flag.
func (s *Settings) SetDebugLogging(ctx context.Context, enabled bool) (SettingsView, error) {
	if s == nil || s.Store == nil {
		return SettingsView{}, fmt.Errorf("settings: not configured")
	}
	stored, err := s.loadSanitized(ctx)
	if err != nil {
		return SettingsView{}, fmt.Errorf("settings debug: %w", err)
	}
	stored.DebugLogging = enabled
	if err := s.Store.Save(ctx, stored); err != nil {
		return SettingsView{}, fmt.Errorf("settings debug: %w", err)
	}
	if s.Level != nil {
		if enabled {
			s.Level.Set(slog.LevelDebug)
		} else {
			s.Level.Set(slog.LevelInfo)
		}
	}
	s.logger().Info("debug logging updated", "component", "settings", "enabled", enabled)
	return s.Get(ctx)
}

// SetLaunchAtLogin registers or removes OS login launch. It does not Apply.
func (s *Settings) SetLaunchAtLogin(ctx context.Context, enabled bool) (SettingsView, error) {
	if s == nil || s.LaunchAtLogin == nil {
		return SettingsView{}, fmt.Errorf("settings: launch at login not configured")
	}
	if err := s.LaunchAtLogin.SetEnabled(ctx, enabled); err != nil {
		return SettingsView{}, fmt.Errorf("settings launch at login: %w", err)
	}
	s.logger().Info("launch at login updated", "component", "settings", "enabled", enabled)
	return s.Get(ctx)
}

// SetAppearance persists UI appearance (system|light|dark). It does not Apply protection.
// Unknown values are coerced to system and written so SQLite config stays valid.
func (s *Settings) SetAppearance(ctx context.Context, appearance string) (SettingsView, error) {
	if s == nil || s.Store == nil {
		return SettingsView{}, fmt.Errorf("settings: not configured")
	}
	normalized := CoerceAppearance(appearance)
	stored, err := s.loadSanitized(ctx)
	if err != nil {
		return SettingsView{}, fmt.Errorf("settings appearance: %w", err)
	}
	stored.Appearance = normalized
	if err := s.Store.Save(ctx, stored); err != nil {
		return SettingsView{}, fmt.Errorf("settings appearance: %w", err)
	}
	s.logger().Info("appearance updated", "component", "settings", "appearance", normalized)
	return s.Get(ctx)
}
