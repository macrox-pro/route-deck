package proxyhealth

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

// LocalPortProber dials 127.0.0.1 for TCP listen checks (ADR-0013).
type LocalPortProber struct {
	Timeout time.Duration
}

var _ port.LocalPortProber = (*LocalPortProber)(nil)

func (p *LocalPortProber) timeout() time.Duration {
	if p != nil && p.Timeout > 0 {
		return p.Timeout
	}
	return 300 * time.Millisecond
}

// Open returns true when TCP dial to 127.0.0.1:port succeeds.
func (p *LocalPortProber) Open(ctx context.Context, port int) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if port < 1 || port > 65535 {
		return false, fmt.Errorf("local port probe: invalid port %d", port)
	}
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	dialer := net.Dialer{Timeout: p.timeout()}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, nil
	}
	_ = conn.Close()
	return true, nil
}
