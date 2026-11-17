// Package buildmeta holds link-time version identity for RouteDeck.
//
// Defaults apply to local and untagged builds. CI injects Version and Commit
// via -ldflags -X (see build/darwin/Taskfile.yml and .github/workflows/build-macos.yml).
package buildmeta

import "strings"

// Version is the marketing version without a required leading "v".
// Link-time default is "dev" when no tag is applied.
var Version = "dev"

// Commit is the short git SHA (typically 7 characters). Empty when unknown.
var Commit = ""

// DisplayVersion returns the short chrome label: "dev" or "v1.2.3".
func DisplayVersion() string {
	v := strings.TrimSpace(Version)
	if v == "" || v == "dev" {
		return "dev"
	}
	if strings.HasPrefix(v, "v") {
		return v
	}
	return "v" + v
}

// PlistVersion returns CFBundleShortVersionString / CFBundleVersion.
// Semver tags use the numeric form (no leading "v"). Untagged builds use "0.0.0".
func PlistVersion() string {
	v := strings.TrimSpace(Version)
	if v == "" || v == "dev" {
		return "0.0.0"
	}
	return strings.TrimPrefix(v, "v")
}
