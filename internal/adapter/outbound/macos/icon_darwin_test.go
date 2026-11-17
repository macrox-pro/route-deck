//go:build darwin

package macos

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
)

func TestIconReaderPNGSafari(t *testing.T) {
	path := "/Applications/Safari.app"
	if _, err := os.Stat(path); err != nil {
		t.Skip("Safari.app not present")
	}
	r := &IconReader{}
	png, err := r.PNG(context.Background(), path)
	if err != nil {
		t.Fatalf("PNG: %v", err)
	}
	if !looksLikePNG(png) {
		t.Fatalf("not a PNG (%d bytes)", len(png))
	}
	// Cache hit
	again, err := r.PNG(context.Background(), path)
	if err != nil {
		t.Fatalf("PNG cache: %v", err)
	}
	if len(again) != len(png) {
		t.Fatalf("cache size mismatch")
	}
}

func TestIconReaderPNGConcurrent(t *testing.T) {
	paths := []string{
		"/Applications/Safari.app",
		"/System/Applications/Utilities/Terminal.app",
		"/System/Applications/System Settings.app",
		"/Applications/Safari.app",
	}
	available := make([]string, 0, len(paths))
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			available = append(available, path)
		}
	}
	if len(available) < 2 {
		t.Skip("need at least two apps for concurrent icon test")
	}

	r := &IconReader{}
	var wg sync.WaitGroup
	errCh := make(chan error, len(available)*4)
	for range 4 {
		for _, path := range available {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				png, err := r.PNG(context.Background(), p)
				if err != nil {
					errCh <- err
					return
				}
				if !looksLikePNG(png) {
					errCh <- errors.New("not a PNG")
				}
			}(path)
		}
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatalf("concurrent PNG: %v", err)
	}
}
