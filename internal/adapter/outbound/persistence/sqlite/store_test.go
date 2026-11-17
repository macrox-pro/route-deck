package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/persistence/sqlite"
	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

func portSettings(debug bool, desired string, bridgePath string) port.Settings {
	return port.Settings{
		DesiredProfileID: domain.ProfileID(desired),
		DesiredSet:       desired != "",
		ProxyBridgePath:  bridgePath,
		DebugLogging:     debug,
	}
}

func openTestDB(t *testing.T) *sqlite.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "route-deck.db")
	db, err := sqlite.Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	return db
}

func TestRoundTripAggregates(t *testing.T) {
	db := openTestDB(t)
	store := db.Store()
	ctx := context.Background()

	app := domain.Application{
		ID:          "app-1",
		DisplayName: "Mail",
		BundleID:    "com.apple.mail",
		IconRef:     "mail",
	}
	if err := store.Applications().Save(ctx, app); err != nil {
		t.Fatalf("SaveApplication: %v", err)
	}
	gotApp, err := store.Applications().Get(ctx, "app-1")
	if err != nil {
		t.Fatalf("GetApplication: %v", err)
	}
	if gotApp != app {
		t.Fatalf("application mismatch: %+v", gotApp)
	}

	identity := domain.ProcessIdentity{
		ID:             "proc-1",
		ApplicationID:  "app-1",
		BundleID:       "com.apple.mail",
		ExecutableName: "Mail",
		Path:           "/Applications/Mail.app",
		SigningHint:    "team",
	}
	if err := store.ProcessIdentities().Save(ctx, identity); err != nil {
		t.Fatalf("SaveProcessIdentity: %v", err)
	}
	gotProc, err := store.ProcessIdentities().Get(ctx, "proc-1")
	if err != nil {
		t.Fatalf("GetProcessIdentity: %v", err)
	}
	if gotProc != identity {
		t.Fatalf("process mismatch: %+v", gotProc)
	}

	endpoint, err := domain.NewProxyEndpointWithUsername("proxy-1", "lab", domain.ProxyKindSOCKS5, "127.0.0.1", 1080, "alice", "kc-ref")
	if err != nil {
		t.Fatalf("NewProxyEndpointWithUsername: %v", err)
	}
	if err := store.Save(ctx, endpoint); err != nil {
		t.Fatalf("Save proxy: %v", err)
	}
	gotProxy, err := store.Get(ctx, "proxy-1")
	if err != nil {
		t.Fatalf("Get proxy: %v", err)
	}
	if gotProxy.ID != endpoint.ID || gotProxy.Host != endpoint.Host || gotProxy.Port != endpoint.Port || gotProxy.Username != "alice" || !gotProxy.SendDomain {
		t.Fatalf("proxy mismatch: %+v", gotProxy)
	}

	action, err := domain.ProxyAction("proxy-1")
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}
	policy, err := domain.NewApplicationPolicy("proc-1", action, domain.ProtocolScopeBOTH, domain.UDPLimitationAbsent)
	if err != nil {
		t.Fatalf("NewApplicationPolicy: %v", err)
	}
	profile, err := domain.NewProfile("profile-1", "Work", []domain.ApplicationPolicy{policy})
	if err != nil {
		t.Fatalf("NewProfile: %v", err)
	}
	if err := store.Profiles().Save(ctx, profile); err != nil {
		t.Fatalf("SaveProfile: %v", err)
	}
	gotProfile, err := store.Profiles().Get(ctx, "profile-1")
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if gotProfile.Name != "Work" || len(gotProfile.Policies) != 1 {
		t.Fatalf("profile mismatch: %+v", gotProfile)
	}
	id, ok := gotProfile.Policies[0].Action.ProxyID()
	if !ok || id != "proxy-1" {
		t.Fatalf("policy proxy id = (%q, %v)", id, ok)
	}
	if gotProfile.Policies[0].ProtocolScope != domain.ProtocolScopeBOTH {
		t.Fatalf("policy scope = %q", gotProfile.Policies[0].ProtocolScope)
	}

	settings := store.Settings()
	if err := settings.Save(ctx, portSettings(true, "profile-1", "/usr/local/bin/proxybridge")); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	gotSettings, err := settings.Get(ctx)
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}
	if !gotSettings.DebugLogging || !gotSettings.DesiredSet || gotSettings.DesiredProfileID != "profile-1" {
		t.Fatalf("settings mismatch: %+v", gotSettings)
	}

	apps, err := store.Applications().List(ctx)
	if err != nil || len(apps) != 1 {
		t.Fatalf("ListApplications: %v len=%d", err, len(apps))
	}
}

func TestUnsetActionRejected(t *testing.T) {
	db := openTestDB(t)
	store := db.Store()
	ctx := context.Background()

	if err := store.Applications().Save(ctx, domain.Application{ID: "app-1", DisplayName: "Mail", BundleID: "com.apple.mail"}); err != nil {
		t.Fatalf("SaveApplication: %v", err)
	}
	if err := store.ProcessIdentities().Save(ctx, domain.ProcessIdentity{
		ID: "proc-1", ApplicationID: "app-1", BundleID: "com.apple.mail", ExecutableName: "Mail",
	}); err != nil {
		t.Fatalf("SaveProcessIdentity: %v", err)
	}

	err := store.Profiles().Save(ctx, domain.Profile{
		ID:   "profile-1",
		Name: "Bad",
		Policies: []domain.ApplicationPolicy{
			{ProcessID: "proc-1", Action: domain.RoutingAction{}},
		},
	})
	if !errors.Is(err, domain.ErrUnsetAction) {
		t.Fatalf("SaveProfile unset error = %v, want ErrUnsetAction", err)
	}
}

