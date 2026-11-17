package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// HealthMonitor periodically checks proxy health and stores the latest registry.
type HealthMonitor struct {
	Proxies  port.ProxyStore
	Secrets  port.SecretStore
	Checker  port.ProxyHealthChecker
	Logger   *slog.Logger
	Interval time.Duration

	mu     sync.Mutex
	health map[domain.ProxyID]domain.ProxyHealth
	stop   chan struct{}
	done   chan struct{}
}

func (m *HealthMonitor) logger() *slog.Logger {
	if m != nil && m.Logger != nil {
		return m.Logger
	}
	return slog.Default()
}

func (m *HealthMonitor) interval() time.Duration {
	if m != nil && m.Interval > 0 {
		return m.Interval
	}
	return 30 * time.Second
}

// Start begins the background ticker. Stop must be called.
func (m *HealthMonitor) Start(ctx context.Context) {
	m.mu.Lock()
	if m.stop != nil {
		m.mu.Unlock()
		return
	}
	m.stop = make(chan struct{})
	m.done = make(chan struct{})
	m.health = map[domain.ProxyID]domain.ProxyHealth{}
	m.mu.Unlock()

	go func() {
		defer close(m.done)
		_ = m.Recheck(ctx)
		ticker := time.NewTicker(m.interval())
		defer ticker.Stop()
		for {
			select {
			case <-m.stop:
				return
			case <-ticker.C:
				_ = m.Recheck(context.Background())
			}
		}
	}()
}

// Stop ends the ticker and waits for the goroutine.
func (m *HealthMonitor) Stop() {
	m.mu.Lock()
	stop := m.stop
	done := m.done
	m.mu.Unlock()
	if stop == nil {
		return
	}
	close(stop)
	<-done
	m.mu.Lock()
	m.stop = nil
	m.done = nil
	m.mu.Unlock()
}

// Snapshot returns a copy of the health registry.
func (m *HealthMonitor) Snapshot() map[domain.ProxyID]domain.ProxyHealth {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[domain.ProxyID]domain.ProxyHealth, len(m.health))
	for id, h := range m.health {
		out[id] = h
	}
	return out
}

// Recheck probes all proxies once.
func (m *HealthMonitor) Recheck(ctx context.Context) error {
	if m.Proxies == nil || m.Checker == nil {
		return fmt.Errorf("health monitor: not configured")
	}
	proxies, err := m.Proxies.List(ctx)
	if err != nil {
		return fmt.Errorf("health recheck: %w", err)
	}
	for _, proxy := range proxies {
		var password []byte
		if m.Secrets != nil && proxy.KeychainRef != "" {
			password, _ = m.Secrets.Load(ctx, proxy.KeychainRef)
		}
		health, err := m.Checker.Check(ctx, proxy, password)
		if err != nil {
			m.logger().Info("proxy health check error", "component", "health", "proxy_id", string(proxy.ID), "err", err)
			health = domain.ProxyHealthUnknown
		}
		m.mu.Lock()
		old := m.health[proxy.ID]
		m.health[proxy.ID] = health
		m.mu.Unlock()
		if old != health {
			m.logger().Info("proxy health changed",
				"component", "health",
				"proxy_id", string(proxy.ID),
				"health", string(health),
			)
		}
	}
	return nil
}
