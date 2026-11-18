package wails

import (
	"context"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
)

// MaybeAutoApplyAtStartup evaluates ADR-0015 gates once, then calls Apply when they pass.
// Failures use the existing Apply / event path. MUST NOT rewrite policies to Direct.
func MaybeAutoApplyAtStartup(
	ctx context.Context,
	logger *slog.Logger,
	settingsStore port.SettingsStore,
	profiles port.ProfileStore,
	apply func(context.Context) (ApplyResultDTO, error),
) {
	if logger == nil {
		logger = slog.Default()
	}
	if settingsStore == nil || profiles == nil || apply == nil {
		logger.Debug("auto apply skipped: not configured", "component", "wails.auto_apply")
		return
	}
	stored, err := settingsStore.Get(ctx)
	if err != nil {
		logger.Warn("auto apply settings failed", "component", "wails.auto_apply", "err", err)
		return
	}
	if !stored.AutoApply {
		logger.Info("auto apply skipped: disabled", "component", "wails.auto_apply")
		return
	}
	if !stored.DesiredSet || stored.DesiredProfileID == "" {
		logger.Info("auto apply skipped: desired profile not set", "component", "wails.auto_apply")
		return
	}
	profile, err := profiles.Get(ctx, stored.DesiredProfileID)
	if err != nil {
		logger.Warn("auto apply profile failed", "component", "wails.auto_apply", "err", err)
		return
	}
	if !service.ShouldAutoApplyOnStartup(stored.AutoApply, stored.DesiredSet, stored.DesiredProfileID, profile) {
		logger.Info("auto apply skipped: no proxy or block rules", "component", "wails.auto_apply", "profile_id", string(stored.DesiredProfileID))
		return
	}
	logger.Info("auto apply starting", "component", "wails.auto_apply", "profile_id", string(stored.DesiredProfileID))
	if _, err := apply(ctx); err != nil {
		logger.Warn("auto apply request failed", "component", "wails.auto_apply", "err", err)
	}
}
