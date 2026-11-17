package logfile

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestOpenEmptyPath(t *testing.T) {
	_, err := Open("", Options{})
	if err == nil {
		t.Fatal("Open() error = nil, want error")
	}
}

func TestWriteAndReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")

	w, err := Open(path, Options{MaxBytes: 1024, Keep: 3})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if _, err := w.Write([]byte("hello\n")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "hello\n" {
		t.Fatalf("file = %q, want %q", data, "hello\n")
	}

	w2, err := Open(path, Options{MaxBytes: 1024, Keep: 3})
	if err != nil {
		t.Fatalf("reopen Open() error = %v", err)
	}
	defer func() {
		if err := w2.Close(); err != nil {
			t.Errorf("reopen Close() error = %v", err)
		}
	}()
	if _, err := w2.Write([]byte("world\n")); err != nil {
		t.Fatalf("reopen Write() error = %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile after reopen error = %v", err)
	}
	if string(data) != "hello\nworld\n" {
		t.Fatalf("file after reopen = %q, want %q", data, "hello\nworld\n")
	}
}

func TestRotationKeepsThree(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")

	w, err := Open(path, Options{MaxBytes: 32, Keep: 3})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if err := w.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	}()

	chunk := []byte("abcdefghijklmnopqrstuvwxyz012345\n") // 33 bytes
	for i := 0; i < 5; i++ {
		if _, err := w.Write(chunk); err != nil {
			t.Fatalf("Write(%d) error = %v", i, err)
		}
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("active log missing: %v", err)
	}
	for i := 1; i <= 3; i++ {
		rotated := path + "." + strconv.Itoa(i)
		if _, err := os.Stat(rotated); err != nil {
			t.Fatalf("rotated %s missing: %v", rotated, err)
		}
	}
	if _, err := os.Stat(path + ".4"); !os.IsNotExist(err) {
		t.Fatalf("path.4 should not exist, err=%v", err)
	}
}
