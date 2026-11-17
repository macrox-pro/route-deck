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
	if len(locator.customHits) < 2 || locator.customHits[0] != "" || locator.customHits[1] != "/custom/pb" {
		t.Fatalf("locate calls = %#v", locator.customHits)
	}
}
