package domain

import (
	"fmt"
	"strings"
)

// policyFingerprint is a comparable key for conflict and collapse (ADR-0014 / ADR-0019).
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
// Prefer ResolveMatchConflictsLastWins before Save/Apply (ADR-0019).
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

// ResolveMatchConflictsLastWins returns a policy list where, for each matchable name,
// the fingerprint of the last policy in profile.Policies that touches that name wins.
// Earlier policies that would conflict are dropped. Order of kept policies is preserved
// relative to the original list (ADR-0019). Unknown process ids are kept as-is.
func ResolveMatchConflictsLastWins(profile Profile, identities []ProcessIdentity) []ApplicationPolicy {
	identityByID := make(map[ProcessIdentityID]ProcessIdentity, len(identities))
	for _, identity := range identities {
		identityByID[identity.ID] = identity
	}

	// Walk last → first: first claim on a name wins (that is last in original order).
	type claim struct {
		fp string
	}
	claimed := make(map[string]claim)
	keep := make([]bool, len(profile.Policies))
	for i := len(profile.Policies) - 1; i >= 0; i-- {
		policy := profile.Policies[i]
		if !policy.Action.IsSet() {
			keep[i] = true
			continue
		}
		identity, ok := identityByID[policy.ProcessID]
		if !ok {
			keep[i] = true
			continue
		}
		fp := policyFingerprint(policy)
		names := identity.MatchableNames()
		conflict := false
		for _, name := range names {
			key := strings.ToLower(name)
			if c, exists := claimed[key]; exists && c.fp != fp {
				conflict = true
				break
			}
		}
		if conflict {
			keep[i] = false
			continue
		}
		keep[i] = true
		for _, name := range names {
			key := strings.ToLower(name)
			if _, exists := claimed[key]; !exists {
				claimed[key] = claim{fp: fp}
			}
		}
	}

	out := make([]ApplicationPolicy, 0, len(profile.Policies))
	for i, policy := range profile.Policies {
		if keep[i] {
			out = append(out, policy)
		}
	}
	return out
}
