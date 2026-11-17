//go:build darwin

package macos

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"sync"
	"unsafe"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

/*
#cgo LDFLAGS: -framework AppKit -framework Foundation
#cgo CFLAGS: -x objective-c
#include <stdlib.h>
#include <AppKit/AppKit.h>
#include <Foundation/Foundation.h>

// AppKit icon extraction must be serialized by the Go caller (iconGenMu).
// Do not dispatch_sync onto the main queue: concurrent callers previously
// raced inside NSGraphicsContext, and main-queue sync can deadlock under
// Wails when the UI thread is busy. CGImageForProposedRect avoids a shared
// graphics-context stack.
static CFDataRef rd_icon_png(const char *path, size_t *outLen) {
    @autoreleasepool {
        if (path == NULL || outLen == NULL) {
            return NULL;
        }
        *outLen = 0;
        NSString *nsPath = [NSString stringWithUTF8String:path];
        if (nsPath == nil) {
            return NULL;
        }
        NSImage *icon = [[[NSWorkspace sharedWorkspace] iconForFile:nsPath] copy];
        if (icon == nil) {
            return NULL;
        }
        [icon setSize:NSMakeSize(32, 32)];
        NSRect proposed = NSMakeRect(0, 0, 32, 32);
        CGImageRef cgImage = [icon CGImageForProposedRect:&proposed context:nil hints:nil];
        if (cgImage == NULL) {
            return NULL;
        }
        NSBitmapImageRep *rep = [[NSBitmapImageRep alloc] initWithCGImage:cgImage];
        if (rep == nil) {
            return NULL;
        }
        NSData *png = [rep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
        if (png == nil || [png length] == 0) {
            return NULL;
        }
        *outLen = (size_t)[png length];
        return CFDataCreate(kCFAllocatorDefault, (const UInt8 *)[png bytes], (CFIndex)[png length]);
    }
}
*/
import "C"

// iconGenMu serializes all AppKit icon generation process-wide.
// Concurrent rd_icon_png previously SIGSEGV'd when Applications loaded many icons.
var iconGenMu sync.Mutex

// IconReader loads Finder icons as PNG bytes.
type IconReader struct {
	mu    sync.Mutex
	cache map[string][]byte
}

var _ port.ApplicationIconReader = (*IconReader)(nil)

// PNG returns a 32×32 PNG for path, cached in memory.
func (r *IconReader) PNG(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("macos icon: %w", err)
	}
	if path == "" {
		return nil, fmt.Errorf("macos icon: empty path")
	}

	r.mu.Lock()
	if r.cache != nil {
		if hit, ok := r.cache[path]; ok {
			out := append([]byte(nil), hit...)
			r.mu.Unlock()
			return out, nil
		}
	}
	r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("macos icon: %w", err)
	}

	iconGenMu.Lock()
	defer iconGenMu.Unlock()

	// Re-check cache after acquiring the generation lock: another caller may
	// have filled it while we waited.
	r.mu.Lock()
	if r.cache != nil {
		if hit, ok := r.cache[path]; ok {
			out := append([]byte(nil), hit...)
			r.mu.Unlock()
			return out, nil
		}
	}
	r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("macos icon: %w", err)
	}

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	var length C.size_t
	data := C.rd_icon_png(cPath, &length)
	if data == 0 || length == 0 {
		return nil, fmt.Errorf("macos icon: unavailable for %s", path)
	}
	defer C.CFRelease(C.CFTypeRef(data))
	raw := C.GoBytes(unsafe.Pointer(C.CFDataGetBytePtr(data)), C.int(length))
	if !looksLikePNG(raw) {
		// Defensive: ensure PNG; re-encode if needed.
		img, _, err := image.Decode(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("macos icon decode: %w", err)
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("macos icon encode: %w", err)
		}
		raw = buf.Bytes()
	}

	r.mu.Lock()
	if r.cache == nil {
		r.cache = map[string][]byte{}
	}
	r.cache[path] = append([]byte(nil), raw...)
	r.mu.Unlock()
	return append([]byte(nil), raw...), nil
}

func looksLikePNG(b []byte) bool {
	return len(b) >= 8 &&
		b[0] == 0x89 && b[1] == 0x50 && b[2] == 0x4e && b[3] == 0x47 &&
		b[4] == 0x0d && b[5] == 0x0a && b[6] == 0x1a && b[7] == 0x0a
}
