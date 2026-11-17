package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type memProxyStore struct {
	byID  map[domain.ProxyID]domain.ProxyEndpoint
	saved int
}

func (m *memProxyStore) Get(ctx context.Context, id domain.ProxyID) (domain.ProxyEndpoint, error) {
	if err := ctx.Err(); err != nil {
		return domain.ProxyEndpoint{}, err
	}
	ep, ok := m.byID[id]
	if !ok {
		return domain.ProxyEndpoint{}, fmt.Errorf("missing")
	}
	return ep, nil
}

func (m *memProxyStore) List(ctx context.Context) ([]domain.ProxyEndpoint, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := make([]domain.ProxyEndpoint, 0, len(m.byID))
	for _, ep := range m.byID {
		out = append(out, ep)
	}
	return out, nil
}

func (m *memProxyStore) Save(ctx context.Context, endpoint domain.ProxyEndpoint) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.byID == nil {
		m.byID = map[domain.ProxyID]domain.ProxyEndpoint{}
	}
	m.byID[endpoint.ID] = endpoint
	m.saved++
	return nil
}

func (m *memProxyStore) Delete(ctx context.Context, id domain.ProxyID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	delete(m.byID, id)
	return nil
}

type memProfileStore struct {
	byID  map[domain.ProfileID]domain.Profile
	saved int
}

func (m *memProfileStore) Get(ctx context.Context, id domain.ProfileID) (domain.Profile, error) {
	if err := ctx.Err(); err != nil {
		return domain.Profile{}, err
	}
	p, ok := m.byID[id]
	if !ok {
		return domain.Profile{}, fmt.Errorf("missing")
	}
	return p, nil
}

func (m *memProfileStore) List(ctx context.Context) ([]domain.Profile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := make([]domain.Profile, 0, len(m.byID))
	for _, p := range m.byID {
		out = append(out, p)
	}
	return out, nil
}

func (m *memProfileStore) Save(ctx context.Context, profile domain.Profile) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.byID == nil {
		m.byID = map[domain.ProfileID]domain.Profile{}
	}
	m.byID[profile.ID] = profile
	m.saved++
	return nil
}

func (m *memProfileStore) Delete(ctx context.Context, id domain.ProfileID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	delete(m.byID, id)
	return nil
}

type noopAppStore struct{}

func (noopAppStore) Get(context.Context, domain.ApplicationID) (domain.Application, error) {
	return domain.Application{}, nil
}
func (noopAppStore) List(context.Context) ([]domain.Application, error) { return nil, nil }
func (noopAppStore) Save(context.Context, domain.Application) error     { return nil }
func (noopAppStore) Delete(context.Context, domain.ApplicationID) error { return nil }

type noopProcessStore struct{}

func (noopProcessStore) Get(context.Context, domain.ProcessIdentityID) (domain.ProcessIdentity, error) {
	return domain.ProcessIdentity{}, nil
}
func (noopProcessStore) List(context.Context) ([]domain.ProcessIdentity, error) { return nil, nil }
func (noopProcessStore) Save(context.Context, domain.ProcessIdentity) error     { return nil }
func (noopProcessStore) Delete(context.Context, domain.ProcessIdentityID) error {
	return nil
}

type recordingProcessStore struct {
	rows []domain.ProcessIdentity
}

func (m *recordingProcessStore) Get(ctx context.Context, id domain.ProcessIdentityID) (domain.ProcessIdentity, error) {
	if err := ctx.Err(); err != nil {
		return domain.ProcessIdentity{}, err
	}
	for _, row := range m.rows {
		if row.ID == id {
			return row, nil
		}
	}
	return domain.ProcessIdentity{}, fmt.Errorf("missing")
}

func (m *recordingProcessStore) List(ctx context.Context) ([]domain.ProcessIdentity, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := make([]domain.ProcessIdentity, len(m.rows))
	copy(out, m.rows)
	return out, nil
}

func (m *recordingProcessStore) Save(ctx context.Context, identity domain.ProcessIdentity) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for i, row := range m.rows {
		if row.ID == identity.ID {
			m.rows[i] = identity
			return nil
		}
	}
	m.rows = append(m.rows, identity)
	return nil
}

func (m *recordingProcessStore) Delete(ctx context.Context, id domain.ProcessIdentityID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	next := m.rows[:0]
	for _, row := range m.rows {
		if row.ID != id {
			next = append(next, row)
		}
	}
	m.rows = next
	return nil
}

