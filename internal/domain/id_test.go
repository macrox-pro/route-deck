package domain_test

import (
	"errors"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestParseIDs(t *testing.T) {
	tests := []struct {
		name    string
		parse   func(string) (string, error)
		in      string
		want    string
		wantErr error
	}{
		{
			name: "proxy_id_empty",
			parse: func(s string) (string, error) {
				id, err := domain.ParseProxyID(s)
				return string(id), err
			},
			in:      "",
			wantErr: domain.ErrEmptyID,
		},
		{
			name: "proxy_id_valid",
			parse: func(s string) (string, error) {
				id, err := domain.ParseProxyID(s)
				return string(id), err
			},
			in:   "proxy-1",
			want: "proxy-1",
		},
		{
			name: "application_id_empty",
			parse: func(s string) (string, error) {
				id, err := domain.ParseApplicationID(s)
				return string(id), err
			},
			in:      "",
			wantErr: domain.ErrEmptyID,
		},
		{
			name: "application_id_valid",
			parse: func(s string) (string, error) {
				id, err := domain.ParseApplicationID(s)
				return string(id), err
			},
			in:   "app-1",
			want: "app-1",
		},
		{
			name: "process_identity_id_empty",
			parse: func(s string) (string, error) {
				id, err := domain.ParseProcessIdentityID(s)
				return string(id), err
			},
			in:      "",
			wantErr: domain.ErrEmptyID,
		},
		{
			name: "process_identity_id_valid",
			parse: func(s string) (string, error) {
				id, err := domain.ParseProcessIdentityID(s)
				return string(id), err
			},
			in:   "proc-1",
			want: "proc-1",
		},
		{
			name: "profile_id_empty",
			parse: func(s string) (string, error) {
				id, err := domain.ParseProfileID(s)
				return string(id), err
			},
			in:      "",
			wantErr: domain.ErrEmptyID,
		},
		{
			name: "profile_id_valid",
			parse: func(s string) (string, error) {
				id, err := domain.ParseProfileID(s)
				return string(id), err
			},
			in:   "profile-1",
			want: "profile-1",
		},
		{
			name: "keychain_ref_empty",
			parse: func(s string) (string, error) {
				id, err := domain.ParseKeychainRef(s)
				return string(id), err
			},
			in:      "",
			wantErr: domain.ErrEmptyID,
		},
		{
			name: "keychain_ref_valid",
			parse: func(s string) (string, error) {
				id, err := domain.ParseKeychainRef(s)
				return string(id), err
			},
			in:   "kc-ref-1",
			want: "kc-ref-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parse(tt.in)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("parse(%q) error = %v", tt.in, err)
				}
				if got != tt.want {
					t.Fatalf("parse(%q) = %q, want %q", tt.in, got, tt.want)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("parse(%q) error = %v, want %v", tt.in, err, tt.wantErr)
			}
			if got != "" {
				t.Fatalf("parse(%q) value = %q, want empty on error", tt.in, got)
			}
		})
	}
}
