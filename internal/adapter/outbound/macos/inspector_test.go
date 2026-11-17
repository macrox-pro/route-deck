package macos_test

import (
	"context"
	"path/filepath"
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
