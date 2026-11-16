package domain_test

import (
	"reflect"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestProcessIdentityMatchableNames(t *testing.T) {
	tests := []struct {
		name string
		in   domain.ProcessIdentity
		want []string
	}{
		{
			name: "both_names",
			in: domain.ProcessIdentity{
				BundleID:       "com.example.app",
				ExecutableName: "Example",
				Path:           "/Applications/Example.app",
				SigningHint:    "Developer ID",
			},
			want: []string{"com.example.app", "Example"},
		},
		{
			name: "empty_bundle_id",
			in: domain.ProcessIdentity{
				ExecutableName: "Example",
				Path:           "/Applications/Example.app",
				SigningHint:    "hint",
			},
			want: []string{"Example"},
		},
		{
			name: "empty_executable_name",
			in: domain.ProcessIdentity{
				BundleID:    "com.example.app",
				Path:        "/Applications/Example.app",
				SigningHint: "hint",
			},
			want: []string{"com.example.app"},
		},
		{
			name: "slash_in_executable_name",
			in: domain.ProcessIdentity{
				BundleID:       "com.example.app",
				ExecutableName: "Helpers/ExampleHelper",
				Path:           "/Applications/Example.app/Contents/MacOS/Example",
				SigningHint:    "hint",
			},
			want: []string{"com.example.app"},
		},
		{
			name: "path_and_signing_hint_absent",
			in: domain.ProcessIdentity{
				BundleID:       "com.example.app",
				ExecutableName: "Example",
				Path:           "/must/not/appear",
				SigningHint:    "must-not-appear",
			},
			want: []string{"com.example.app", "Example"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.MatchableNames()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("MatchableNames() = %#v, want %#v", got, tt.want)
			}
			for _, name := range got {
				if name == tt.in.Path || name == tt.in.SigningHint {
					t.Fatalf("MatchableNames() included Path or SigningHint: %#v", got)
				}
			}
		})
	}
}
