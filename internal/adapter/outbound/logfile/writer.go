// Package logfile provides a size-rotating file writer for RouteDeck logs.
package logfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	// DefaultMaxBytes is 5 MiB per active log file.
	DefaultMaxBytes int64 = 5 << 20
	// DefaultKeep is the number of rotated files retained (.1 .. .N).
	DefaultKeep = 3
)

// Writer appends to a log file and rotates when size exceeds MaxBytes.
type Writer struct {
	path     string
	maxBytes int64
	keep     int

	mu   sync.Mutex
	file *os.File
	size int64
}

// Options configures Open.
type Options struct {
	MaxBytes int64
	Keep     int
}

// Open creates or appends path with rotation. Directory is created at 0o700.
func Open(path string, opts Options) (*Writer, error) {
	if path == "" {
		return nil, fmt.Errorf("logfile open: empty path")
	}
	maxBytes := opts.MaxBytes
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}
	keep := opts.Keep
	if keep < 0 {
		return nil, fmt.Errorf("logfile open: keep must be >= 0")
	}
	if keep == 0 {
		keep = DefaultKeep
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("logfile open mkdir: %w", err)
	}
	w := &Writer{
		path:     path,
		maxBytes: maxBytes,
		keep:     keep,
	}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

// DefaultPath returns ~/Library/Logs/route-deck/route-deck.log on macOS.
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("logfile default path: %w", err)
	}
	return filepath.Join(home, "Library", "Logs", "route-deck", "route-deck.log"), nil
}

func (w *Writer) openFile() error {
	file, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("logfile open file: %w", err)
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("logfile stat: %w", err)
	}
	w.file = file
	w.size = info.Size()
	return nil
}

// Write implements io.Writer. Concurrent writes are serialized.
func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return 0, fmt.Errorf("logfile write: closed")
	}
	if w.size > 0 && w.size+int64(len(p)) > w.maxBytes {
		if err := w.rotateLocked(); err != nil {
			return 0, err
		}
	}
	n, err := w.file.Write(p)
	w.size += int64(n)
	if err != nil {
		return n, fmt.Errorf("logfile write: %w", err)
	}
	return n, nil
}

func (w *Writer) rotateLocked() error {
	if err := w.file.Close(); err != nil {
		return fmt.Errorf("logfile rotate close: %w", err)
	}
	w.file = nil
	for i := w.keep; i >= 1; i-- {
		src := w.rotatedPath(i)
		if i == w.keep {
			_ = os.Remove(src)
			continue
		}
		dst := w.rotatedPath(i + 1)
		if _, err := os.Stat(src); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("logfile rotate stat %s: %w", src, err)
		}
		if err := os.Rename(src, dst); err != nil {
			return fmt.Errorf("logfile rotate rename %s: %w", src, err)
		}
	}
	if err := os.Rename(w.path, w.rotatedPath(1)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("logfile rotate active: %w", err)
	}
	return w.openFile()
}

func (w *Writer) rotatedPath(n int) string {
	return fmt.Sprintf("%s.%d", w.path, n)
}

// Close closes the underlying file.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	if err != nil {
		return fmt.Errorf("logfile close: %w", err)
	}
	return nil
}
