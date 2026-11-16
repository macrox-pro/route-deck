package domain_test

import (
	"errors"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestNewProfile(t *testing.T) {
	tests := []struct {
		name     string
		id       domain.ProfileID
		profName string
		policies []domain.ApplicationPolicy
		wantErr  error
		want     domain.Profile
	}{
		{
			name:     "empty_name",
			id:       "profile-1",
			profName: "",
			policies: nil,
			wantErr:  domain.ErrEmptyName,
		},
		{
			name:     "whitespace_name_and_empty_id",
			id:       "",
			profName: "   ",
			policies: nil,
			want: domain.Profile{
				ID:       "",
				Name:     "   ",
				Policies: nil,
			},
		},
		{
			name:     "valid_with_policies",
			id:       "profile-2",
			profName: "Work",
			policies: []domain.ApplicationPolicy{
				{ProcessID: "proc-1", Action: domain.DirectAction()},
			},
			want: domain.Profile{
				ID:   "profile-2",
				Name: "Work",
				Policies: []domain.ApplicationPolicy{
					{ProcessID: "proc-1", Action: domain.DirectAction()},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewProfile(tt.id, tt.profName, tt.policies)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewProfile() error = %v, want %v", err, tt.wantErr)
				}
				if got.ID != "" || got.Name != "" || got.Policies != nil {
					t.Fatalf("NewProfile() on error = %#v, want zero value", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewProfile() error = %v", err)
			}
			if got.ID != tt.want.ID || got.Name != tt.want.Name {
				t.Fatalf("NewProfile() = %#v, want %#v", got, tt.want)
			}
			if len(got.Policies) != len(tt.want.Policies) {
				t.Fatalf("NewProfile() Policies len = %d, want %d", len(got.Policies), len(tt.want.Policies))
			}
			for i := range got.Policies {
				if got.Policies[i] != tt.want.Policies[i] {
					t.Fatalf("NewProfile() Policies[%d] = %#v, want %#v", i, got.Policies[i], tt.want.Policies[i])
				}
			}
		})
	}
}

func TestValidateProxyRefs(t *testing.T) {
	knownProxy, err := domain.ProxyAction("proxy-known")
	if err != nil {
		t.Fatalf("ProxyAction(proxy-known): %v", err)
	}
	missingProxy, err := domain.ProxyAction("proxy-missing")
	if err != nil {
		t.Fatalf("ProxyAction(proxy-missing): %v", err)
	}
	direct := domain.DirectAction()

	tests := []struct {
		name    string
		profile domain.Profile
		known   []domain.ProxyID
		wantErr error
	}{
		{
			name: "known_id",
			profile: domain.Profile{
				ID:   "p1",
				Name: "Known",
				Policies: []domain.ApplicationPolicy{
					{ProcessID: "proc-1", Action: knownProxy},
					{ProcessID: "proc-2", Action: direct},
				},
			},
			known:   []domain.ProxyID{"proxy-known"},
			wantErr: nil,
		},
		{
			name: "missing_id",
			profile: domain.Profile{
				ID:   "p2",
				Name: "Missing",
				Policies: []domain.ApplicationPolicy{
					{ProcessID: "proc-1", Action: knownProxy},
					{ProcessID: "proc-2", Action: missingProxy},
				},
			},
			known:   []domain.ProxyID{"proxy-known"},
			wantErr: domain.ErrUnresolvedProxy,
		},
		{
			name: "direct_is_not_a_proxy_ref",
			profile: domain.Profile{
				ID:   "p3",
				Name: "DirectOnly",
				Policies: []domain.ApplicationPolicy{
					{ProcessID: "proc-1", Action: direct},
				},
			},
			known:   nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := make([]domain.RoutingAction, len(tt.profile.Policies))
			for i, p := range tt.profile.Policies {
				before[i] = p.Action
			}

			err := domain.ValidateProxyRefs(tt.profile, tt.known)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("ValidateProxyRefs() error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ValidateProxyRefs() error = %v, want %v", err, tt.wantErr)
			}

			for i, p := range tt.profile.Policies {
				if p.Action != before[i] {
					t.Fatalf("ValidateProxyRefs() Action[%d] = %#v, want %#v (unchanged)", i, p.Action, before[i])
				}
				if tt.name == "missing_id" && p.Action.IsDirect() {
					t.Fatalf("ValidateProxyRefs() Action[%d] became Direct; must stay Proxy", i)
				}
			}
		})
	}
}
