package domain_test

import (
	"errors"
	"testing"

	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestNewProxyEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		id        domain.ProxyID
		proxyName string
		kind      domain.ProxyKind
		host      string
		port      int
		ref       domain.KeychainRef
		wantErr   error
		want      domain.ProxyEndpoint
	}{
		{
			name:      "port_zero",
			id:        "proxy-1",
			proxyName: "Primary",
			kind:      domain.ProxyKindSOCKS5,
			host:      "127.0.0.1",
			port:      0,
			ref:       "kc-ref",
			wantErr:   domain.ErrInvalidPort,
		},
		{
			name:      "port_65536",
			id:        "proxy-1",
			proxyName: "Primary",
			kind:      domain.ProxyKindSOCKS5,
			host:      "127.0.0.1",
			port:      65536,
			ref:       "kc-ref",
			wantErr:   domain.ErrInvalidPort,
		},
		{
			name:      "empty_host",
			id:        "proxy-1",
			proxyName: "Primary",
			kind:      domain.ProxyKindHTTP,
			host:      "",
			port:      8080,
			ref:       "kc-ref",
			wantErr:   domain.ErrEmptyHost,
		},
		{
			name:      "port_one",
			id:        "proxy-1",
			proxyName: "Primary",
			kind:      domain.ProxyKindSOCKS5,
			host:      "127.0.0.1",
			port:      1,
			ref:       "kc-ref",
			want: domain.ProxyEndpoint{
				ID:          "proxy-1",
				Name:        "Primary",
				Kind:        domain.ProxyKindSOCKS5,
				Host:        "127.0.0.1",
				Port:        1,
				KeychainRef: "kc-ref",
				SendDomain:  true,
			},
		},
		{
			name:      "port_65535",
			id:        "proxy-2",
			proxyName: "Edge",
			kind:      domain.ProxyKindHTTP,
			host:      "proxy.example",
			port:      65535,
			ref:       "kc-edge",
			want: domain.ProxyEndpoint{
				ID:          "proxy-2",
				Name:        "Edge",
				Kind:        domain.ProxyKindHTTP,
				Host:        "proxy.example",
				Port:        65535,
				KeychainRef: "kc-edge",
				SendDomain:  true,
			},
		},
		{
			name:      "valid_send_domain_true",
			id:        "proxy-3",
			proxyName: "Main",
			kind:      domain.ProxyKindSOCKS5,
			host:      "10.0.0.1",
			port:      1080,
			ref:       "kc-main",
			want: domain.ProxyEndpoint{
				ID:          "proxy-3",
				Name:        "Main",
				Kind:        domain.ProxyKindSOCKS5,
				Host:        "10.0.0.1",
				Port:        1080,
				KeychainRef: "kc-main",
				SendDomain:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewProxyEndpoint(tt.id, tt.proxyName, tt.kind, tt.host, tt.port, tt.ref)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewProxyEndpoint() error = %v, want %v", err, tt.wantErr)
				}
				if got != (domain.ProxyEndpoint{}) {
					t.Fatalf("NewProxyEndpoint() on error = %#v, want zero value", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewProxyEndpoint() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("NewProxyEndpoint() = %#v, want %#v", got, tt.want)
			}
			if !got.SendDomain {
				t.Fatal("NewProxyEndpoint() SendDomain = false, want true")
			}
		})
	}
}
