package proxybridge_test

import (
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/proxybridge"
	"github.com/macrox-pro/route-deck/internal/application/port"
)

func TestClassifyLine(t *testing.T) {
	tests := []struct {
		line string
		want port.BridgeLineClass
	}{
		{"Redirecting traffic. Press Ctrl+C to stop.", port.BridgeLineRedirectReady},
		{`Profile "Work": 1 proxies, 2 rules applied`, port.BridgeLineProfileLoaded},
		{"The ProxyBridge app is running...", port.BridgeLineGUIConflict},
		{"Another ProxyBridge CLI is already running.", port.BridgeLineAnotherCLI},
		{"The tunnel was stopped by something else...", port.BridgeLineTunnelStopped},
		{"Timed out waiting for the extension...", port.BridgeLineExtensionTimeout},
		{"warning: 2 rule(s) were rejected by the extension (invalid fields)", port.BridgeLineRejectedRules},
		{"random noise", port.BridgeLineUnclassified},
	}
	for _, tt := range tests {
		if got := proxybridge.ClassifyLine(tt.line); got != tt.want {
			t.Fatalf("ClassifyLine(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}
}
