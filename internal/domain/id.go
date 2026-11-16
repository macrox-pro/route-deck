package domain

import "fmt"

// ProxyID identifies a proxy in the global catalog.
type ProxyID string

// ApplicationID identifies a user-facing application entry.
type ApplicationID string

// ProcessIdentityID identifies a process identity used for matching.
type ProcessIdentityID string

// ProfileID identifies a named profile.
type ProfileID string

// KeychainRef is a Keychain reference for a secret. SQLite stores this, not the secret.
type KeychainRef string

// ParseProxyID returns a ProxyID or ErrEmptyID when s is empty.
func ParseProxyID(s string) (ProxyID, error) {
	if s == "" {
		return "", fmt.Errorf("parse proxy id: %w", ErrEmptyID)
	}
	return ProxyID(s), nil
}

// ParseApplicationID returns an ApplicationID or ErrEmptyID when s is empty.
func ParseApplicationID(s string) (ApplicationID, error) {
	if s == "" {
		return "", fmt.Errorf("parse application id: %w", ErrEmptyID)
	}
	return ApplicationID(s), nil
}

// ParseProcessIdentityID returns a ProcessIdentityID or ErrEmptyID when s is empty.
func ParseProcessIdentityID(s string) (ProcessIdentityID, error) {
	if s == "" {
		return "", fmt.Errorf("parse process identity id: %w", ErrEmptyID)
	}
	return ProcessIdentityID(s), nil
}

// ParseProfileID returns a ProfileID or ErrEmptyID when s is empty.
func ParseProfileID(s string) (ProfileID, error) {
	if s == "" {
		return "", fmt.Errorf("parse profile id: %w", ErrEmptyID)
	}
	return ProfileID(s), nil
}

// ParseKeychainRef returns a KeychainRef or ErrEmptyID when s is empty.
func ParseKeychainRef(s string) (KeychainRef, error) {
	if s == "" {
		return "", fmt.Errorf("parse keychain ref: %w", ErrEmptyID)
	}
	return KeychainRef(s), nil
}
