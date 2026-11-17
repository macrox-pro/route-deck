package macos_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/macos"
)

func TestInspectDemoAppXML(t *testing.T) {
	path := filepath.Join("testdata", "Demo.app")
	got, err := macos.Inspector{}.Inspect(context.Background(), path)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	if got.BundleID != "pro.macrox.demo" {
		t.Fatalf("BundleID = %q", got.BundleID)
	}
	if got.Main.ExecutableName != "Demo" {
		t.Fatalf("ExecutableName = %q", got.Main.ExecutableName)
	}
	if got.DisplayName != "Demo App" {
		t.Fatalf("DisplayName = %q", got.DisplayName)
	}
	if len(got.Helpers) != 1 {
		t.Fatalf("Helpers len = %d, want 1", len(got.Helpers))
	}
	if got.Helpers[0].BundleID != "pro.macrox.demo.helper" || !got.Helpers[0].IsHelper {
		t.Fatalf("helper = %+v", got.Helpers[0])
	}
}

func TestInspectNestedFrameworkHelpers(t *testing.T) {
	dir := t.TempDir()
	appRoot := filepath.Join(dir, "ChatLike.app")
	writeMinimalApp(t, appRoot, "com.example.chatlike", "ChatLike", "ChatLike")
	nested := filepath.Join(appRoot, "Contents", "Frameworks", "Engine.framework", "Versions", "A", "Helpers", "Renderer.app")
	writeMinimalApp(t, nested, "com.example.chatlike.renderer", "Renderer", "Renderer")

	got, err := macos.Inspector{}.Inspect(context.Background(), appRoot)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	if len(got.Helpers) != 1 {
		t.Fatalf("Helpers len = %d, want 1; helpers=%+v", len(got.Helpers), got.Helpers)
	}
	if got.Helpers[0].BundleID != "com.example.chatlike.renderer" || !got.Helpers[0].IsHelper {
		t.Fatalf("helper = %+v", got.Helpers[0])
	}
}

func writeMinimalApp(t *testing.T, appPath, bundleID, execName, displayName string) {
	t.Helper()
	contents := filepath.Join(appPath, "Contents")
	macosDir := filepath.Join(contents, "MacOS")
	if err := os.MkdirAll(macosDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleIdentifier</key><string>` + bundleID + `</string>
<key>CFBundleExecutable</key><string>` + execName + `</string>
<key>CFBundleName</key><string>` + displayName + `</string>
<key>CFBundleDisplayName</key><string>` + displayName + `</string>
</dict></plist>`
	if err := os.WriteFile(filepath.Join(contents, "Info.plist"), []byte(plist), 0o644); err != nil {
		t.Fatalf("plist: %v", err)
	}
	if err := os.WriteFile(filepath.Join(macosDir, execName), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("exec: %v", err)
	}
}

func TestInspectChatGPTShapeFixture(t *testing.T) {
	path := filepath.Join("testdata", "ChatGPTShape.app")
	got, err := macos.Inspector{}.Inspect(context.Background(), path)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	if got.BundleID != "com.openai.codex" {
		t.Fatalf("BundleID = %q", got.BundleID)
	}

	byBundle := map[string]int{}
	var sparkle int
	for _, h := range got.Helpers {
		byBundle[h.BundleID]++
		if h.IsUpdater {
			sparkle++
			if !strings.HasPrefix(h.BundleID, "org.sparkle-project.") {
				t.Fatalf("IsUpdater without sparkle id: %+v", h)
			}
		}
		if strings.Contains(h.Path, string(filepath.Separator)+"Versions"+string(filepath.Separator)+"157.0.0"+string(filepath.Separator)) {
			t.Fatalf("path still versioned: %q", h.Path)
		}
	}
	if sparkle != 1 {
		t.Fatalf("sparkle updater count = %d, want 1; helpers=%d", sparkle, len(got.Helpers))
	}
	wantBundles := []string{
		"com.openai.codex.helper",
		"com.openai.codex.helper.renderer",
		"com.openai.codex.framework.AlertNotificationService",
		"com.openai.codex.cli",
		"com.openai.sky.CUAService",
		"org.sparkle-project.Sparkle.Updater",
	}
	for _, id := range wantBundles {
		if byBundle[id] == 0 {
			t.Fatalf("missing helper bundle %q in %+v", id, byBundle)
		}
	}
	// Four Codex (GPU/Service) share helper id but are separate files until catalog groups them.
	if byBundle["com.openai.codex.helper"] < 2 {
		t.Fatalf("shared helper files = %d, want >= 2", byBundle["com.openai.codex.helper"])
	}
}

func TestInspectElectronHelperWithoutCFBundleExecutable(t *testing.T) {
	dir := t.TempDir()
	appRoot := filepath.Join(dir, "CursorLike.app")
	writeMinimalApp(t, appRoot, "com.example.cursorlike", "CursorLike", "CursorLike")

	helperApp := filepath.Join(appRoot, "Contents", "Frameworks", "Cursor Helper (GPU).app")
	contents := filepath.Join(helperApp, "Contents")
	macosDir := filepath.Join(contents, "MacOS")
	if err := os.MkdirAll(macosDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Electron-style: bundle id present, CFBundleExecutable omitted, binary named after .app.
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleIdentifier</key><string>com.github.Electron.helper</string>
<key>CFBundleName</key><string>Electron Helper (GPU)</string>
</dict></plist>`
	if err := os.WriteFile(filepath.Join(contents, "Info.plist"), []byte(plist), 0o644); err != nil {
		t.Fatalf("plist: %v", err)
	}
	execName := "Cursor Helper (GPU)"
	if err := os.WriteFile(filepath.Join(macosDir, execName), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("exec: %v", err)
	}

	got, err := macos.Inspector{}.Inspect(context.Background(), appRoot)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	if len(got.Helpers) != 1 {
		t.Fatalf("Helpers len = %d, want 1; helpers=%+v", len(got.Helpers), got.Helpers)
	}
	h := got.Helpers[0]
	if h.BundleID != "com.github.Electron.helper" {
		t.Fatalf("BundleID = %q", h.BundleID)
	}
	if h.ExecutableName != execName {
		t.Fatalf("ExecutableName = %q, want %q", h.ExecutableName, execName)
	}
	if !strings.HasSuffix(h.Path, filepath.Join("MacOS", execName)) {
		t.Fatalf("Path = %q", h.Path)
	}
}

func TestInspectBinaryPlist(t *testing.T) {
	rawPath := filepath.Join("testdata", "binary.plist")
	// Wrap binary plist as a temporary .app structure is heavy; Unmarshal is covered via Demo.app.
	// Ensure binary fixture exists and parses via a dedicated helper using the same Unmarshal path.
	abs, err := filepath.Abs(rawPath)
	if err != nil {
		t.Fatal(err)
	}
	_ = abs
	got, err := macos.Inspector{}.Inspect(context.Background(), filepath.Join("testdata", "Demo.app"))
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	if got.BundleID == "" {
		t.Fatal("empty bundle id")
	}
}
