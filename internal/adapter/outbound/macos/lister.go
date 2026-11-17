package macos

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

// Lister lists installed applications from an injectable index and roots.
type Lister struct {
	Roots []string
	// HomeDir optional; empty uses os.UserHomeDir for path priority and
	// default user-facing ~/Applications.
	HomeDir string
	// UserFacingRoots optional; empty defaults to /Applications and
	// <HomeDir>/Applications. Paths outside these roots get IsSystem=true.
	UserFacingRoots []string
	IndexPaths      func(ctx context.Context) ([]string, error)
	Logger          *slog.Logger
}

var _ port.InstalledApplicationLister = (*Lister)(nil)

func (l *Lister) logger() *slog.Logger {
	if l != nil && l.Logger != nil {
		return l.Logger
	}
	return slog.Default()
}

// DefaultRoots returns the three Application folders from ADR-0009.
func DefaultRoots() []string {
	home, err := os.UserHomeDir()
	userApps := ""
	if err == nil {
		userApps = filepath.Join(home, "Applications")
	}
	roots := []string{"/Applications"}
	if userApps != "" {
		roots = append(roots, userApps)
	}
	roots = append(roots, "/System/Applications")
	return roots
}

// List unions the index with roots, filters .app bundles, and dedupes by bundle id.
func (l *Lister) List(ctx context.Context) ([]port.InstalledApplication, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("macos list installed: %w", err)
	}
	paths := map[string]struct{}{}
	indexFn := l.IndexPaths
	if indexFn == nil {
		indexFn = queryApplicationIndex
	}
	indexed, err := indexFn(ctx)
	if err != nil {
		l.logger().Debug("application index failed; using roots only", "component", "macos.lister", "err", err)
	} else {
		for _, path := range indexed {
			paths[path] = struct{}{}
		}
	}
	roots := l.Roots
	if roots == nil {
		roots = DefaultRoots()
	}
	for _, root := range roots {
		entries, readErr := os.ReadDir(root)
		if readErr != nil {
			if os.IsNotExist(readErr) {
				continue
			}
			l.logger().Debug("application root unreadable", "component", "macos.lister", "path", root, "err", readErr)
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasSuffix(name, ".app") {
				continue
			}
			paths[filepath.Join(root, name)] = struct{}{}
		}
	}

	type candidate struct {
		app      port.InstalledApplication
		priority int
	}
	byBundle := map[string]candidate{}
	for path := range paths {
		if !isTopLevelAppBundle(path) {
			continue
		}
		info, plistErr := readInfoPlist(path)
		if plistErr != nil || info.CFBundleIdentifier == "" {
			l.logger().Debug("skip installed app", "component", "macos.lister", "path", path, "err", plistErr)
			continue
		}
		display := displayNameFromInfo(info)
		if display == "" {
			display = strings.TrimSuffix(filepath.Base(path), ".app")
		}
		pri := pathPriority(path, l.homeDir())
		next := candidate{
			app: port.InstalledApplication{
				DisplayName: display,
				BundleID:    info.CFBundleIdentifier,
				Path:        path,
				IsSystem:    l.isSystemPath(path),
			},
			priority: pri,
		}
		if existing, ok := byBundle[info.CFBundleIdentifier]; ok && existing.priority <= pri {
			continue
		}
		byBundle[info.CFBundleIdentifier] = next
	}

	out := make([]port.InstalledApplication, 0, len(byBundle))
	for _, row := range byBundle {
		out = append(out, row.app)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DisplayName != out[j].DisplayName {
			return strings.ToLower(out[i].DisplayName) < strings.ToLower(out[j].DisplayName)
		}
		return out[i].Path < out[j].Path
	})
	return out, nil
}

func isTopLevelAppBundle(path string) bool {
	if !strings.HasSuffix(path, ".app") {
		return false
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	// Drop nested .app inside another .app (helpers).
	parent := filepath.Dir(path)
	for parent != "." && parent != string(filepath.Separator) {
		if strings.HasSuffix(parent, ".app") {
			return false
		}
		next := filepath.Dir(parent)
		if next == parent {
			break
		}
		parent = next
	}
	return true
}

func (l *Lister) homeDir() string {
	if l != nil && l.HomeDir != "" {
		return l.HomeDir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

func (l *Lister) userFacingRoots() []string {
	if l != nil && len(l.UserFacingRoots) > 0 {
		return l.UserFacingRoots
	}
	roots := []string{"/Applications"}
	if home := l.homeDir(); home != "" {
		roots = append(roots, filepath.Join(home, "Applications"))
	}
	return roots
}

func (l *Lister) isSystemPath(path string) bool {
	for _, root := range l.userFacingRoots() {
		root = filepath.Clean(root)
		if path == root || strings.HasPrefix(path, root+string(filepath.Separator)) {
			return false
		}
	}
	return true
}

// pathPriority prefers ~/…, then /Applications, then /System/Applications.
func pathPriority(path, home string) int {
	if home != "" {
		home = filepath.Clean(home)
		if path == home || strings.HasPrefix(path, home+string(filepath.Separator)) {
			return 0
		}
	}
	if strings.HasPrefix(path, "/Applications"+string(filepath.Separator)) || path == "/Applications" {
		return 1
	}
	if strings.HasPrefix(path, "/System/Applications"+string(filepath.Separator)) || path == "/System/Applications" {
		return 2
	}
	return 3
}
