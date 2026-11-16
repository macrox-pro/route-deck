package domain

import "fmt"

// ApplicationPolicy is a policy for one process identity inside a profile.
type ApplicationPolicy struct {
	ProcessID ProcessIdentityID
	Action    RoutingAction
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
