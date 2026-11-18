package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type fakeLister struct {
	apps []port.InstalledApplication
	err  error
}

func (f fakeLister) List(context.Context) ([]port.InstalledApplication, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.apps, nil
}

type fakeProber struct {
	open map[int]bool
}

func (f fakeProber) Open(_ context.Context, port int) (bool, error) {
	return f.open[port], nil
}

type memProxyForGuides struct {
	byID map[domain.ProxyID]domain.ProxyEndpoint
}

func (m *memProxyForGuides) Get(_ context.Context, id domain.ProxyID) (domain.ProxyEndpoint, error) {
	ep, ok := m.byID[id]
	if !ok {
		return domain.ProxyEndpoint{}, errors.New("not found")
	}
	return ep, nil
}

func (m *memProxyForGuides) List(context.Context) ([]domain.ProxyEndpoint, error) {
	out := make([]domain.ProxyEndpoint, 0, len(m.byID))
	for _, ep := range m.byID {
		out = append(out, ep)
	}
	return out, nil
}

func (m *memProxyForGuides) Save(_ context.Context, endpoint domain.ProxyEndpoint) error {
	if m.byID == nil {
		m.byID = map[domain.ProxyID]domain.ProxyEndpoint{}
	}
	m.byID[endpoint.ID] = endpoint
	return nil
}

func (m *memProxyForGuides) Delete(_ context.Context, id domain.ProxyID) error {
	delete(m.byID, id)
	return nil
}

type memAppsForGuides struct {
	apps []domain.Application
}

func (m *memAppsForGuides) Get(context.Context, domain.ApplicationID) (domain.Application, error) {
	return domain.Application{}, errors.New("not found")
}
func (m *memAppsForGuides) List(context.Context) ([]domain.Application, error) { return m.apps, nil }
func (m *memAppsForGuides) Save(context.Context, domain.Application) error     { return nil }
func (m *memAppsForGuides) Delete(context.Context, domain.ApplicationID) error { return nil }

func TestDefaultGuideCatalog(t *testing.T) {
	catalog := service.DefaultGuideCatalog()
	if len(catalog) != 8 {
		t.Fatalf("len = %d, want 8", len(catalog))
	}
	wantIDs := []string{"proxybridge", "shadowrocket", "v2rayn", "hiddify", "happ", "singbox", "surge", "generic"}
	for i, id := range wantIDs {
		if catalog[i].ID != id {
			t.Fatalf("catalog[%d].ID = %q, want %q", i, catalog[i].ID, id)
		}
	}
	for _, g := range catalog {
		t.Run(g.ID, func(t *testing.T) {
			if g.Title == "" {
				t.Fatal("empty title")
			}
			if len(g.Steps) == 0 {
				t.Fatal("no steps")
			}
			if g.ID != "proxybridge" && g.ID != "generic" {
				hasWarn := false
				for _, step := range g.Steps {
					if step.Warning {
						hasWarn = true
						break
					}
				}
				if !hasWarn {
					t.Fatal("client guide missing warning step")
				}
			}
			if g.DocsURL != "" && !strings.HasPrefix(g.DocsURL, "https://") {
				t.Fatalf("docs url must be https: %q", g.DocsURL)
			}
			switch g.ID {
			case "shadowrocket":
				if len(g.Endpoints) != 1 || g.Endpoints[0].Port != 1082 || g.Endpoints[0].Kind != domain.ProxyKindSOCKS5 {
					t.Fatalf("shadowrocket endpoints = %+v", g.Endpoints)
				}
			case "v2rayn":
				if g.Endpoints[0].Port != 10808 {
					t.Fatalf("port = %d", g.Endpoints[0].Port)
				}
			case "hiddify":
				if len(g.BundleIDs) != 0 {
					t.Fatal("hiddify must have Detection unavailable (no bundle IDs)")
				}
				if g.Endpoints[0].Port != 12334 {
					t.Fatalf("port = %d", g.Endpoints[0].Port)
				}
			case "happ":
				if g.Endpoints[0].Port != 10808 {
					t.Fatalf("port = %d", g.Endpoints[0].Port)
				}
			case "surge":
				if g.Endpoints[0].Port != 6153 {
					t.Fatalf("port = %d", g.Endpoints[0].Port)
				}
			case "singbox", "generic":
				if g.AddMode != service.AddModeFormFallback {
					t.Fatalf("add mode = %q", g.AddMode)
				}
				if len(g.Endpoints) != 0 {
					t.Fatal("expected no default endpoints")
				}
			case "proxybridge":
				if g.AddMode != service.AddModeNone {
					t.Fatalf("add mode = %q", g.AddMode)
				}
			}
		})
	}
}

