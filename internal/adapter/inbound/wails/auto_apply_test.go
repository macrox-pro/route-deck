package wails_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

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

type recordingDialog struct {
	warnings []wails.UserMsg
}

func (d *recordingDialog) Info(wails.UserMsg) {}
func (d *recordingDialog) Warning(msg wails.UserMsg) {
	d.warnings = append(d.warnings, msg)
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
			var calls atomic.Int32
			wails.MaybeAutoApplyAtStartup(context.Background(), wails.AutoApplyConfig{
				Logger:   logger,
				Settings: autoApplySettingsStore{settings: tc.settings},
				Profiles: autoApplyProfileStore{profile: tc.profile},
				Apply: func(context.Context, domain.ProfileID) error {
					calls.Add(1)
					return nil
				},
			})
			got := calls.Load() > 0
			if got != tc.wantCall {
				t.Fatalf("apply called = %v, want %v", got, tc.wantCall)
			}
		})
	}
}

func TestShouldRetryAutoApply(t *testing.T) {
	t.Parallel()
	proxy, err := domain.ProxyAction("p1")
	if err != nil {
		t.Fatal(err)
	}
	policy, err := domain.NewApplicationPolicy("proc", proxy, domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatal(err)
	}
	profile, err := domain.NewProfile("prof-1", "Work", []domain.ApplicationPolicy{policy})
	if err != nil {
		t.Fatal(err)
	}
	healthUnavailable := map[domain.ProxyID]domain.ProxyHealth{
		"p1": domain.ProxyHealthUnavailable,
	}
	healthAuth := map[domain.ProxyID]domain.ProxyHealth{
		"p1": domain.ProxyHealthAuthenticationFailed,
	}

	cases := []struct {
		name   string
		err    error
		health map[domain.ProxyID]domain.ProxyHealth
		want   bool
	}{
		{name: "nil_error", err: nil, health: healthUnavailable, want: false},
		{name: "unavailable", err: errors.New("dial refused"), health: healthUnavailable, want: true},
		{name: "auth_health", err: errors.New("dial refused"), health: healthAuth, want: false},
		{name: "conflict", err: errors.New(string(port.BridgeLineGUIConflict)), health: healthUnavailable, want: false},
		{name: "another_cli", err: errors.New(string(port.BridgeLineAnotherCLI)), health: healthUnavailable, want: false},
		{name: "extension_timeout", err: errors.New(string(port.BridgeLineExtensionTimeout)), health: healthUnavailable, want: false},
		{name: "not_installed", err: errors.New("protection apply: proxybridge not installed"), health: healthUnavailable, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := wails.ShouldRetryAutoApply(tc.err, profile, tc.health)
			if got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMaybeAutoApplyRetriesOnceThenWarns(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	proxy, err := domain.ProxyAction("p1")
	if err != nil {
		t.Fatal(err)
	}
	policy, err := domain.NewApplicationPolicy("proc", proxy, domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatal(err)
	}
	profile, err := domain.NewProfile("prof-1", "Work", []domain.ApplicationPolicy{policy})
	if err != nil {
		t.Fatal(err)
	}
	dialog := &recordingDialog{}
	var calls atomic.Int32
	wails.MaybeAutoApplyAtStartup(context.Background(), wails.AutoApplyConfig{
		Logger: logger,
		Settings: autoApplySettingsStore{settings: port.Settings{
			AutoApply: true, DesiredSet: true, DesiredProfileID: "prof-1",
		}},
		Profiles: autoApplyProfileStore{profile: profile},
		HealthSnapshot: func() map[domain.ProxyID]domain.ProxyHealth {
			return map[domain.ProxyID]domain.ProxyHealth{"p1": domain.ProxyHealthUnavailable}
		},
		Apply: func(context.Context, domain.ProfileID) error {
			calls.Add(1)
			return errors.New("start failed")
		},
		Dialog:     dialog,
		RetryDelay: time.Millisecond,
	})
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want 2", calls.Load())
	}
	if len(dialog.warnings) != 1 || dialog.warnings[0].ID != "protection.auto_apply_failed" {
		t.Fatalf("warnings = %+v", dialog.warnings)
	}
}

func TestMaybeAutoApplyNoRetryOnConflict(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	proxy, err := domain.ProxyAction("p1")
	if err != nil {
		t.Fatal(err)
	}
	policy, err := domain.NewApplicationPolicy("proc", proxy, domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatal(err)
	}
	profile, err := domain.NewProfile("prof-1", "Work", []domain.ApplicationPolicy{policy})
	if err != nil {
		t.Fatal(err)
	}
	dialog := &recordingDialog{}
	var calls atomic.Int32
	wails.MaybeAutoApplyAtStartup(context.Background(), wails.AutoApplyConfig{
		Logger: logger,
		Settings: autoApplySettingsStore{settings: port.Settings{
			AutoApply: true, DesiredSet: true, DesiredProfileID: "prof-1",
		}},
		Profiles: autoApplyProfileStore{profile: profile},
		HealthSnapshot: func() map[domain.ProxyID]domain.ProxyHealth {
			return map[domain.ProxyID]domain.ProxyHealth{"p1": domain.ProxyHealthUnavailable}
		},
		Apply: func(context.Context, domain.ProfileID) error {
			calls.Add(1)
			return errors.New(string(port.BridgeLineGUIConflict))
		},
		Dialog:     dialog,
		RetryDelay: time.Millisecond,
	})
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1", calls.Load())
	}
	if len(dialog.warnings) != 1 || dialog.warnings[0].ID != "protection.auto_apply_failed" {
		t.Fatalf("warnings = %+v", dialog.warnings)
	}
}
