// Package proxyhealth probes TCP and SOCKS5 proxy health.
package proxyhealth

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// Checker implements port.ProxyHealthChecker.
type Checker struct {
	DialTimeout time.Duration
}

var _ port.ProxyHealthChecker = (*Checker)(nil)

func (c *Checker) timeout() time.Duration {
	if c != nil && c.DialTimeout > 0 {
		return c.DialTimeout
	}
	return 3 * time.Second
}

// Check probes the endpoint. Password may be nil/empty for no-auth SOCKS5.
func (c *Checker) Check(ctx context.Context, endpoint domain.ProxyEndpoint, password []byte) (domain.ProxyHealth, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("proxy health: %w", err)
	}
	addr := net.JoinHostPort(endpoint.Host, fmt.Sprintf("%d", endpoint.Port))
	dialer := net.Dialer{Timeout: c.timeout()}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return domain.ProxyHealthUnavailable, nil
	}
	defer func() { _ = conn.Close() }()

	switch endpoint.Kind {
	case domain.ProxyKindHTTP:
		return domain.ProxyHealthAvailable, nil
	case domain.ProxyKindSOCKS5:
		return checkSOCKS5(ctx, conn, endpoint.Username, password)
	default:
		return domain.ProxyHealthUnknown, fmt.Errorf("proxy health: unknown kind %q", endpoint.Kind)
	}
}

func checkSOCKS5(ctx context.Context, conn net.Conn, username string, password []byte) (domain.ProxyHealth, error) {
	_ = ctx
	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	methods := []byte{0x00}
	if username != "" {
		methods = []byte{0x02, 0x00}
	}
	greeting := append([]byte{0x05, byte(len(methods))}, methods...)
	if _, err := conn.Write(greeting); err != nil {
		return domain.ProxyHealthUnavailable, nil
	}
	resp := make([]byte, 2)
	if _, err := conn.Read(resp); err != nil {
		return domain.ProxyHealthUnavailable, nil
	}
	if resp[0] != 0x05 {
		return domain.ProxyHealthUnavailable, nil
	}
	switch resp[1] {
	case 0x00:
		// no auth
	case 0x02:
		if err := socks5UserPass(conn, username, password); err != nil {
			return domain.ProxyHealthAuthenticationFailed, nil
		}
	case 0xFF:
		return domain.ProxyHealthAuthenticationFailed, nil
	default:
		return domain.ProxyHealthUnavailable, nil
	}
	// UDP ASSOCIATE request: CMD=0x03, RSV=0, ATYP=1, ADDR=0.0.0.0, PORT=0
	udpReq := []byte{0x05, 0x03, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
	if _, err := conn.Write(udpReq); err != nil {
		return domain.ProxyHealthTCPAvailableUDPUnavailable, nil
	}
	udpResp := make([]byte, 10)
	n, err := conn.Read(udpResp)
	if err != nil || n < 2 || udpResp[1] != 0x00 {
		return domain.ProxyHealthTCPAvailableUDPUnavailable, nil
	}
	return domain.ProxyHealthAvailable, nil
}

func socks5UserPass(conn net.Conn, username string, password []byte) error {
	u := []byte(username)
	p := password
	if len(u) > 255 || len(p) > 255 {
		return fmt.Errorf("credentials too long")
	}
	msg := make([]byte, 0, 3+len(u)+len(p))
	msg = append(msg, 0x01, byte(len(u)))
	msg = append(msg, u...)
	msg = append(msg, byte(len(p)))
	msg = append(msg, p...)
	if _, err := conn.Write(msg); err != nil {
		return err
	}
	resp := make([]byte, 2)
	if _, err := conn.Read(resp); err != nil {
		return err
	}
	if resp[1] != 0x00 {
		return fmt.Errorf("auth rejected")
	}
	return nil
}