func TestSetupGuidesListSortAndDetection(t *testing.T) {
	proxies := &memProxyForGuides{byID: map[domain.ProxyID]domain.ProxyEndpoint{}}
	sg := &service.SetupGuides{
		Lister: fakeLister{apps: []port.InstalledApplication{
			{DisplayName: "Shadowrocket", BundleID: "com.liguangming.Shadowrocket", Path: "/Applications/Shadowrocket.app"},
		}},
		Prober:  fakeProber{open: map[int]bool{1082: true}},
		Proxies: proxies,
		Catalog: &service.Catalog{Proxies: proxies},
	}
	got, err := sg.List(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}
	if got[0].ID != "proxybridge" {
		t.Fatalf("first = %q", got[0].ID)
	}
	if got[len(got)-1].ID != "generic" {
		t.Fatalf("last = %q", got[len(got)-1].ID)
	}
	var shadow *service.GuideView
	var hiddify *service.GuideView
	for i := range got {
		if got[i].ID == "shadowrocket" {
			shadow = &got[i]
		}
		if got[i].ID == "hiddify" {
			hiddify = &got[i]
		}
	}
	if shadow == nil || shadow.Detection != service.DetectionInstalled || !shadow.PortOpen || shadow.InstalledPath == "" {
		t.Fatalf("shadowrocket = %+v", shadow)
	}
	if hiddify == nil || hiddify.Detection != service.DetectionUnavailable {
		t.Fatalf("hiddify = %+v", hiddify)
	}
}

func TestSetupGuidesAddProxy(t *testing.T) {
	tests := []struct {
		name    string
		guideID string
		wantErr error
		port    int
	}{
		{name: "shadowrocket", guideID: "shadowrocket", port: 1082},
		{name: "v2rayn", guideID: "v2rayn", port: 10808},
		{name: "hiddify", guideID: "hiddify", port: 12334},
		{name: "happ", guideID: "happ", port: 10808},
		{name: "surge", guideID: "surge", port: 6153},
		{name: "singbox", guideID: "singbox", wantErr: service.ErrGuideNeedsPort},
		{name: "generic", guideID: "generic", wantErr: service.ErrGuideNeedsPort},
		{name: "proxybridge", guideID: "proxybridge", wantErr: service.ErrGuideHasNoProxy},
		{name: "unknown", guideID: "nope", wantErr: service.ErrUnknownGuideID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxies := &memProxyForGuides{byID: map[domain.ProxyID]domain.ProxyEndpoint{}}
			sg := &service.SetupGuides{
				Proxies: proxies,
				Catalog: &service.Catalog{Proxies: proxies},
			}
			ep, err := sg.AddProxy(context.Background(), tt.guideID)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if ep.Host != "127.0.0.1" || ep.Port != tt.port || ep.Kind != domain.ProxyKindSOCKS5 {
				t.Fatalf("endpoint = %+v", ep)
			}
			_, err = sg.AddProxy(context.Background(), tt.guideID)
			if !errors.Is(err, service.ErrEndpointAlreadyAdded) {
				t.Fatalf("duplicate err = %v", err)
			}
		})
	}
}

func TestSetupGuidesSuggestionsAndDismiss(t *testing.T) {
	proxies := &memProxyForGuides{byID: map[domain.ProxyID]domain.ProxyEndpoint{}}
	sg := &service.SetupGuides{
		Lister: fakeLister{apps: []port.InstalledApplication{
			{BundleID: "com.liguangming.Shadowrocket", Path: "/Applications/Shadowrocket.app"},
		}},
		Prober:  fakeProber{open: map[int]bool{1082: true}},
		Proxies: proxies,
		Catalog: &service.Catalog{Proxies: proxies},
	}
	sug, err := sg.Suggestions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(sug) != 1 || sug[0].GuideID != "shadowrocket" || sug[0].Port != 1082 {
		t.Fatalf("suggestions = %+v", sug)
	}
	if err := sg.DismissSuggestion("shadowrocket", 1082); err != nil {
		t.Fatal(err)
	}
	sug, err = sg.Suggestions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(sug) != 0 {
		t.Fatalf("after dismiss = %+v", sug)
	}
}

