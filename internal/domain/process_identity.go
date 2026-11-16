package domain

import "strings"

// ProcessIdentity holds bundle id, executable name, path, and signing hints for matching.
type ProcessIdentity struct {
	ID             ProcessIdentityID
	ApplicationID  ApplicationID
	BundleID       string
	ExecutableName string
	Path           string
	SigningHint    string
}

// MatchableNames returns bundle id and executable name for processNames patterns.
// Empty strings and strings containing "/" are skipped. Path and SigningHint are omitted.
func (p ProcessIdentity) MatchableNames() []string {
	candidates := []string{p.BundleID, p.ExecutableName}
	names := make([]string, 0, len(candidates))
	for _, name := range candidates {
		if name == "" || strings.Contains(name, "/") {
			continue
		}
		names = append(names, name)
	}
	return names
}
