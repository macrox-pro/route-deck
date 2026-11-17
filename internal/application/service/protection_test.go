package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type memProfiles map[domain.ProfileID]domain.Profile

func (m memProfiles) Get(_ context.Context, id domain.ProfileID) (domain.Profile, error) {
	p, ok := m[id]
	if !ok {
		return domain.Profile{}, errors.New("not found")
	}
	return p, nil
}
func (m memProfiles) List(context.Context) ([]domain.Profile, error) { return nil, nil }
func (m memProfiles) Save(context.Context, domain.Profile) error     { return nil }
func (m memProfiles) Delete(context.Context, domain.ProfileID) error { return nil }

type memProxies []domain.ProxyEndpoint

func (m memProxies) Get(context.Context, domain.ProxyID) (domain.ProxyEndpoint, error) {
	return domain.ProxyEndpoint{}, errors.New("unused")
}
func (m memProxies) List(context.Context) ([]domain.ProxyEndpoint, error) { return m, nil }
func (m memProxies) Save(context.Context, domain.ProxyEndpoint) error     { return nil }
func (m memProxies) Delete(context.Context, domain.ProxyID) error         { return nil }

type memProcesses []domain.ProcessIdentity

func (m memProcesses) Get(context.Context, domain.ProcessIdentityID) (domain.ProcessIdentity, error) {
	return domain.ProcessIdentity{}, errors.New("unused")
}
func (m memProcesses) List(context.Context) ([]domain.ProcessIdentity, error) { return m, nil }
func (m memProcesses) Save(context.Context, domain.ProcessIdentity) error     { return nil }
func (m memProcesses) Delete(context.Context, domain.ProcessIdentityID) error { return nil }

type memSettings struct {
	mu       sync.Mutex
	settings port.Settings
}

func (m *memSettings) Get(context.Context) (port.Settings, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.settings, nil
}
func (m *memSettings) Save(_ context.Context, settings port.Settings) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings = settings
	return nil
}

type fakeLocator struct {
	disc port.BridgeDiscovery
}

func (f fakeLocator) Locate(context.Context, string) (port.BridgeDiscovery, error) {
	return f.disc, nil
}

func (f fakeLocator) PlatformPaths() []string {
	if f.disc.Path == "" {
		return nil
	}
	return []string{f.disc.Path}
}

type fakeCompiler struct {
	json []byte
	err  error
}

func (f fakeCompiler) Compile(context.Context, domain.Profile, []domain.ProxyEndpoint, []domain.ProcessIdentity) (port.CompiledProfile, error) {
	if f.err != nil {
		return port.CompiledProfile{}, f.err
	}
	return port.CompiledProfile{JSON: f.json}, nil
}

type fakeProcess struct {
	mu            sync.Mutex
	running       bool
	startErr      error
	confirmErr    error
	confirmDelay  time.Duration
	starts        int
	stops         int
	lastProfile   []byte
	blockConfirm  chan struct{}
	confirmClosed bool
}

func (f *fakeProcess) Start(_ context.Context, _ string, profileJSON []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.starts++
	f.lastProfile = append([]byte(nil), profileJSON...)
	if f.startErr != nil {
		return f.startErr
	}
	f.running = true
	return nil
}

func (f *fakeProcess) Stop(context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stops++
	f.running = false
	return nil
}

func (f *fakeProcess) WaitConfirm(ctx context.Context) error {
	if f.blockConfirm != nil && !f.confirmClosed {
		select {
		case <-f.blockConfirm:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if f.confirmDelay > 0 {
		time.Sleep(f.confirmDelay)
	}
	return f.confirmErr
}

func (f *fakeProcess) Running() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.running
}

func testProtection(t *testing.T) (*service.Protection, *fakeProcess, *memSettings) {
	t.Helper()
	settings := &memSettings{settings: port.Settings{}}
	proc := &fakeProcess{}
	p := &service.Protection{
		Profiles: memProfiles{
			"prof-1": {ID: "prof-1", Name: "one"},
		},
		Proxies:   memProxies{},
		Processes: memProcesses{},
		Settings:  settings,
		Locator: fakeLocator{disc: port.BridgeDiscovery{
			State:        domain.BridgeStateDetected,
			Path:         "/bin/proxybridge",
			VersionKnown: true,
			Version:      "4.1.3",
		}},
		Compiler: fakeCompiler{json: []byte(`{"dnsEnabled":true}`)},
		Process:  proc,
	}
	return p, proc, settings
}

func TestProtectionApplyConfirm(t *testing.T) {
	p, proc, settings := testProtection(t)
	ctx := context.Background()
	if err := p.Apply(ctx, "prof-1", "op-1"); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	id, ok := p.Applied()
	if !ok || id != "prof-1" {
		t.Fatalf("Applied = %q %v", id, ok)
	}
	if p.BridgeState() != domain.BridgeStateReady {
		t.Fatalf("bridge = %s", p.BridgeState())
	}
	if !settings.settings.DesiredSet || settings.settings.DesiredProfileID != "prof-1" {
		t.Fatalf("settings = %+v", settings.settings)
	}
	if proc.starts != 1 {
		t.Fatalf("starts = %d", proc.starts)
	}
}

func TestProtectionApplyRollback(t *testing.T) {
	p, proc, _ := testProtection(t)
	ctx := context.Background()
	if err := p.Apply(ctx, "prof-1", "op-1"); err != nil {
		t.Fatalf("first Apply: %v", err)
	}
	proc.confirmErr = errors.New(string(port.BridgeLineRejectedRules))
	err := p.Apply(ctx, "prof-1", "op-2")
	if err == nil {
		t.Fatal("expected confirm failure")
	}
	if !errors.Is(err, err) && err.Error() == "" {
		t.Fatal("empty error")
	}
	// rollback restarts previous JSON
	if proc.starts < 2 {
		t.Fatalf("starts = %d, want rollback start", proc.starts)
	}
}

func TestProtectionConflict(t *testing.T) {
	p, proc, _ := testProtection(t)
	proc.startErr = errors.New(string(port.BridgeLineAnotherCLI))
	err := p.Apply(context.Background(), "prof-1", "op-c")
	if err == nil {
		t.Fatal("expected conflict")
	}
	if p.BridgeState() != domain.BridgeStateConflict {
		t.Fatalf("bridge = %s", p.BridgeState())
	}
}

func TestProtectionConcurrentOp(t *testing.T) {
	p, proc, _ := testProtection(t)
	block := make(chan struct{})
	proc.blockConfirm = block
	done := make(chan error, 1)
	go func() {
		done <- p.Apply(context.Background(), "prof-1", "op-a")
	}()
	time.Sleep(20 * time.Millisecond)
	err := p.Apply(context.Background(), "prof-1", "op-b")
	if err == nil || err.Error() != "protection: operation in progress" {
		t.Fatalf("concurrent = %v", err)
	}
	close(block)
	if err := <-done; err != nil {
		t.Fatalf("first Apply: %v", err)
	}
}

func TestProtectionStop(t *testing.T) {
	p, _, _ := testProtection(t)
	ctx := context.Background()
	if err := p.Apply(ctx, "prof-1", "op-1"); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if err := p.Stop(ctx, "op-stop"); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	if _, ok := p.Applied(); ok {
		t.Fatal("applied still set")
	}
	if p.BridgeState() != domain.BridgeStateStopped {
		t.Fatalf("bridge = %s", p.BridgeState())
	}
}
