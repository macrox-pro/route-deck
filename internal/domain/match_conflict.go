package domain

import (
	"fmt"
	"strings"
)

// policyFingerprint is a comparable key for conflict and collapse (ADR-0014).
func policyFingerprint(policy ApplicationPolicy) string {
	scope := policy.ProtocolScope
	if scope == "" {
		scope = ProtocolScopeBOTH
	}
	limitation := policy.UDPLimitation
	if limitation == "" {
		limitation = UDPLimitationAbsent
	}
	var action string
	switch {
	case policy.Action.IsDirect():
		action = "direct"
	case policy.Action.IsBlock():
		action = "block"
	default:
		id, ok := policy.Action.ProxyID()
		if ok {
			action = "proxy:" + string(id)
		} else {
			action = "unset"
		}
	}
	return action + "|" + string(scope) + "|" + string(limitation)
}

// ValidateMatchConflicts returns ErrConflictingMatch when policies assign different
// fingerprints to identities that share a matchable name (bundle id or executable).
func ValidateMatchConflicts(profile Profile, identities []ProcessIdentity) error {
	identityByID := make(map[ProcessIdentityID]ProcessIdentity, len(identities))
	for _, identity := range identities {
		identityByID[identity.ID] = identity
	}
	// matchable name -> fingerprint seen first
	seen := make(map[string]string)
	for _, policy := range profile.Policies {
		if !policy.Action.IsSet() {
			continue
		}
		identity, ok := identityByID[policy.ProcessID]
		if !ok {
			continue
		}
		fp := policyFingerprint(policy)
		for _, name := range identity.MatchableNames() {
			key := strings.ToLower(name)
			prev, exists := seen[key]
			if !exists {
				seen[key] = fp
				continue
			}
			if prev != fp {
				return fmt.Errorf("validate match conflicts %q: %w", name, ErrConflictingMatch)
			}
		}
	}
	return nil
}
