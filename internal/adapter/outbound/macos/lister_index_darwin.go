//go:build darwin

package macos

import (
	"context"
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -framework CoreServices -framework CoreFoundation
#include <CoreServices/CoreServices.h>
#include <CoreFoundation/CoreFoundation.h>

static CFArrayRef rd_query_app_paths(CFErrorRef *outErr) {
    CFStringRef q = CFSTR("kMDItemContentTypeTree == 'com.apple.application-bundle'");
    MDQueryRef query = MDQueryCreate(kCFAllocatorDefault, q, NULL, NULL);
    if (query == NULL) {
        return NULL;
    }
    Boolean ok = MDQueryExecute(query, kMDQuerySynchronous);
    if (!ok) {
        CFRelease(query);
        return NULL;
    }
    CFIndex count = MDQueryGetResultCount(query);
    CFMutableArrayRef paths = CFArrayCreateMutable(kCFAllocatorDefault, count, &kCFTypeArrayCallBacks);
    for (CFIndex i = 0; i < count; i++) {
        MDItemRef item = (MDItemRef)MDQueryGetResultAtIndex(query, i);
        if (item == NULL) {
            continue;
        }
        CFTypeRef pathRef = MDItemCopyAttribute(item, kMDItemPath);
        if (pathRef == NULL) {
            continue;
        }
        if (CFGetTypeID(pathRef) == CFStringGetTypeID()) {
            CFArrayAppendValue(paths, pathRef);
        }
        CFRelease(pathRef);
    }
    CFRelease(query);
    return paths;
}
*/
import "C"

func queryApplicationIndex(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	arr := C.rd_query_app_paths(nil)
	if arr == 0 {
		return nil, fmt.Errorf("macos application index: query failed")
	}
	defer C.CFRelease(C.CFTypeRef(arr))
	count := int(C.CFArrayGetCount(arr))
	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		ref := C.CFArrayGetValueAtIndex(arr, C.CFIndex(i))
		if ref == nil {
			continue
		}
		path := cfStringToGo(C.CFStringRef(ref))
		if path != "" {
			out = append(out, path)
		}
	}
	return out, nil
}

func cfStringToGo(s C.CFStringRef) string {
	if s == 0 {
		return ""
	}
	length := C.CFStringGetLength(s)
	if length == 0 {
		return ""
	}
	bufSize := C.CFStringGetMaximumSizeForEncoding(length, C.kCFStringEncodingUTF8) + 1
	buf := make([]byte, bufSize)
	if C.CFStringGetCString(s, (*C.char)(unsafe.Pointer(&buf[0])), bufSize, C.kCFStringEncodingUTF8) == 0 {
		return ""
	}
	n := 0
	for n < len(buf) && buf[n] != 0 {
		n++
	}
	return string(buf[:n])
}
