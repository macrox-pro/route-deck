package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// Settings is the application use case for non-secret settings and bridge discovery.
type Settings struct {
	Store   port.SettingsStore
	Locator port.BridgeLocator
	Level   *slog.LevelVar
	Logger  *slog.Logger
}

// SettingsView is store state plus locator results for the UI.
type SettingsView struct {
	ProxyBridgePath           string
	DiscoveredProxyBridgePath string
	DebugLogging              bool
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
	stored, err := s.Store.Get(ctx)
	if err != nil {
		return SettingsView{}, fmt.Errorf("settings get: %w", err)
	}
	view := SettingsView{
		ProxyBridgePath: stored.ProxyBridgePath,
		DebugLogging:    stored.DebugLogging,
		DesiredSet:      stored.DesiredSet,
		DesiredID:       stored.DesiredProfileID,
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

// SetProxyBridgePath stores a custom path and returns a refreshed view.
func (s *Settings) SetProxyBridgePath(ctx context.Context, path string) (SettingsView, error) {
	if s == nil || s.Store == nil {
		return SettingsView{}, fmt.Errorf("settings: not configured")
	}
	stored, err := s.Store.Get(ctx)
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
	stored, err := s.Store.Get(ctx)
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
