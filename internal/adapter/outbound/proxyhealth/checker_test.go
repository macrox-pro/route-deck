package proxyhealth_test

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/proxyhealth"
	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestCheckHTTPAvailable(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		time.Sleep(20 * time.Millisecond)
	}()
	host, port := splitAddr(t, ln.Addr().String())
	endpoint, err := domain.NewProxyEndpoint("p", "http", domain.ProxyKindHTTP, host, port, "kc")
	if err != nil {
		t.Fatal(err)
	}
	got, err := (&proxyhealth.Checker{DialTimeout: time.Second}).Check(context.Background(), endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != domain.ProxyHealthAvailable {
		t.Fatalf("health = %q", got)
	}
}

func TestCheckSOCKS5Available(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	go serveSOCKS5(t, ln, false)
	host, port := splitAddr(t, ln.Addr().String())
	endpoint, err := domain.NewProxyEndpoint("p", "socks", domain.ProxyKindSOCKS5, host, port, "kc")
	if err != nil {
		t.Fatal(err)
	}
	got, err := (&proxyhealth.Checker{DialTimeout: time.Second}).Check(context.Background(), endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != domain.ProxyHealthAvailable {
		t.Fatalf("health = %q", got)
	}
}

func TestCheckUnavailable(t *testing.T) {
	endpoint, err := domain.NewProxyEndpoint("p", "x", domain.ProxyKindSOCKS5, "127.0.0.1", 1, "kc")
	if err != nil {
		t.Fatal(err)
	}
	got, err := (&proxyhealth.Checker{DialTimeout: 50 * time.Millisecond}).Check(context.Background(), endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != domain.ProxyHealthUnavailable {
		t.Fatalf("health = %q", got)
	}
}

func serveSOCKS5(t *testing.T, ln net.Listener, requireAuth bool) {
	t.Helper()
	conn, err := ln.Accept()
	if err != nil {
		return
	}
	defer func() { _ = conn.Close() }()
	buf := make([]byte, 16)
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		return
	}
	nmethods := int(buf[1])
	if _, err := io.ReadFull(conn, buf[:nmethods]); err != nil {
		return
	}
	if requireAuth {
		_, _ = conn.Write([]byte{0x05, 0x02})
		return
	}
	_, _ = conn.Write([]byte{0x05, 0x00})
	req := make([]byte, 10)
	if _, err := io.ReadFull(conn, req); err != nil {
		return
	}
	_, _ = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
}

func splitAddr(t *testing.T, addr string) (string, int) {
	t.Helper()
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatal(err)
	}
	var port int
	_, err = fmtSscanf(portStr, &port)
	if err != nil {
		t.Fatal(err)
	}
	return host, port
}

func fmtSscanf(s string, port *int) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, net.InvalidAddrError(s)
		}
		n = n*10 + int(c-'0')
	}
	*port = n
	return 1, nil
}