func TestSetupGuidesChecklist(t *testing.T) {
	sg := &service.SetupGuides{
		Apps: &memAppsForGuides{apps: []domain.Application{{ID: "a1", DisplayName: "App"}}},
		Snapshot: func(context.Context) (service.SnapshotResult, error) {
			return service.SnapshotResult{
				Bridge:           domain.BridgeStateReady,
				AppliedConfirmed: true,
				Proxies: []service.ProxyHealthRow{
					{ID: "p1", Health: domain.ProxyHealthAvailable},
				},
			}, nil
		},
	}
	got, err := sg.Checklist(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !got.Complete || !got.BridgeFound || !got.ProxyAvailable || !got.ApplicationAdded || !got.ProfileApplied {
		t.Fatalf("checklist = %+v", got)
	}
	if got.ProxyPartial || got.ProxyStep != service.ChecklistProxyAvailable {
		t.Fatalf("full available fields = %+v", got)
	}

	sg.Snapshot = func(context.Context) (service.SnapshotResult, error) {
		return service.SnapshotResult{Bridge: domain.BridgeStateNotInstalled}, nil
	}
	got, err = sg.Checklist(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got.BridgeFound || got.Complete {
		t.Fatalf("regression checklist = %+v", got)
	}
}

func TestSetupGuidesChecklistProxyPriority(t *testing.T) {
	tests := []struct {
		name      string
		proxies   []service.ProxyHealthRow
		available bool
		partial   bool
		step      service.ChecklistProxyStep
		complete  bool
	}{
		{
			name: "available beats udp limited",
			proxies: []service.ProxyHealthRow{
				{ID: "a", Health: domain.ProxyHealthTCPAvailableUDPUnavailable},
				{ID: "b", Health: domain.ProxyHealthAvailable},
			},
			available: true,
			partial:   false,
			step:      service.ChecklistProxyAvailable,
			complete:  true,
		},
		{
			name: "udp limited beats checking",
			proxies: []service.ProxyHealthRow{
				{ID: "a", Health: domain.ProxyHealthChecking},
				{ID: "b", Health: domain.ProxyHealthTCPAvailableUDPUnavailable},
			},
			available: true,
			partial:   true,
			step:      service.ChecklistProxyUDPLimited,
			complete:  true,
		},
		{
			name: "checking only",
			proxies: []service.ProxyHealthRow{
				{ID: "a", Health: domain.ProxyHealthChecking},
			},
			available: true,
			partial:   true,
			step:      service.ChecklistProxyChecking,
			complete:  true,
		},
		{
			name: "unavailable does not count",
			proxies: []service.ProxyHealthRow{
				{ID: "a", Health: domain.ProxyHealthUnavailable},
			},
			available: false,
			partial:   false,
			step:      service.ChecklistProxyMissing,
			complete:  false,
		},
		{
			name: "unknown does not count",
			proxies: []service.ProxyHealthRow{
				{ID: "a", Health: domain.ProxyHealthUnknown},
			},
			available: false,
			partial:   false,
			step:      service.ChecklistProxyMissing,
			complete:  false,
		},
		{
			name: "auth failed does not count",
			proxies: []service.ProxyHealthRow{
				{ID: "a", Health: domain.ProxyHealthAuthenticationFailed},
			},
			available: false,
			partial:   false,
			step:      service.ChecklistProxyMissing,
			complete:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sg := &service.SetupGuides{
				Apps: &memAppsForGuides{apps: []domain.Application{{ID: "a1", DisplayName: "App"}}},
				Snapshot: func(context.Context) (service.SnapshotResult, error) {
					return service.SnapshotResult{
						Bridge:           domain.BridgeStateReady,
						AppliedConfirmed: true,
						Proxies:          tt.proxies,
					}, nil
				},
			}
			got, err := sg.Checklist(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if got.ProxyAvailable != tt.available || got.ProxyPartial != tt.partial || got.ProxyStep != tt.step || got.Complete != tt.complete {
				t.Fatalf("got available=%v partial=%v step=%q complete=%v", got.ProxyAvailable, got.ProxyPartial, got.ProxyStep, got.Complete)
			}
		})
	}
}

func TestSetupGuidesChecklistDismiss(t *testing.T) {
	bridge := domain.BridgeStateReady
	var applied bool
	apps := &memAppsForGuides{apps: []domain.Application{{ID: "a1", DisplayName: "App"}}}
	sg := &service.SetupGuides{
		Apps: apps,
		Snapshot: func(context.Context) (service.SnapshotResult, error) {
			return service.SnapshotResult{
				Bridge:           bridge,
				AppliedConfirmed: applied,
				Proxies: []service.ProxyHealthRow{
					{ID: "p1", Health: domain.ProxyHealthTCPAvailableUDPUnavailable},
				},
			}, nil
		},
	}

	got, err := sg.Checklist(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got.Dismissed || !got.ProxyAvailable || !got.ProxyPartial || got.Complete {
		t.Fatalf("before dismiss = %+v", got)
	}

	if err := sg.DismissChecklist(context.Background()); err != nil {
		t.Fatal(err)
	}
	got, err = sg.Checklist(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !got.Dismissed {
		t.Fatalf("after dismiss = %+v", got)
	}

	apps.apps = append(apps.apps, domain.Application{ID: "a2", DisplayName: "App2"})
	got, err = sg.Checklist(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !got.Dismissed {
		t.Fatalf("improve after dismiss cleared dismiss = %+v", got)
	}

	bridge = domain.BridgeStateNotInstalled
	got, err = sg.Checklist(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got.Dismissed || got.BridgeFound {
		t.Fatalf("after regression = %+v", got)
	}
}
