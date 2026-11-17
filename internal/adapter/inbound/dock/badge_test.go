package dock

import (
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestBadgeFor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   SnapshotFacts
		want string
	}{
		{
			name: "no applied profile",
			in:   SnapshotFacts{AppliedConfirmed: false, Protection: domain.ProtectionStateUnprotected},
			want: "",
		},
		{
			name: "applied active",
			in:   SnapshotFacts{AppliedConfirmed: true, Protection: domain.ProtectionStateActive},
			want: "",
		},
		{
			name: "applied degraded",
			in:   SnapshotFacts{AppliedConfirmed: true, Protection: domain.ProtectionStateDegraded},
			want: "!",
		},
		{
			name: "applied failed",
			in:   SnapshotFacts{AppliedConfirmed: true, Protection: domain.ProtectionStateFailed},
			want: "!",
		},
		{
			name: "applied unprotected",
			in:   SnapshotFacts{AppliedConfirmed: true, Protection: domain.ProtectionStateUnprotected},
			want: "!",
		},
		{
			name: "applied inactive",
			in:   SnapshotFacts{AppliedConfirmed: true, Protection: domain.ProtectionStateInactive},
			want: "!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := BadgeFor(tt.in); got != tt.want {
				t.Fatalf("BadgeFor() = %q, want %q", got, tt.want)
			}
		})
	}
}
