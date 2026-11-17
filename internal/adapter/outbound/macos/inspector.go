// Package macos implements ApplicationInspector for .app bundles.
package macos

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

// Inspector reads Info.plist metadata from macOS application bundles.
type Inspector struct{}

var _ port.ApplicationInspector = Inspector{}

type infoPlist struct {
	CFBundleIdentifier  string `plist:"CFBundleIdentifier"`
	CFBundleExecutable  string `plist:"CFBundleExecutable"`
	CFBundleName        string `plist:"CFBundleName"`
	CFBundleDisplayName string `plist:"CFBundleDisplayName"`
}

// Inspect reads path as a .app and nested helper apps (ADR-0014).
func (Inspector) Inspect(ctx context.Context, path string) (port.InspectedApplication, error) {
	if err := ctx.Err(); err != nil {
		return port.InspectedApplication{}, fmt.Errorf("macos inspect: %w", err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return port.InspectedApplication{}, fmt.Errorf("macos inspect: %w", err)
	}
	main, err := readBundle(abs, false)
	if err != nil {
		return port.InspectedApplication{}, err
	}
	helpers, err := findHelpers(abs)
	if err != nil {
		return port.InspectedApplication{}, err
	}
	display := main.DisplayName
	if display == "" {
		display = main.ExecutableName
	}
	return port.InspectedApplication{
		DisplayName: display,
		BundleID:    main.BundleID,
		Main:        main,
		Helpers:     helpers,
	}, nil
}

func readBundle(appPath string, isHelper bool) (port.InspectedProcess, error) {
	info, err := readInfoPlist(appPath)
	if err != nil {
		return port.InspectedProcess{}, err
	}
	if info.CFBundleIdentifier == "" {
		return port.InspectedProcess{}, fmt.Errorf("macos inspect: missing CFBundleIdentifier")
	}
	display := displayNameFromInfo(info)
	bundleRoot := preferCurrentVersion(appPath)
	execName := info.CFBundleExecutable
	if execName == "" {
		// Electron helpers (Cursor, VS Code, …) often omit CFBundleExecutable
		// but still ship a binary under Contents/MacOS (ADR-0014: .app with an executable).
		execName = discoverMacOSExecutable(bundleRoot, appPath)
	}
	if execName == "" && isHelper {
		return port.InspectedProcess{}, fmt.Errorf("macos inspect: missing executable")
	}
	execPath := ""
	if execName != "" {
		execPath = filepath.Join(bundleRoot, "Contents", "MacOS", execName)
	}
	realPath := appPath
	if resolved, err := filepath.EvalSymlinks(appPath); err == nil {
		realPath = resolved
	}
	return port.InspectedProcess{
		BundleID:       info.CFBundleIdentifier,
		ExecutableName: execName,
		Path:           execPath,
		DisplayName:    display,
		IsHelper:       isHelper,
		IsUpdater:      isSparkleUpdater(info.CFBundleIdentifier),
		RealPath:       realPath,
	}, nil
}

// discoverMacOSExecutable finds Contents/MacOS binary when Info.plist omits CFBundleExecutable.
func discoverMacOSExecutable(bundleRoot, appPath string) string {
	macosDir := filepath.Join(bundleRoot, "Contents", "MacOS")
	entries, err := os.ReadDir(macosDir)
	if err != nil {
		return ""
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// Skip AppleDouble / resource forks.
		if strings.HasPrefix(entry.Name(), "._") {
			continue
		}
		names = append(names, entry.Name())
	}
	if len(names) == 0 {
		return ""
	}
	if len(names) == 1 {
		return names[0]
	}
	base := strings.TrimSuffix(filepath.Base(appPath), filepath.Ext(appPath))
	for _, name := range names {
		if name == base {
			return name
		}
	}
	return names[0]
}

func isSparkleUpdater(bundleID string) bool {
	return strings.HasPrefix(bundleID, "org.sparkle-project.")
}

// preferCurrentVersion rewrites .../Versions/<ver>/... to .../Versions/Current/... when Current exists.
func preferCurrentVersion(appPath string) string {
	const marker = string(filepath.Separator) + "Versions" + string(filepath.Separator)
	idx := strings.Index(appPath, marker)
	if idx < 0 {
		return appPath
	}
	after := appPath[idx+len(marker):]
	parts := strings.SplitN(after, string(filepath.Separator), 2)
	if len(parts) == 0 || parts[0] == "" || parts[0] == "Current" {
		return appPath
	}
	frameworkRoot := appPath[:idx]
	currentLink := filepath.Join(frameworkRoot, "Versions", "Current")
	if _, err := os.Lstat(currentLink); err != nil {
		return appPath
	}
	rest := ""
	if len(parts) == 2 {
		rest = parts[1]
	}
	if rest == "" {
		return currentLink
	}
	return filepath.Join(currentLink, rest)
}

func findHelpers(appPath string) ([]port.InspectedProcess, error) {
	roots := []string{
		filepath.Join(appPath, "Contents", "Frameworks"),
		filepath.Join(appPath, "Contents", "Helpers"),
		filepath.Join(appPath, "Contents", "Resources"),
	}
	var helpers []port.InspectedProcess
	seen := make(map[string]struct{})
	for _, root := range roots {
		if err := walkHelperApps(root, &helpers, seen); err != nil {
			return nil, err
		}
	}
	return helpers, nil
}

// walkHelperApps collects nested .app bundles (ADR-0014).
// A found .app is a leaf: we do not walk inside it for further helpers.
func walkHelperApps(root string, helpers *[]port.InspectedProcess, seen map[string]struct{}) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("macos inspect helpers: %w", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(root, name)
		if strings.HasSuffix(strings.ToLower(name), ".app") {
			proc, err := readBundle(path, true)
			if err != nil {
				continue
			}
			key := proc.RealPath
			if key == "" {
				key = path
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			*helpers = append(*helpers, proc)
			continue
		}
		if !entry.IsDir() {
			continue
		}
		if err := walkHelperApps(path, helpers, seen); err != nil {
			return err
		}
	}
	return nil
}
