package service

import "github.com/macrox-pro/route-deck/internal/domain"

// ProfileHasProtectableRules reports whether any policy uses Proxy or Block (ADR-0015 gate 3).
func ProfileHasProtectableRules(profile domain.Profile) bool {
	for _, policy := range profile.Policies {
		if policy.Action.IsBlock() {
			return true
		}
		if _, isProxy := policy.Action.ProxyID(); isProxy {
			return true
		}
	}
	return false
}

// ShouldAutoApplyOnStartup is the ADR-0015 closed gate list (without invoking Apply).
// Callers that pass must use the existing ProtectionService.Apply path.
func ShouldAutoApplyOnStartup(autoApply, desiredSet bool, desiredID domain.ProfileID, profile domain.Profile) bool {
	if !autoApply {
		return false
	}
	if !desiredSet || desiredID == "" {
		return false
	}
	if profile.ID != "" && profile.ID != desiredID {
		return false
	}
	return ProfileHasProtectableRules(profile)
}
