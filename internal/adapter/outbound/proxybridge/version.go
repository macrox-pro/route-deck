package proxybridge

import (
	"regexp"
	"strings"
)

var versionRE = regexp.MustCompile(`(?i)(\d+\.\d+(?:\.\d+)?)`)

// ParseVersion extracts a dotted version from --version output.
// Known means a version string was parsed. TestedVersion is 4.1.3.
func ParseVersion(output string) (version string, known bool, warning string) {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return "", false, ""
	}
	match := versionRE.FindString(trimmed)
	if match == "" {
		return "", false, ""
	}
	warning = ""
	if match != "4.1.3" {
		warning = "ProxyBridge version " + match + " is not the tested version 4.1.3"
	}
	return match, true, warning
}
