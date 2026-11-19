package service_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

type healthProxyStore struct {
	proxies []domain.ProxyEndpoint
}

func (s healthProxyStore) Get(context.Context, domain.ProxyID) (domain.ProxyEndpoint, error) {
	return domain.ProxyEndpoint{}, nil
}
func (s healthProxyStore) List(context.Context) ([]domain.ProxyEndpoint, error) {
	return s.proxies, nil
}
func (s healthProxyStore) Save(context.Context, domain.ProxyEndpoint) error { return nil }
func (s healthProxyStore) Delete(context.Context, domain.ProxyID) error     { return nil }

type flipChecker struct {
	calls atomic.Int32
}

func (c *flipChecker) Check(_ context.Context, _ domain.ProxyEndpoint, _ []byte) (domain.ProxyHealth, error) {
	n := c.calls.Add(1)
	if n == 1 {
		return domain.ProxyHealthUnavailable, nil
	}
	return domain.ProxyHealthAvailable, nil
}

func TestHealthMonitorOnChangeOnlyWhenFlips(t *testing.T) {
	t.Parallel()
	checker := &flipChecker{}
	var notifies atomic.Int32
	mon := &service.HealthMonitor{
		Proxies: healthProxyStore{proxies: []domain.ProxyEndpoint{{
			ID: "p1", Name: "local", Kind: domain.ProxyKindSOCKS5, Host: "127.0.0.1", Port: 1,
		}}},
		Checker: checker,
		OnChange: func() {
			notifies.Add(1)
		},
	}
	if err := mon.Recheck(context.Background()); err != nil {
		t.Fatal(err)
	}
	if notifies.Load() != 1 {
		t.Fatalf("first recheck notifies = %d", notifies.Load())
	}
	if err := mon.Recheck(context.Background()); err != nil {
		t.Fatal(err)
	}
	if notifies.Load() != 2 {
		t.Fatalf("second recheck (flip) notifies = %d", notifies.Load())
	}
	if err := mon.Recheck(context.Background()); err != nil {
		t.Fatal(err)
	}
	if notifies.Load() != 2 {
		t.Fatalf("third recheck (unchanged) notifies = %d, want 2", notifies.Load())
	}
}

type slowChecker struct {
	started chan struct{}
	release chan struct{}
}

func (c *slowChecker) Check(context.Context, domain.ProxyEndpoint, []byte) (domain.ProxyHealth, error) {
	select {
	case c.started <- struct{}{}:
	default:
	}
	<-c.release
	return domain.ProxyHealthAvailable, nil
}

func TestHealthMonitorOneRecheckAtATime(t *testing.T) {
	t.Parallel()
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	mon := &service.HealthMonitor{
		Proxies: healthProxyStore{proxies: []domain.ProxyEndpoint{{
			ID: "p1", Name: "local", Kind: domain.ProxyKindSOCKS5, Host: "127.0.0.1", Port: 1,
		}}},
		Checker: &slowChecker{started: started, release: release},
	}
	done1 := make(chan error, 1)
	go func() { done1 <- mon.Recheck(context.Background()) }()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("first recheck did not start")
	}
	done2 := make(chan error, 1)
	go func() { done2 <- mon.Recheck(context.Background()) }()
	select {
	case err := <-done2:
		t.Fatalf("second recheck finished before first released: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	if err := <-done1; err != nil {
		t.Fatal(err)
	}
	if err := <-done2; err != nil {
		t.Fatal(err)
	}
}

var _ port.ProxyHealthChecker = (*flipChecker)(nil)
var _ port.ProxyHealthChecker = (*slowChecker)(nil)
