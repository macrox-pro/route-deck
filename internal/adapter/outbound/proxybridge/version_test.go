package proxybridge_test

import (
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/proxybridge"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		wantVersion string
		wantKnown   bool
		wantWarn    bool
	}{
		{name: "tested", in: "ProxyBridge CLI 4.1.3", wantVersion: "4.1.3", wantKnown: true},
		{name: "other", in: "ProxyBridge CLI 4.2.0", wantVersion: "4.2.0", wantKnown: true, wantWarn: true},
		{name: "empty", in: "", wantKnown: false},
		{name: "garbage", in: "not a version", wantKnown: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, known, warn := proxybridge.ParseVersion(tt.in)
			if v != tt.wantVersion || known != tt.wantKnown {
				t.Fatalf("got (%q,%v), want (%q,%v)", v, known, tt.wantVersion, tt.wantKnown)
			}
			if tt.wantWarn && warn == "" {
				t.Fatal("expected warning")
			}
			if !tt.wantWarn && warn != "" {
				t.Fatalf("unexpected warning %q", warn)
			}
		})
	}
}
