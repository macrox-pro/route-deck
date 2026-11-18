package service_test

import (
	"context"
	"testing"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type settingsStoreFake struct {
	settings port.Settings
}

func (f *settingsStoreFake) Get(context.Context) (port.Settings, error) {
	return f.settings, nil
}

func (f *settingsStoreFake) Save(_ context.Context, settings port.Settings) error {
	f.settings = settings
	return nil
}

type settingsLocatorFake struct {
	platformPath string
	customHits   []string
}

func (f *settingsLocatorFake) PlatformPaths() []string {
	if f.platformPath == "" {
		return nil
	}
	return []string{f.platformPath}
}

func (f *settingsLocatorFake) Locate(_ context.Context, customPath string) (port.BridgeDiscovery, error) {
	f.customHits = append(f.customHits, customPath)
	if customPath == "" {
		if f.platformPath == "" {
			return port.BridgeDiscovery{State: domain.BridgeStateNotInstalled}, nil
		}
		return port.BridgeDiscovery{
			State: domain.BridgeStateDetected,
			Path:  f.platformPath,
		}, nil
	}
	return port.BridgeDiscovery{
		State:        domain.BridgeStateDetected,
		Path:         customPath,
		Version:      "4.1.3",
		VersionKnown: true,
	}, nil
}

func TestSettingsGetDiscoveredPathFromLocator(t *testing.T) {
	store := &settingsStoreFake{settings: port.Settings{ProxyBridgePath: "/custom/pb"}}
	locator := &settingsLocatorFake{platformPath: "/usr/local/bin/proxybridge"}
	svc := &service.Settings{Store: store, Locator: locator}

	view, err := svc.Get(context.Background())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if view.DiscoveredProxyBridgePath != "/usr/local/bin/proxybridge" {
		t.Fatalf("discovered = %q", view.DiscoveredProxyBridgePath)
	}
	if view.ProxyBridgePath != "/custom/pb" {
		t.Fatalf("custom = %q", view.ProxyBridgePath)
	}
	if view.BridgeVersion != "4.1.3" || !view.VersionKnown {
		t.Fatalf("configured discovery = %+v", view)
	}
	if view.Appearance != service.AppearanceSystem {
		t.Fatalf("appearance default = %q", view.Appearance)
	}
	if len(locator.customHits) < 2 || locator.customHits[0] != "" || locator.customHits[1] != "/custom/pb" {
		t.Fatalf("locate calls = %#v", locator.customHits)
	}
}

func TestNormalizeAppearance(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{in: "", want: service.AppearanceSystem, ok: true},
		{in: "system", want: service.AppearanceSystem, ok: true},
		{in: "light", want: service.AppearanceLight, ok: true},
		{in: "dark", want: service.AppearanceDark, ok: true},
		{in: "nope", want: "", ok: false},
	}
	for _, tc := range cases {
		got, ok := service.NormalizeAppearance(tc.in)
		if got != tc.want || ok != tc.ok {
			t.Fatalf("NormalizeAppearance(%q) = (%q, %v); want (%q, %v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestSettingsSetAppearance(t *testing.T) {
	store := &settingsStoreFake{}
	svc := &service.Settings{Store: store}

	view, err := svc.SetAppearance(context.Background(), "dark")
	if err != nil {
		t.Fatalf("SetAppearance: %v", err)
	}
	if view.Appearance != service.AppearanceDark {
		t.Fatalf("view = %q", view.Appearance)
	}
	if store.settings.Appearance != service.AppearanceDark {
		t.Fatalf("stored = %q", store.settings.Appearance)
	}

	view, err = svc.SetAppearance(context.Background(), "neon")
	if err != nil {
		t.Fatalf("SetAppearance coerce: %v", err)
	}
	if view.Appearance != service.AppearanceSystem {
		t.Fatalf("coerced view = %q", view.Appearance)
	}
	if store.settings.Appearance != service.AppearanceSystem {
		t.Fatalf("coerced stored = %q", store.settings.Appearance)
	}
}

func TestSettingsGetRewritesInvalidAppearance(t *testing.T) {
	store := &settingsStoreFake{settings: port.Settings{Appearance: "neon"}}
	svc := &service.Settings{Store: store}

	view, err := svc.Get(context.Background())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if view.Appearance != service.AppearanceSystem {
		t.Fatalf("view = %q", view.Appearance)
	}
	if store.settings.Appearance != service.AppearanceSystem {
		t.Fatalf("stored after heal = %q", store.settings.Appearance)
	}
}

func TestSettingsSetAutoApply(t *testing.T) {
	store := &settingsStoreFake{settings: port.Settings{AutoApply: true, Appearance: service.AppearanceSystem}}
	svc := &service.Settings{Store: store}

	view, err := svc.SetAutoApply(context.Background(), false)
	if err != nil {
		t.Fatalf("SetAutoApply false: %v", err)
	}
	if view.AutoApply {
		t.Fatal("view.AutoApply = true, want false")
	}
	if store.settings.AutoApply {
		t.Fatal("stored AutoApply = true, want false")
	}

	view, err = svc.SetAutoApply(context.Background(), true)
	if err != nil {
		t.Fatalf("SetAutoApply true: %v", err)
	}
	if !view.AutoApply || !store.settings.AutoApply {
		t.Fatal("AutoApply not restored to true")
	}
}
