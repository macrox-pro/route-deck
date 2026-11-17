package macos

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"howett.net/plist"
)

func writeApp(t *testing.T, dir, name, bundleID, display string) string {
	t.Helper()
	appPath := filepath.Join(dir, name+".app")
	contents := filepath.Join(appPath, "Contents")
	if err := os.MkdirAll(contents, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	info := map[string]any{
		"CFBundleIdentifier":  bundleID,
		"CFBundleName":        display,
		"CFBundleDisplayName": display,
		"CFBundleExecutable":  name,
	}
	raw, err := plist.Marshal(info, plist.XMLFormat)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(contents, "Info.plist"), raw, 0o644); err != nil {
		t.Fatalf("write plist: %v", err)
	}
	return appPath
}

func TestListerDedupeAndNested(t *testing.T) {
	root := t.TempDir()
	apps := filepath.Join(root, "Applications")
	homeApps := filepath.Join(root, "HomeApplications")
	sysApps := filepath.Join(root, "SystemApplications")
	for _, dir := range []string{apps, homeApps, sysApps} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	homeSafari := writeApp(t, homeApps, "Safari", "com.apple.Safari", "Safari Home")
	_ = writeApp(t, apps, "Safari", "com.apple.Safari", "Safari Apps")
	_ = writeApp(t, sysApps, "Preview", "com.apple.Preview", "Preview")
	nestedParent := writeApp(t, apps, "Parent", "com.example.Parent", "Parent")
	nested := filepath.Join(nestedParent, "Contents", "Frameworks")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	_ = writeApp(t, nested, "Helper", "com.example.Helper", "Helper")

	missingRoot := filepath.Join(root, "MissingApplications")
	lister := &Lister{
		Roots:           []string{apps, homeApps, sysApps, missingRoot},
		HomeDir:         homeApps, // synthetic "~/Applications" for priority
		UserFacingRoots: []string{apps, homeApps},
		IndexPaths: func(context.Context) ([]string, error) {
			return []string{
				filepath.Join(nested, "Helper.app"), // nested from index must drop
			}, nil
		},
	}
	got, err := lister.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	byID := map[string]struct {
		path     string
		isSystem bool
	}{}
	for _, app := range got {
		byID[app.BundleID] = struct {
			path     string
			isSystem bool
		}{path: app.Path, isSystem: app.IsSystem}
		if app.BundleID == "com.example.Helper" {
			t.Fatalf("nested helper listed: %+v", app)
		}
	}
	if byID["com.apple.Safari"].path != homeSafari {
		t.Fatalf("safari path = %q, want home %q", byID["com.apple.Safari"].path, homeSafari)
	}
	if byID["com.apple.Safari"].isSystem {
		t.Fatal("safari under user-facing root must not be system")
	}
	if _, ok := byID["com.apple.Preview"]; !ok {
		t.Fatal("preview missing")
	}
	if !byID["com.apple.Preview"].isSystem {
		t.Fatal("preview under system root must be system")
	}
	if _, ok := byID["com.example.Parent"]; !ok {
		t.Fatal("parent missing")
	}
	if byID["com.example.Parent"].isSystem {
		t.Fatal("parent under /Applications stand-in must not be system")
	}
}

func TestListerSkipsBadPlist(t *testing.T) {
	root := t.TempDir()
	apps := filepath.Join(root, "Applications")
	if err := os.MkdirAll(apps, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	bad := filepath.Join(apps, "Bad.app", "Contents")
	if err := os.MkdirAll(bad, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bad, "Info.plist"), []byte("not-plist"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	good := writeApp(t, apps, "Good", "com.example.Good", "Good")
	lister := &Lister{
		Roots:      []string{apps},
		IndexPaths: func(context.Context) ([]string, error) { return nil, nil },
	}
	got, err := lister.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 || got[0].Path != good {
		t.Fatalf("got %#v", got)
	}
}

func TestListerIndexFailureFallsBackToRoots(t *testing.T) {
	root := t.TempDir()
	apps := filepath.Join(root, "Applications")
	if err := os.MkdirAll(apps, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	want := writeApp(t, apps, "Only", "com.example.Only", "Only")
	lister := &Lister{
		Roots: []string{apps},
		IndexPaths: func(context.Context) ([]string, error) {
			return nil, os.ErrPermission
		},
	}
	got, err := lister.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 || got[0].Path != want {
		t.Fatalf("got %#v", got)
	}
}
