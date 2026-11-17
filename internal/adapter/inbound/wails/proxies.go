package wails

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProxiesService exposes proxy catalog CRUD over Wails.
type ProxiesService struct {
	Catalog  *service.Catalog
	Secrets  port.SecretStore
	Health   *service.HealthMonitor
	Settings port.SettingsStore
	Emit     Emitter
	Logger   *slog.Logger
}

// ProxyDTO is a catalog proxy without secrets.
type ProxyDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	HasPassword bool   `json:"hasPassword"`
}

// ProxySaveDTO is the save input. Empty Password keeps the stored secret on edit.
type ProxySaveDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *ProxiesService) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// List returns proxy catalog rows.
func (s *ProxiesService) List(ctx context.Context) ([]ProxyDTO, error) {
	if s == nil || s.Catalog == nil {
		return nil, fmt.Errorf("proxies service: not configured")
	}
	proxies, err := s.Catalog.ListProxies(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]ProxyDTO, 0, len(proxies))
	for _, proxy := range proxies {
		hasPassword := false
		if s.Secrets != nil && proxy.KeychainRef != "" {
			secret, loadErr := s.Secrets.Load(ctx, proxy.KeychainRef)
			hasPassword = loadErr == nil && len(secret) > 0
		}
		out = append(out, ProxyDTO{
			ID:          string(proxy.ID),
			Name:        proxy.Name,
			Kind:        string(proxy.Kind),
			Host:        proxy.Host,
			Port:        proxy.Port,
			Username:    proxy.Username,
			HasPassword: hasPassword,
		})
	}
	return out, nil
}

// Save upserts a proxy. Empty password on existing id keeps Keychain secret.
func (s *ProxiesService) Save(ctx context.Context, input ProxySaveDTO) (ProxyDTO, error) {
	if s == nil || s.Catalog == nil {
		return ProxyDTO{}, fmt.Errorf("proxies service: not configured")
	}
	kind, err := domain.ParseProxyKind(input.Kind)
	if err != nil {
		return ProxyDTO{}, fmt.Errorf("proxies save: %w", err)
	}
	id := input.ID
	isNew := id == ""
	if isNew {
		generated, genErr := service.NewID()
		if genErr != nil {
			return ProxyDTO{}, genErr
		}
		id = generated
	}
	ref := domain.KeychainRef(id)
	endpoint, err := domain.NewProxyEndpointWithUsername(
		domain.ProxyID(id),
		input.Name,
		kind,
		input.Host,
		input.Port,
		input.Username,
		ref,
	)
	if err != nil {
		return ProxyDTO{}, fmt.Errorf("proxies save: %w", err)
	}
	if input.Password != "" {
		if s.Secrets == nil {
			return ProxyDTO{}, fmt.Errorf("proxies save: secret store not configured")
		}
		if err := s.Secrets.Save(ctx, ref, []byte(input.Password)); err != nil {
			return ProxyDTO{}, fmt.Errorf("proxies save secret: %w", err)
		}
	}
	if err := s.Catalog.SaveProxy(ctx, endpoint); err != nil {
		s.logger().Info("proxy save failed", "component", "wails.proxies", "err", err)
		return ProxyDTO{}, err
	}
	if s.Health != nil {
		_ = s.Health.Recheck(ctx)
		emit(s.Emit, EventProxyHealthChanged, map[string]any{"proxy_id": id})
	}
	hasPassword := input.Password != ""
	if !hasPassword && s.Secrets != nil {
		secret, loadErr := s.Secrets.Load(ctx, ref)
		hasPassword = loadErr == nil && len(secret) > 0
	}
	return ProxyDTO{
		ID:          id,
		Name:        endpoint.Name,
		Kind:        string(endpoint.Kind),
		Host:        endpoint.Host,
		Port:        endpoint.Port,
		Username:    endpoint.Username,
		HasPassword: hasPassword,
	}, nil
}

// Delete removes a proxy when unreferenced.
func (s *ProxiesService) Delete(ctx context.Context, id string) error {
	if s == nil || s.Catalog == nil {
		return fmt.Errorf("proxies service: not configured")
	}
	proxyID, err := domain.ParseProxyID(id)
	if err != nil {
		return fmt.Errorf("proxies delete: %w", err)
	}
	if err := s.Catalog.DeleteProxy(ctx, proxyID); err != nil {
		s.logger().Info("proxy delete failed", "component", "wails.proxies", "err", err)
		return err
	}
	if s.Secrets != nil {
		_ = s.Secrets.Delete(ctx, domain.KeychainRef(id))
	}
	return nil
}

// ClearPassword removes the Keychain secret for a proxy.
func (s *ProxiesService) ClearPassword(ctx context.Context, id string) error {
	if s == nil || s.Secrets == nil {
		return fmt.Errorf("proxies clear password: not configured")
	}
	if _, err := domain.ParseProxyID(id); err != nil {
		return fmt.Errorf("proxies clear password: %w", err)
	}
	if err := s.Secrets.Delete(ctx, domain.KeychainRef(id)); err != nil {
		return fmt.Errorf("proxies clear password: %w", err)
	}
	return nil
}

// Recheck runs a health probe for all proxies.
func (s *ProxiesService) Recheck(ctx context.Context) error {
	if s == nil || s.Health == nil {
		return fmt.Errorf("proxies recheck: health monitor not configured")
	}
	if err := s.Health.Recheck(ctx); err != nil {
		return err
	}
	emit(s.Emit, EventProxyHealthChanged, map[string]any{"scope": "all"})
	return nil
}
