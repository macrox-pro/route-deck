package logfile

import (
	"log/slog"
	"testing"
)

func TestRedactAttr(t *testing.T) {
	tests := []struct {
		name string
		attr slog.Attr
		want string
	}{
		{name: "password", attr: slog.String("password", "secret-value"), want: "[REDACTED]"},
		{name: "Password_case", attr: slog.String("Password", "x"), want: "[REDACTED]"},
		{name: "secret", attr: slog.String("secret", "x"), want: "[REDACTED]"},
		{name: "token", attr: slog.String("token", "x"), want: "[REDACTED]"},
		{name: "profile", attr: slog.String("profile", `{"id":"1"}`), want: "[REDACTED]"},
		{name: "proxy_id_kept", attr: slog.String("proxy_id", "abc"), want: "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RedactAttr(nil, tt.attr)
			if got.Value.String() != tt.want {
				t.Fatalf("RedactAttr() = %q, want %q", got.Value.String(), tt.want)
			}
		})
	}
}
