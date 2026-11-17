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

// Inspect reads path as a .app (or nested helper apps under Frameworks/Helpers).
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
	execName := info.CFBundleExecutable
	execPath := ""
	if execName != "" {
		execPath = filepath.Join(appPath, "Contents", "MacOS", execName)
	}
	return port.InspectedProcess{
		BundleID:       info.CFBundleIdentifier,
		ExecutableName: execName,
		Path:           execPath,
		DisplayName:    display,
		IsHelper:       isHelper,
	}, nil
}

func findHelpers(appPath string) ([]port.InspectedProcess, error) {
	roots := []string{
		filepath.Join(appPath, "Contents", "Frameworks"),
		filepath.Join(appPath, "Contents", "Helpers"),
	}
	var helpers []port.InspectedProcess
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("macos inspect helpers: %w", err)
		}
		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasSuffix(name, ".app") {
				continue
			}
			helperPath := filepath.Join(root, name)
			proc, err := readBundle(helperPath, true)
			if err != nil {
				continue
			}
			helpers = append(helpers, proc)
		}
	}
	return helpers, nil
}
