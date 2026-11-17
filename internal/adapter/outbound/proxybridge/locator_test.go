package proxybridge_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/proxybridge"
	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestLocateCustomPathPreferred(t *testing.T) {
	dir := t.TempDir()
	custom := filepath.Join(dir, "proxybridge")
	if err := os.WriteFile(custom, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	locator := &proxybridge.Locator{
		Runner: func(_ context.Context, name string, _ ...string) ([]byte, error) {
			if name != custom {
				t.Fatalf("runner name = %q, want custom", name)
			}
			return []byte("ProxyBridge CLI 4.1.3\n"), nil
		},
	}
	got, err := locator.Locate(context.Background(), custom)
	if err != nil {
		t.Fatalf("Locate: %v", err)
	}
	if got.State != domain.BridgeStateDetected || got.Path != custom || !got.VersionKnown || got.Version != "4.1.3" {
		t.Fatalf("discovery = %+v", got)
	}
}

func TestPlatformPathsDefault(t *testing.T) {
	locator := &proxybridge.Locator{}
	paths := locator.PlatformPaths()
	if len(paths) == 0 {
		// Non-darwin adapters intentionally return no documented paths yet.
		return
	}
	for _, path := range paths {
		if path == "" {
			t.Fatal("empty platform path")
		}
	}
}

func TestLocateUsesInjectablePaths(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "pb")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	locator := &proxybridge.Locator{
		Paths: []string{bin},
		Runner: func(_ context.Context, name string, _ ...string) ([]byte, error) {
			if name != bin {
				t.Fatalf("runner name = %q", name)
			}
			return []byte("ProxyBridge CLI 4.1.3\n"), nil
		},
	}
	got, err := locator.Locate(context.Background(), "")
	if err != nil {
		t.Fatalf("Locate: %v", err)
	}
	if got.Path != bin || got.State != domain.BridgeStateDetected {
		t.Fatalf("discovery = %+v", got)
	}
}

func TestLocateNotInstalled(t *testing.T) {
	locator := &proxybridge.Locator{
		Paths: []string{"/no/such/proxybridge"},
		Stat: func(string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
		Runner: func(context.Context, string, ...string) ([]byte, error) {
			t.Fatal("runner must not run")
			return nil, nil
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, err := locator.Locate(ctx, "")
	if err != nil {
		t.Fatalf("Locate: %v", err)
	}
	if got.State != domain.BridgeStateNotInstalled {
		t.Fatalf("state = %q", got.State)
	}
}
