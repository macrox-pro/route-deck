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

func (f fakeProxyStore) Delete(ctx context.Context, id domain.ProxyID) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("proxy store delete: %w", err)
	}
	_ = id
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

func (f fakeSecretStore) Save(ctx context.Context, ref domain.KeychainRef, secret []byte) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("secret store save: %w", err)
	}
	_ = ref
	_ = secret
	return nil
}

func (f fakeSecretStore) Delete(ctx context.Context, ref domain.KeychainRef) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("secret store delete: %w", err)
	}
	_ = ref
	return nil
}

type fakeSettingsStore struct {
	settings port.Settings
}

func (f fakeSettingsStore) Get(ctx context.Context) (port.Settings, error) {
	if err := ctx.Err(); err != nil {
		return port.Settings{}, fmt.Errorf("settings get: %w", err)
	}
	return f.settings, nil
}

func (f fakeSettingsStore) Save(ctx context.Context, settings port.Settings) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("settings save: %w", err)
	}
	_ = settings
	return nil
}

type fakeInspector struct{}

func (fakeInspector) Inspect(ctx context.Context, path string) (port.InspectedApplication, error) {
	if err := ctx.Err(); err != nil {
		return port.InspectedApplication{}, err
	}
	_ = path
	return port.InspectedApplication{
		DisplayName: "Demo",
		BundleID:    "pro.macrox.demo",
		Main: port.InspectedProcess{
			BundleID:       "pro.macrox.demo",
			ExecutableName: "Demo",
			Path:           path,
			DisplayName:    "Demo",
		},
	}, nil
}

var _ port.ProxyStore = fakeProxyStore{}
var _ port.SecretStore = fakeSecretStore{}
var _ port.SettingsStore = fakeSettingsStore{}
var _ port.ApplicationInspector = fakeInspector{}

func TestStores(t *testing.T) {
	endpoint, err := domain.NewProxyEndpointWithUsername(
		"proxy-1",
		"lab",
		domain.ProxyKindSOCKS5,
		"127.0.0.1",
		1080,
		"user",
		"kc-ref",
	)
	if err != nil {
		t.Fatalf("NewProxyEndpointWithUsername: %v", err)
	}

	proxies := fakeProxyStore{endpoint: endpoint}
	secrets := fakeSecretStore{secret: []byte{1}}
	settings := fakeSettingsStore{settings: port.Settings{DebugLogging: true}}
	inspector := fakeInspector{}

	got, err := proxies.Get(context.Background(), "proxy-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Username != "user" || got.ID != endpoint.ID {
		t.Fatalf("Get() endpoint mismatch: got %+v", got)
	}

	if err := proxies.Delete(context.Background(), "proxy-1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := secrets.Save(context.Background(), "kc-ref", []byte{2}); err != nil {
		t.Fatalf("Save secret: %v", err)
	}
	if err := secrets.Delete(context.Background(), "kc-ref"); err != nil {
		t.Fatalf("Delete secret: %v", err)
	}

	loaded, err := secrets.Load(context.Background(), "kc-ref")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !bytes.Equal(loaded, []byte{1}) {
		t.Fatal("secret mismatch")
	}

	gotSettings, err := settings.Get(context.Background())
	if err != nil || !gotSettings.DebugLogging {
		t.Fatalf("settings Get: %v %#v", err, gotSettings)
	}

	inspected, err := inspector.Inspect(context.Background(), "/Applications/Demo.app")
	if err != nil || inspected.BundleID != "pro.macrox.demo" {
		t.Fatalf("Inspect: %v %#v", err, inspected)
	}
}
