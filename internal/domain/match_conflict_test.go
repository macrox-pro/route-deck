package domain_test

import (
	"errors"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestResolveMatchConflictsLastWins(t *testing.T) {
	identities := []domain.ProcessIdentity{
		{ID: "h1", ApplicationID: "app", BundleID: "com.example.helper", ExecutableName: "HelperA"},
		{ID: "h2", ApplicationID: "app", BundleID: "com.example.helper", ExecutableName: "HelperB"},
		{ID: "m1", ApplicationID: "app", BundleID: "com.example.app", ExecutableName: "App"},
	}

	t.Run("last_wins_shared_bundle", func(t *testing.T) {
		profile := domain.Profile{
			ID:   "p",
			Name: "Work",
			Policies: []domain.ApplicationPolicy{
				{ProcessID: "h1", Action: domain.BlockAction()},
				{ProcessID: "h2", Action: domain.DirectAction()},
			},
		}
		got := domain.ResolveMatchConflictsLastWins(profile, identities)
		if len(got) != 1 {
			t.Fatalf("len = %d, want 1: %+v", len(got), got)
		}
		if got[0].ProcessID != "h2" || !got[0].Action.IsDirect() {
			t.Fatalf("got = %+v, want h2 Direct", got[0])
		}
		resolved := profile
		resolved.Policies = got
		if err := domain.ValidateMatchConflicts(resolved, identities); err != nil {
			t.Fatalf("ValidateMatchConflicts after resolve: %v", err)
		}
	})

	t.Run("no_conflict_kept", func(t *testing.T) {
		profile := domain.Profile{
			ID:   "p",
			Name: "Work",
			Policies: []domain.ApplicationPolicy{
				{ProcessID: "m1", Action: domain.BlockAction()},
				{ProcessID: "h1", Action: domain.DirectAction()},
			},
		}
		got := domain.ResolveMatchConflictsLastWins(profile, identities)
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2", len(got))
		}
	})

	t.Run("same_fingerprint_keeps_both", func(t *testing.T) {
		profile := domain.Profile{
			ID:   "p",
			Name: "Work",
			Policies: []domain.ApplicationPolicy{
				{ProcessID: "h1", Action: domain.BlockAction()},
				{ProcessID: "h2", Action: domain.BlockAction()},
			},
		}
		got := domain.ResolveMatchConflictsLastWins(profile, identities)
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2", len(got))
		}
	})

	t.Run("unknown_process_kept", func(t *testing.T) {
		profile := domain.Profile{
			ID:   "p",
			Name: "Work",
			Policies: []domain.ApplicationPolicy{
				{ProcessID: "orphan", Action: domain.BlockAction()},
				{ProcessID: "h1", Action: domain.DirectAction()},
			},
		}
		got := domain.ResolveMatchConflictsLastWins(profile, identities)
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2", len(got))
		}
	})
}

func TestValidateMatchConflictsStillDetects(t *testing.T) {
	identities := []domain.ProcessIdentity{
		{ID: "h1", BundleID: "com.example.helper", ExecutableName: "A"},
		{ID: "h2", BundleID: "com.example.helper", ExecutableName: "B"},
	}
	err := domain.ValidateMatchConflicts(domain.Profile{
		Policies: []domain.ApplicationPolicy{
			{ProcessID: "h1", Action: domain.BlockAction()},
			{ProcessID: "h2", Action: domain.DirectAction()},
		},
	}, identities)
	if !errors.Is(err, domain.ErrConflictingMatch) {
		t.Fatalf("error = %v, want ErrConflictingMatch", err)
	}
}
