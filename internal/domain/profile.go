package domain

import "fmt"

// ApplicationPolicy is a policy for one process identity inside a profile.
type ApplicationPolicy struct {
	ProcessID     ProcessIdentityID
	Action        RoutingAction
	ProtocolScope ProtocolScope
	UDPLimitation UDPLimitation
}

// NewApplicationPolicy returns a policy. Proxy actions require a non-empty ProtocolScope.
// Direct and Block may omit scope (treated as BOTH when compiling). Empty UDPLimitation
// is normalized to Absent.
func NewApplicationPolicy(processID ProcessIdentityID, action RoutingAction, scope ProtocolScope, limitation UDPLimitation) (ApplicationPolicy, error) {
	if !action.IsSet() {
		return ApplicationPolicy{}, fmt.Errorf("new application policy: %w", ErrUnsetAction)
	}
	if _, isProxy := action.ProxyID(); isProxy && scope == "" {
		return ApplicationPolicy{}, fmt.Errorf("new application policy: %w", ErrUnsetScope)
	}
	if limitation == "" {
		limitation = UDPLimitationAbsent
	}
	switch limitation {
	case UDPLimitationAbsent, UDPLimitationDirectUDP:
	default:
		return ApplicationPolicy{}, fmt.Errorf("new application policy: %w", ErrUnknownUDPLimitation)
	}
	if scope != "" {
		switch scope {
		case ProtocolScopeTCP, ProtocolScopeUDP, ProtocolScopeBOTH:
		default:
			return ApplicationPolicy{}, fmt.Errorf("new application policy: %w", ErrUnknownProtocolScope)
		}
	}
	return ApplicationPolicy{
		ProcessID:     processID,
		Action:        action,
		ProtocolScope: scope,
		UDPLimitation: limitation,
	}, nil
}

// Profile is a named set of application policies that reference proxy ids.
type Profile struct {
	ID       ProfileID
	Name     string
	Policies []ApplicationPolicy
}

// NewProfile returns a Profile. Empty name returns ErrEmptyName.
// Empty ID, empty policies, and policy contents are not checked.
func NewProfile(id ProfileID, name string, policies []ApplicationPolicy) (Profile, error) {
	if name == "" {
		return Profile{}, fmt.Errorf("new profile: %w", ErrEmptyName)
	}
	return Profile{
		ID:       id,
		Name:     name,
		Policies: policies,
	}, nil
}

// ValidateProxyRefs returns ErrUnresolvedProxy when a Proxy action's id is not in known.
// It does not change Action. Direct, Block, and zero RoutingAction are not proxy refs.
func ValidateProxyRefs(profile Profile, known []ProxyID) error {
	for _, policy := range profile.Policies {
		id, ok := policy.Action.ProxyID()
		if !ok {
			continue
		}
		found := false
		for _, k := range known {
			if k == id {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("validate proxy refs: %w", ErrUnresolvedProxy)
		}
	}
	return nil
}
