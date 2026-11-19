package wails_test

import (
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/wails"
)

func TestDirtySaveableKinds(t *testing.T) {
	t.Parallel()
	d := &wails.EditorDialogs{}
	_ = (&wails.EditorDialogsService{Dialogs: d}).SetDirty("profile", true)
	_ = (&wails.EditorDialogsService{Dialogs: d}).SetDirty("helpers", true)
	got := d.DirtySaveableKinds()
	if len(got) != 1 || got[0] != "profile" {
		t.Fatalf("got %#v", got)
	}
}

func TestAckSaveForQuitFailureDoesNotQuit(t *testing.T) {
	t.Parallel()
	d := &wails.EditorDialogs{}
	_ = (&wails.EditorDialogsService{Dialogs: d}).SetDirty("proxy", true)
	done := make(chan error, 1)
	go func() {
		done <- d.RunSaveForQuit()
	}()
	// Wait until wait map is ready.
	time.Sleep(20 * time.Millisecond)
	d.AckSaveForQuit("proxy", "save failed")
	select {
	case err := <-done:
		if err == nil || err.Error() != "save failed" {
			t.Fatalf("err = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}
