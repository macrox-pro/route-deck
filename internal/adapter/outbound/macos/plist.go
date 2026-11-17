package macos

import (
	"fmt"
	"os"
	"path/filepath"

	"howett.net/plist"
)

func readInfoPlist(appPath string) (infoPlist, error) {
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")
	raw, err := os.ReadFile(plistPath)
	if err != nil {
		return infoPlist{}, fmt.Errorf("macos read plist: %w", err)
	}
	var info infoPlist
	if _, err := plist.Unmarshal(raw, &info); err != nil {
		return infoPlist{}, fmt.Errorf("macos parse plist: %w", err)
	}
	return info, nil
}

func displayNameFromInfo(info infoPlist) string {
	if info.CFBundleDisplayName != "" {
		return info.CFBundleDisplayName
	}
	return info.CFBundleName
}
