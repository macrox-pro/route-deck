package proxybridge

import (
	"strings"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

// ClassifyLine maps one CLI output line to a BridgeLineClass.
func ClassifyLine(line string) port.BridgeLineClass {
	s := strings.TrimSpace(line)
	switch {
	case strings.Contains(s, "Redirecting traffic. Press Ctrl+C to stop."):
		return port.BridgeLineRedirectReady
	case strings.Contains(s, "proxies,") && strings.Contains(s, "rules applied"):
		return port.BridgeLineProfileLoaded
	case strings.Contains(s, "The ProxyBridge app is running"):
		return port.BridgeLineGUIConflict
	case strings.Contains(s, "Another ProxyBridge CLI is already running"):
		return port.BridgeLineAnotherCLI
	case strings.Contains(s, "The tunnel was stopped by something else"):
		return port.BridgeLineTunnelStopped
	case strings.Contains(s, "Timed out waiting for the extension"):
		return port.BridgeLineExtensionTimeout
	case strings.Contains(s, "rule(s) were rejected by the extension"):
		return port.BridgeLineRejectedRules
	default:
		return port.BridgeLineUnclassified
	}
}