type memAppStore struct {
	byID map[domain.ApplicationID]domain.Application
}

func (m *memAppStore) Get(ctx context.Context, id domain.ApplicationID) (domain.Application, error) {
	if err := ctx.Err(); err != nil {
		return domain.Application{}, err
	}
	app, ok := m.byID[id]
	if !ok {
		return domain.Application{}, fmt.Errorf("missing")
	}
	return app, nil
}
func (m *memAppStore) List(context.Context) ([]domain.Application, error) { return nil, nil }
func (m *memAppStore) Save(context.Context, domain.Application) error     { return nil }
func (m *memAppStore) Delete(context.Context, domain.ApplicationID) error { return nil }

func TestConfirmHelpersSkipsDuplicateBundleID(t *testing.T) {
	processes := &recordingProcessStore{rows: []domain.ProcessIdentity{
		{
			ID:             "proc-main",
			ApplicationID:  "app-1",
			BundleID:       "com.example.app",
			ExecutableName: "App",
		},
		{
			ID:             "proc-helper",
			ApplicationID:  "app-1",
			BundleID:       "com.example.helper",
			ExecutableName: "HelperA",
			Path:           "/tmp/HelperA.app/Contents/MacOS/HelperA",
		},
	}}
	catalog := &service.Catalog{
		Apps: &memAppStore{byID: map[domain.ApplicationID]domain.Application{
			"app-1": {ID: "app-1", BundleID: "com.example.app", DisplayName: "App"},
		}},
		Processes: processes,
		Profiles:  &memProfileStore{},
	}
	err := catalog.ConfirmHelpers(context.Background(), "app-1", []port.InspectedProcess{
		{BundleID: "com.example.helper", ExecutableName: "HelperB", Path: "/tmp/HelperB.app/Contents/MacOS/HelperB"},
		{BundleID: "com.example.helper", ExecutableName: "HelperC", Path: "/tmp/HelperC.app/Contents/MacOS/HelperC"},
		{BundleID: "com.example.other", ExecutableName: "Other", Path: "/tmp/Other.app/Contents/MacOS/Other"},
		{BundleID: "com.example.other", ExecutableName: "OtherCopy", Path: "/tmp/OtherCopy.app/Contents/MacOS/OtherCopy"},
	})
	if err != nil {
		t.Fatalf("ConfirmHelpers: %v", err)
	}
	var helperCount, otherCount int
	var helperPath string
	for _, row := range processes.rows {
		switch row.BundleID {
		case "com.example.helper":
			helperCount++
			helperPath = row.Path
		case "com.example.other":
			otherCount++
		}
	}
	if helperCount != 1 {
		t.Fatalf("helper bundle rows = %d, want 1", helperCount)
	}
	if helperPath != "/tmp/HelperB.app/Contents/MacOS/HelperB" {
		t.Fatalf("helper path refreshed = %q, want HelperB", helperPath)
	}
	if otherCount != 1 {
		t.Fatalf("other bundle rows = %d, want 1 (batch deduped)", otherCount)
	}
}

func TestConfirmHelpersRemovesUnchecked(t *testing.T) {
	processes := &recordingProcessStore{rows: []domain.ProcessIdentity{
		{ID: "proc-main", ApplicationID: "app-1", BundleID: "com.example.app", ExecutableName: "App"},
		{ID: "proc-helper", ApplicationID: "app-1", BundleID: "com.example.helper", ExecutableName: "Helper"},
	}}
	catalog := &service.Catalog{
		Apps: &memAppStore{byID: map[domain.ApplicationID]domain.Application{
			"app-1": {ID: "app-1", BundleID: "com.example.app"},
		}},
		Processes: processes,
		Profiles:  &memProfileStore{},
	}
	err := catalog.ConfirmHelpers(context.Background(), "app-1", nil)
	if err != nil {
		t.Fatalf("ConfirmHelpers: %v", err)
	}
	for _, row := range processes.rows {
		if row.BundleID == "com.example.helper" {
			t.Fatalf("unchecked helper still present: %+v", row)
		}
	}
}

