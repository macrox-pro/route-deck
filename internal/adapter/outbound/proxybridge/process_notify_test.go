package proxybridge

import (
	"strings"
	"sync/atomic"
	"testing"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

func TestExtensionApprovalNotifyOnce(t *testing.T) {
	var calls atomic.Int32
	p := &Process{
		LineNotify: func(class port.BridgeLineClass) {
			if class != port.BridgeLineExtensionApproval {
				t.Fatalf("class = %s", class)
			}
			calls.Add(1)
		},
	}
	p.ready = make(chan struct{})
	p.failed = make(chan error, 1)
	p.approvalNotified.Store(false)

	line := "Extension: waiting for your approval. Open System Settings > Privacy & Security and click Allow.\n"
	p.readLines(strings.NewReader(line + line))
	if got := calls.Load(); got != 1 {
		t.Fatalf("LineNotify calls = %d, want 1", got)
	}
}
