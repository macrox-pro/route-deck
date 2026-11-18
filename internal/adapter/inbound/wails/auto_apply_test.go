package wails_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/wails"
	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type autoApplySettingsStore struct {
	settings port.Settings
	err      error
}

func (s autoApplySettingsStore) Get(context.Context) (port.Settings, error) {
	if s.err != nil {
		return port.Settings{}, s.err
	}
	return s.settings, nil
}

func (s autoApplySettingsStore) Save(context.Context, port.Settings) error {
	return nil
}

type autoApplyProfileStore struct {
	profile domain.Profile
	err     error
}

func (s autoApplyProfileStore) Get(_ context.Context, id domain.ProfileID) (domain.Profile, error) {
	if s.err != nil {
		return domain.Profile{}, s.err
	}
	if s.profile.ID != "" && s.profile.ID != id {
		return domain.Profile{}, errors.New("not found")
	}
	return s.profile, nil
}

func (s autoApplyProfileStore) List(context.Context) ([]domain.Profile, error) { return nil, nil }
func (s autoApplyProfileStore) Save(context.Context, domain.Profile) error     { return nil }
func (s autoApplyProfileStore) Delete(context.Context, domain.ProfileID) error {
	return nil
}

func TestMaybeAutoApplyAtStartup(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	proxy, err := domain.ProxyAction("p1")
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}
	policy, err := domain.NewApplicationPolicy("proc", proxy, domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatalf("policy: %v", err)
	}
	profile, err := domain.NewProfile("prof-1", "Work", []domain.ApplicationPolicy{policy})
	if err != nil {
		t.Fatalf("profile: %v", err)
	}
	directPolicy, err := domain.NewApplicationPolicy("proc", domain.DirectAction(), domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatalf("direct: %v", err)
	}
	directProfile, err := domain.NewProfile("prof-1", "Direct", []domain.ApplicationPolicy{directPolicy})
	if err != nil {
		t.Fatalf("direct profile: %v", err)
	}

	cases := []struct {
		name     string
		settings port.Settings
		profile  domain.Profile
		wantCall bool
	}{
		{
			name: "gates_pass",
			settings: port.Settings{
				AutoApply:        true,
				DesiredSet:       true,
				DesiredProfileID: "prof-1",
			},
			profile:  profile,
			wantCall: true,
		},
		{
			name: "disabled",
			settings: port.Settings{
				AutoApply:        false,
				DesiredSet:       true,
				DesiredProfileID: "prof-1",
			},
			profile:  profile,
			wantCall: false,
		},
		{
			name: "direct_only",
			settings: port.Settings{
				AutoApply:        true,
				DesiredSet:       true,
				DesiredProfileID: "prof-1",
			},
			profile:  directProfile,
			wantCall: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			called := false
			wails.MaybeAutoApplyAtStartup(
				context.Background(),
				logger,
				autoApplySettingsStore{settings: tc.settings},
				autoApplyProfileStore{profile: tc.profile},
				func(context.Context) (wails.ApplyResultDTO, error) {
					called = true
					return wails.ApplyResultDTO{OperationID: "op"}, nil
				},
			)
			if called != tc.wantCall {
				t.Fatalf("apply called = %v, want %v", called, tc.wantCall)
			}
		})
	}
}