func TestConfirmHelpersReferenced(t *testing.T) {
	processes := &recordingProcessStore{rows: []domain.ProcessIdentity{
		{ID: "proc-main", ApplicationID: "app-1", BundleID: "com.example.app", ExecutableName: "App"},
		{ID: "proc-helper", ApplicationID: "app-1", BundleID: "com.example.helper", ExecutableName: "Helper"},
	}}
	catalog := &service.Catalog{
		Apps: &memAppStore{byID: map[domain.ApplicationID]domain.Application{
			"app-1": {ID: "app-1", BundleID: "com.example.app"},
		}},
		Processes: processes,
		Profiles: &memProfileStore{byID: map[domain.ProfileID]domain.Profile{
			"p1": {ID: "p1", Name: "Work", Policies: []domain.ApplicationPolicy{
				{ProcessID: "proc-helper", Action: domain.BlockAction()},
			}},
		}},
	}
	err := catalog.ConfirmHelpers(context.Background(), "app-1", nil)
	if !errors.Is(err, domain.ErrReferenced) {
		t.Fatalf("error = %v, want ErrReferenced", err)
	}
}

func TestSaveProfileRejectsUnresolvedAndUnset(t *testing.T) {
	proxies := &memProxyStore{byID: map[domain.ProxyID]domain.ProxyEndpoint{}}
	profiles := &memProfileStore{}
	catalog := &service.Catalog{
		Proxies:   proxies,
		Apps:      noopAppStore{},
		Processes: noopProcessStore{},
		Profiles:  profiles,
	}
	ctx := context.Background()

	t.Run("unset_action", func(t *testing.T) {
		err := catalog.SaveProfile(ctx, domain.Profile{
			ID:   "p1",
			Name: "Work",
			Policies: []domain.ApplicationPolicy{
				{ProcessID: "proc-1", Action: domain.RoutingAction{}},
			},
		})
		if !errors.Is(err, domain.ErrUnsetAction) {
			t.Fatalf("error = %v, want ErrUnsetAction", err)
		}
		if profiles.saved != 0 {
			t.Fatal("nothing must be saved after unset rejection")
		}
	})

	t.Run("unresolved_proxy", func(t *testing.T) {
		action, err := domain.ProxyAction("missing")
		if err != nil {
			t.Fatalf("ProxyAction: %v", err)
		}
		err = catalog.SaveProfile(ctx, domain.Profile{
			ID:       "p1",
			Name:     "Work",
			Policies: []domain.ApplicationPolicy{{ProcessID: "proc-1", Action: action}},
		})
		if !errors.Is(err, domain.ErrUnresolvedProxy) {
			t.Fatalf("error = %v, want ErrUnresolvedProxy", err)
		}
		if profiles.saved != 0 {
			t.Fatal("nothing must be saved after unresolved rejection")
		}
	})
}

func TestSnapshotUnknownNeverActive(t *testing.T) {
	endpoint, err := domain.NewProxyEndpoint("proxy-1", "lab", domain.ProxyKindSOCKS5, "127.0.0.1", 1080, "kc")
	if err != nil {
		t.Fatalf("NewProxyEndpoint: %v", err)
	}
	proxies := &memProxyStore{byID: map[domain.ProxyID]domain.ProxyEndpoint{
		endpoint.ID: endpoint,
	}}
	snap := &service.Snapshotter{Proxies: proxies}
	got, err := snap.Snapshot(context.Background())
	if err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if got.Protection != domain.ProtectionStateUnknown {
		t.Fatalf("Protection = %q, want Unknown", got.Protection)
	}
	if got.Protection == domain.ProtectionStateActive {
		t.Fatal("must never be Active without CLI")
	}
	if len(got.Proxies) != 1 || got.Proxies[0].Health != domain.ProxyHealthUnknown {
		t.Fatalf("proxy rows = %+v", got.Proxies)
	}
	if got.DatabaseError != "" {
		t.Fatalf("DatabaseError = %q, want empty", got.DatabaseError)
	}
}

func TestSnapshotDatabaseFailed(t *testing.T) {
	snap := &service.Snapshotter{
		Proxies:       &memProxyStore{byID: map[domain.ProxyID]domain.ProxyEndpoint{}},
		DatabaseError: "open failed",
	}
	got, err := snap.Snapshot(context.Background())
	if err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if got.Protection != domain.ProtectionStateUnknown {
		t.Fatalf("Protection = %q, want Unknown", got.Protection)
	}
	if got.DatabaseError != "open failed" {
		t.Fatalf("DatabaseError = %q", got.DatabaseError)
	}
	if len(got.Proxies) != 0 {
		t.Fatalf("proxies must be empty on DB failure, got %d", len(got.Proxies))
	}
}
