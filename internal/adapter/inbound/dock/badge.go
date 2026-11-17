package dock

import "github.com/macrox-pro/route-deck/internal/domain"

// SnapshotFacts are the fields badgeFor needs from a live snapshot.
type SnapshotFacts struct {
	AppliedConfirmed bool
	Protection       domain.ProtectionState
}

// BadgeFor returns Dock badge text for the given snapshot facts (ADR-0014).
// Only "!" when a profile is applied and protection is not Active; otherwise empty.
func BadgeFor(facts SnapshotFacts) string {
	if !facts.AppliedConfirmed {
		return ""
	}
	if facts.Protection == domain.ProtectionStateActive {
		return ""
	}
	return "!"
}
