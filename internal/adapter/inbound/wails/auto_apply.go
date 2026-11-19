package wails

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// AutoApplyRetryDelay is the ADR-0017 wait between the first failed Apply and the single retry.
const AutoApplyRetryDelay = 5 * time.Second

// AutoApplyConfig wires startup Auto Apply (ADR-0015 / ADR-0017).
type AutoApplyConfig struct {
	Logger   *slog.Logger
	Settings port.SettingsStore
	Profiles port.ProfileStore
	// HealthSnapshot returns the latest proxy health registry after Apply's Recheck.
	HealthSnapshot func() map[domain.ProxyID]domain.ProxyHealth
	// Apply is the synchronous coordinator Apply (not the async Wails method).
	Apply  func(ctx context.Context, profileID domain.ProfileID) error
	Dialog UserDialog
	// RetryDelay defaults to AutoApplyRetryDelay when non-positive (tests may set a short value).
	RetryDelay time.Duration
}

// MaybeAutoApplyAtStartup evaluates ADR-0015 gates, runs synchronous Apply, and may retry once (ADR-0017).
// MUST NOT rewrite policies to Direct.
func MaybeAutoApplyAtStartup(ctx context.Context, cfg AutoApplyConfig) {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	if cfg.Settings == nil || cfg.Profiles == nil || cfg.Apply == nil {
		logger.Debug("auto apply skipped: not configured", "component", "wails.auto_apply")
		return
	}
	stored, err := cfg.Settings.Get(ctx)
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
	profile, err := cfg.Profiles.Get(ctx, stored.DesiredProfileID)
	if err != nil {
		logger.Warn("auto apply profile failed", "component", "wails.auto_apply", "err", err)
		return
	}
	if !service.ShouldAutoApplyOnStartup(stored.AutoApply, stored.DesiredSet, stored.DesiredProfileID, profile) {
		logger.Info("auto apply skipped: no proxy or block rules", "component", "wails.auto_apply", "profile_id", string(stored.DesiredProfileID))
		return
	}

	profileID := stored.DesiredProfileID
	logger.Info("auto apply starting", "component", "wails.auto_apply", "profile_id", string(profileID))
	err = cfg.Apply(ctx, profileID)
	if err == nil {
		return
	}
	logger.Warn("auto apply failed", "component", "wails.auto_apply", "err", err)

	health := map[domain.ProxyID]domain.ProxyHealth{}
	if cfg.HealthSnapshot != nil {
		health = cfg.HealthSnapshot()
	}
	if !ShouldRetryAutoApply(err, profile, health) {
		warnAutoApplyFailed(cfg.Dialog, err)
		return
	}

	delay := cfg.RetryDelay
	if delay <= 0 {
		delay = AutoApplyRetryDelay
	}
	logger.Info("auto apply retry scheduled", "component", "wails.auto_apply", "delay", delay.String())
	if waitErr := waitRetryDelay(ctx, delay); waitErr != nil {
		logger.Warn("auto apply retry cancelled", "component", "wails.auto_apply", "err", waitErr)
		warnAutoApplyFailed(cfg.Dialog, err)
		return
	}
	logger.Info("auto apply retry starting", "component", "wails.auto_apply", "profile_id", string(profileID))
	if err := cfg.Apply(ctx, profileID); err != nil {
		logger.Warn("auto apply retry failed", "component", "wails.auto_apply", "err", err)
		warnAutoApplyFailed(cfg.Dialog, err)
	}
}

func waitRetryDelay(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// ShouldRetryAutoApply is the ADR-0017 gate: error + referenced proxy Unavailable, excluding Conflict classes.
func ShouldRetryAutoApply(err error, profile domain.Profile, health map[domain.ProxyID]domain.ProxyHealth) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if strings.Contains(msg, string(port.BridgeLineGUIConflict)) ||
		strings.Contains(msg, string(port.BridgeLineAnotherCLI)) ||
		strings.Contains(msg, string(port.BridgeLineExtensionTimeout)) ||
		strings.Contains(msg, "not installed") {
		return false
	}
	return profileReferencesUnavailable(profile, health)
}

func profileReferencesUnavailable(profile domain.Profile, health map[domain.ProxyID]domain.ProxyHealth) bool {
	for _, policy := range profile.Policies {
		id, ok := policy.Action.ProxyID()
		if !ok {
			continue
		}
		if h, ok := health[id]; ok && h == domain.ProxyHealthUnavailable {
			return true
		}
	}
	return false
}

func warnAutoApplyFailed(dialog UserDialog, err error) {
	if dialog == nil {
		return
	}
	if err != nil && strings.Contains(err.Error(), string(port.BridgeLineExtensionTimeout)) {
		dialog.Warning(UserMsgExtensionTimeout)
		return
	}
	dialog.Warning(UserMsgAutoApplyFailed)
}
