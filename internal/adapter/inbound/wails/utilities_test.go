package wails_test

import (
	"strings"
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/wails"
)

func TestUtilityWindowsService_OpenGuidesStepValidation(t *testing.T) {
	t.Parallel()
	svc := &wails.UtilityWindowsService{Windows: &wails.UtilityWindows{}}

	cases := []struct {
		step    string
		wantErr bool
	}{
		{step: "", wantErr: false},
		{step: "bridge", wantErr: false},
		{step: "proxy", wantErr: false},
		{step: "application", wantErr: false},
		{step: "profile", wantErr: false},
		{step: "nope", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.step+"_", func(t *testing.T) {
			t.Parallel()
			err := svc.OpenGuides(tc.step)
			// Without App configured, valid steps still fail with "not configured" once past validation,
			// or succeed validation then hit Open. Invalid steps must fail with invalid step.
			if tc.wantErr {
				if err == nil || !strings.Contains(err.Error(), "invalid guides step") {
					t.Fatalf("OpenGuides(%q) error = %v, want invalid guides step", tc.step, err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected not configured error with empty App")
			}
			if strings.Contains(err.Error(), "invalid guides step") {
				t.Fatalf("valid step %q rejected: %v", tc.step, err)
			}
		})
	}
}

func TestUtilityWindowsService_ShowMainRouteValidation(t *testing.T) {
	t.Parallel()
	svc := &wails.UtilityWindowsService{Windows: &wails.UtilityWindows{}}

	cases := []struct {
		route   string
		wantErr string
	}{
		{route: "/", wantErr: ""},
		{route: "/settings", wantErr: ""},
		{route: "/help", wantErr: "invalid main route"},
		{route: "", wantErr: "invalid main route"},
	}
	for _, tc := range cases {
		t.Run(tc.route, func(t *testing.T) {
			t.Parallel()
			err := svc.ShowMain(tc.route)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("ShowMain(%q) error = %v, want %q", tc.route, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ShowMain(%q) unexpected error: %v", tc.route, err)
			}
		})
	}
}
