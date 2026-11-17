package wails_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/wails"
	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type fakeDialog struct {
	mu       sync.Mutex
	warnings []wails.UserMsg
}

func (f *fakeDialog) Info(wails.UserMsg) {}

func (f *fakeDialog) Warning(msg wails.UserMsg) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.warnings = append(f.warnings, msg)
}

type fakeCoordinator struct {
	err error
}

func (f *fakeCoordinator) Apply(context.Context, domain.ProfileID, string) error {
	return f.err
}

func (f *fakeCoordinator) Stop(context.Context, string) error { return nil }

type fakeSettings struct{}

func (fakeSettings) Get(context.Context) (domain.ProfileID, bool, error) {
	return "prof-1", true, nil
}

func TestApplyShowsWarningOnGUIConflict(t *testing.T) {
	dialog := &fakeDialog{}
	svc := &wails.ProtectionService{
		Coordinator: &fakeCoordinator{err: errors.New("protection apply confirm: bridge process: " + string(port.BridgeLineGUIConflict))},
		Settings:    fakeSettings{},
		NewID:       func() (string, error) { return "op-1", nil },
		Dialog:      dialog,
	}
	if _, err := svc.Apply(context.Background()); err != nil {
		t.Fatalf("Apply sync: %v", err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for {
		dialog.mu.Lock()
		n := len(dialog.warnings)
		dialog.mu.Unlock()
		if n == 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("warnings = %d, want 1", n)
		}
		time.Sleep(10 * time.Millisecond)
	}
	dialog.mu.Lock()
	got := dialog.warnings[0]
	dialog.mu.Unlock()
	if got.ID != wails.UserMsgGUIConflict.ID {
		t.Fatalf("warning id = %q", got.ID)
	}
}
