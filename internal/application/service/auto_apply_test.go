package service_test

import (
	"testing"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestProfileHasProtectableRules(t *testing.T) {
	t.Parallel()
	proxy, err := domain.ProxyAction("p1")
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}
	proxyPolicy, err := domain.NewApplicationPolicy("proc", proxy, domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatalf("proxy policy: %v", err)
	}
	blockPolicy, err := domain.NewApplicationPolicy("proc", domain.BlockAction(), domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatalf("block policy: %v", err)
	}
	directPolicy, err := domain.NewApplicationPolicy("proc", domain.DirectAction(), domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatalf("direct policy: %v", err)
	}

	cases := []struct {
		name     string
		policies []domain.ApplicationPolicy
		want     bool
	}{
		{name: "empty", want: false},
		{name: "direct_only", policies: []domain.ApplicationPolicy{directPolicy}, want: false},
		{name: "proxy", policies: []domain.ApplicationPolicy{proxyPolicy}, want: true},
		{name: "block", policies: []domain.ApplicationPolicy{blockPolicy}, want: true},
		{name: "direct_and_proxy", policies: []domain.ApplicationPolicy{directPolicy, proxyPolicy}, want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			profile, err := domain.NewProfile("prof", "Test", tc.policies)
			if err != nil {
				t.Fatalf("NewProfile: %v", err)
			}
			if got := service.ProfileHasProtectableRules(profile); got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestShouldAutoApplyOnStartup(t *testing.T) {
	t.Parallel()
	proxy, err := domain.ProxyAction("p1")
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}
	policy, err := domain.NewApplicationPolicy("proc", proxy, domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatalf("policy: %v", err)
	}
	protectable, err := domain.NewProfile("prof-1", "Work", []domain.ApplicationPolicy{policy})
	if err != nil {
		t.Fatalf("profile: %v", err)
	}
	directPolicy, err := domain.NewApplicationPolicy("proc", domain.DirectAction(), domain.ProtocolScopeBOTH, "")
	if err != nil {
		t.Fatalf("direct: %v", err)
	}
	directOnly, err := domain.NewProfile("prof-1", "Direct", []domain.ApplicationPolicy{directPolicy})
	if err != nil {
		t.Fatalf("direct profile: %v", err)
	}

	cases := []struct {
		name       string
		autoApply  bool
		desiredSet bool
		desiredID  domain.ProfileID
		profile    domain.Profile
		want       bool
	}{
		{name: "all_gates", autoApply: true, desiredSet: true, desiredID: "prof-1", profile: protectable, want: true},
		{name: "flag_off", autoApply: false, desiredSet: true, desiredID: "prof-1", profile: protectable, want: false},
		{name: "no_desired", autoApply: true, desiredSet: false, desiredID: "", profile: protectable, want: false},
		{name: "empty_desired_id", autoApply: true, desiredSet: true, desiredID: "", profile: protectable, want: false},
		{name: "direct_only", autoApply: true, desiredSet: true, desiredID: "prof-1", profile: directOnly, want: false},
		{name: "id_mismatch", autoApply: true, desiredSet: true, desiredID: "other", profile: protectable, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := service.ShouldAutoApplyOnStartup(tc.autoApply, tc.desiredSet, tc.desiredID, tc.profile)
			if got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}
