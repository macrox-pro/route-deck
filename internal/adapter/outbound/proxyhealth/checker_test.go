package proxyhealth_test

import (
	"bufio"
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/proxyhealth"
	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestCheckHTTP(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password []byte
		respond  func(conn net.Conn, req string)
		want     domain.ProxyHealth
	}{
		{
			name: "200_available",
			respond: func(conn net.Conn, _ string) {
				_, _ = io.WriteString(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
			},
			want: domain.ProxyHealthAvailable,
		},
		{
			name: "407_no_creds",
			respond: func(conn net.Conn, _ string) {
				_, _ = io.WriteString(conn, "HTTP/1.1 407 Proxy Authentication Required\r\n\r\n")
			},
			want: domain.ProxyHealthAuthenticationFailed,
		},
		{
			name:     "407_with_password",
			username: "user",
			password: []byte("fake-password"),
			respond: func(conn net.Conn, _ string) {
				_, _ = io.WriteString(conn, "HTTP/1.1 407 Proxy Authentication Required\r\n\r\n")
			},
			want: domain.ProxyHealthAuthenticationFailed,
		},
		{
			name: "502_unavailable",
			respond: func(conn net.Conn, _ string) {
				_, _ = io.WriteString(conn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			},
			want: domain.ProxyHealthUnavailable,
		},
		{
			name: "short_read_unavailable",
			respond: func(conn net.Conn, _ string) {
				_ = conn.Close()
			},
			want: domain.ProxyHealthUnavailable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = ln.Close() }()
			type probeResult struct {
				req string
				err string
			}
			done := make(chan probeResult, 1)
			go func() {
				conn, err := ln.Accept()
				if err != nil {
					done <- probeResult{err: err.Error()}
					return
				}
				defer func() { _ = conn.Close() }()
				req, readErr := readHTTPRequestHead(conn)
				if readErr != nil {
					done <- probeResult{err: readErr.Error()}
					return
				}
				tt.respond(conn, req)
				done <- probeResult{req: req}
			}()
			host, port := splitAddr(t, ln.Addr().String())
			endpoint, err := domain.NewProxyEndpointWithUsername("p", "http", domain.ProxyKindHTTP, host, port, tt.username, "kc")
			if err != nil {
				t.Fatal(err)
			}
			got, err := (&proxyhealth.Checker{DialTimeout: time.Second}).Check(context.Background(), endpoint, tt.password)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("health = %q, want %q", got, tt.want)
			}
			select {
			case res := <-done:
				if res.err != "" && tt.want == domain.ProxyHealthAvailable {
					t.Fatalf("server: %s", res.err)
				}
				if res.req != "" {
					if !strings.Contains(res.req, "CONNECT example.com:443") {
						t.Fatalf("missing CONNECT target in %q", res.req)
					}
					hasAuth := strings.Contains(res.req, "Proxy-Authorization:")
					expectAuth := tt.username != "" || len(tt.password) > 0
					if hasAuth != expectAuth {
						t.Fatalf("auth header present=%v want=%v", hasAuth, expectAuth)
					}
					if strings.Contains(res.req, "fake-password") {
						t.Fatal("password appeared in cleartext request")
					}
				}
			case <-time.After(2 * time.Second):
				t.Fatal("server did not finish")
			}
		})
	}
}

func TestCheckHTTPDialUnavailable(t *testing.T) {
	endpoint, err := domain.NewProxyEndpoint("p", "http", domain.ProxyKindHTTP, "127.0.0.1", 1, "kc")
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

func readHTTPRequestHead(conn net.Conn) (string, error) {
	br := bufio.NewReader(conn)
	var lines []string
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			return strings.Join(lines, ""), err
		}
		lines = append(lines, line)
		if line == "\r\n" {
			break
		}
	}
	return strings.Join(lines, ""), nil
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