func TestMissingProxyFKNotDirect(t *testing.T) {
	db := openTestDB(t)
	store := db.Store()
	ctx := context.Background()

	if err := store.Applications().Save(ctx, domain.Application{ID: "app-1", DisplayName: "Mail", BundleID: "com.apple.mail"}); err != nil {
		t.Fatalf("SaveApplication: %v", err)
	}
	if err := store.ProcessIdentities().Save(ctx, domain.ProcessIdentity{
		ID: "proc-1", ApplicationID: "app-1", BundleID: "com.apple.mail", ExecutableName: "Mail",
	}); err != nil {
		t.Fatalf("SaveProcessIdentity: %v", err)
	}

	action, err := domain.ProxyAction("missing-proxy")
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}
	err = store.Profiles().Save(ctx, domain.Profile{
		ID:       "profile-1",
		Name:     "Work",
		Policies: []domain.ApplicationPolicy{{ProcessID: "proc-1", Action: action}},
	})
	if err == nil {
		t.Fatal("SaveProfile with missing proxy must fail")
	}
	if errors.Is(err, domain.ErrUnsetAction) {
		t.Fatal("missing proxy must not be treated as unset/Direct")
	}

	_, getErr := store.Profiles().Get(ctx, "profile-1")
	if !errors.Is(getErr, sqlite.ErrNotFound) {
		t.Fatalf("GetProfile after failed save: %v, want ErrNotFound", getErr)
	}
}

func TestGetMissingIsNotFound(t *testing.T) {
	db := openTestDB(t)
	_, err := db.Store().Get(context.Background(), "nope")
	if !errors.Is(err, sqlite.ErrNotFound) {
		t.Fatalf("Get missing = %v, want ErrNotFound", err)
	}
}

func TestDeleteReferencedProxyRejected(t *testing.T) {
	db := openTestDB(t)
	store := db.Store()
	ctx := context.Background()

	if err := store.Applications().Save(ctx, domain.Application{ID: "app-1", DisplayName: "Mail", BundleID: "com.apple.mail"}); err != nil {
		t.Fatalf("SaveApplication: %v", err)
	}
	if err := store.ProcessIdentities().Save(ctx, domain.ProcessIdentity{
		ID: "proc-1", ApplicationID: "app-1", BundleID: "com.apple.mail", ExecutableName: "Mail",
	}); err != nil {
		t.Fatalf("SaveProcessIdentity: %v", err)
	}
	endpoint, err := domain.NewProxyEndpoint("proxy-1", "lab", domain.ProxyKindSOCKS5, "127.0.0.1", 1080, "kc")
	if err != nil {
		t.Fatalf("NewProxyEndpoint: %v", err)
	}
	if err := store.Save(ctx, endpoint); err != nil {
		t.Fatalf("Save proxy: %v", err)
	}
	action, err := domain.ProxyAction("proxy-1")
	if err != nil {
		t.Fatalf("ProxyAction: %v", err)
	}
	policy, err := domain.NewApplicationPolicy("proc-1", action, domain.ProtocolScopeTCP, "")
	if err != nil {
		t.Fatalf("NewApplicationPolicy: %v", err)
	}
	if err := store.Profiles().Save(ctx, domain.Profile{ID: "p1", Name: "Work", Policies: []domain.ApplicationPolicy{policy}}); err != nil {
		t.Fatalf("SaveProfile: %v", err)
	}
	if err := store.Delete(ctx, "proxy-1"); !errors.Is(err, domain.ErrReferenced) {
		t.Fatalf("Delete proxy = %v, want ErrReferenced", err)
	}
	if err := store.Applications().Delete(ctx, "app-1"); !errors.Is(err, domain.ErrReferenced) {
		t.Fatalf("Delete app = %v, want ErrReferenced", err)
	}
}

func TestDuplicateBundleAndProfileName(t *testing.T) {
	db := openTestDB(t)
	store := db.Store()
	ctx := context.Background()

	if err := store.Applications().Save(ctx, domain.Application{ID: "app-1", DisplayName: "Mail", BundleID: "com.apple.mail"}); err != nil {
		t.Fatalf("SaveApplication: %v", err)
	}
	err := store.Applications().Save(ctx, domain.Application{ID: "app-2", DisplayName: "Other", BundleID: "com.apple.mail"})
	if !errors.Is(err, domain.ErrDuplicateBundleID) {
		t.Fatalf("duplicate bundle = %v, want ErrDuplicateBundleID", err)
	}
	if err := store.Profiles().Save(ctx, domain.Profile{ID: "p1", Name: "Work"}); err != nil {
		t.Fatalf("SaveProfile: %v", err)
	}
	err = store.Profiles().Save(ctx, domain.Profile{ID: "p2", Name: "Work"})
	if !errors.Is(err, domain.ErrDuplicateProfileName) {
		t.Fatalf("duplicate name = %v, want ErrDuplicateProfileName", err)
	}
}
