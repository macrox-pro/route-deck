package port_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type fakeProxyStore struct {
	endpoint domain.ProxyEndpoint
}

func (f fakeProxyStore) Get(ctx context.Context, id domain.ProxyID) (domain.ProxyEndpoint, error) {
	if err := ctx.Err(); err != nil {
		return domain.ProxyEndpoint{}, fmt.Errorf("proxy store get: %w", err)
	}
	if id == f.endpoint.ID {
		return f.endpoint, nil
	}
	return domain.ProxyEndpoint{}, fmt.Errorf("proxy store get: id not in fake")
}

func (f fakeProxyStore) List(ctx context.Context) ([]domain.ProxyEndpoint, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("proxy store list: %w", err)
	}
	return []domain.ProxyEndpoint{f.endpoint}, nil
}

func (f fakeProxyStore) Save(ctx context.Context, endpoint domain.ProxyEndpoint) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("proxy store save: %w", err)
	}
	_ = endpoint
	return nil
}

type fakeSecretStore struct {
	secret []byte
}

func (f fakeSecretStore) Load(ctx context.Context, ref domain.KeychainRef) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("secret store load: %w", err)
	}
	if ref == "kc-ref" {
		return f.secret, nil
	}
	return nil, fmt.Errorf("secret store load: ref not in fake")
}

var _ port.ProxyStore = fakeProxyStore{}
var _ port.SecretStore = fakeSecretStore{}

func TestStores(t *testing.T) {
	endpoint, err := domain.NewProxyEndpoint(
		"proxy-1",
		"lab",
		domain.ProxyKindSOCKS5,
		"127.0.0.1",
		1080,
		"kc-ref",
	)
	if err != nil {
		t.Fatalf("NewProxyEndpoint: %v", err)
	}

	proxies := fakeProxyStore{endpoint: endpoint}
	secrets := fakeSecretStore{secret: []byte{1}}

	got, err := proxies.Get(context.Background(), "proxy-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != endpoint.ID || got.Host != endpoint.Host || got.Port != endpoint.Port || got.Kind != endpoint.Kind || !got.SendDomain {
		t.Fatalf("Get() endpoint mismatch: got %+v", got)
	}

	listed, err := proxies.List(context.Background())
	if err != nil || len(listed) != 1 {
		t.Fatalf("List: %v len=%d", err, len(listed))
	}

	loaded, err := secrets.Load(context.Background(), "kc-ref")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !bytes.Equal(loaded, []byte{1}) {
		t.Fatal("secret mismatch")
	}
}
