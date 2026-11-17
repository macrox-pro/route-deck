// Package service holds application use cases over domain and ports.
package service

import (
	"context"
	"fmt"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// Catalog saves and lists catalog aggregates through ports.
type Catalog struct {
	Proxies   port.ProxyStore
	Apps      port.ApplicationStore
	Processes port.ProcessIdentityStore
	Profiles  port.ProfileStore
}

// SaveApplication upserts an application.
func (c *Catalog) SaveApplication(ctx context.Context, app domain.Application) error {
	if err := c.Apps.Save(ctx, app); err != nil {
		return fmt.Errorf("catalog save application: %w", err)
	}
	return nil
}

// ListApplications returns all applications.
func (c *Catalog) ListApplications(ctx context.Context) ([]domain.Application, error) {
	apps, err := c.Apps.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("catalog list applications: %w", err)
	}
	return apps, nil
}

// SaveProcessIdentity upserts a process identity.
func (c *Catalog) SaveProcessIdentity(ctx context.Context, identity domain.ProcessIdentity) error {
	if err := c.Processes.Save(ctx, identity); err != nil {
		return fmt.Errorf("catalog save process identity: %w", err)
	}
	return nil
}

// ListProcessIdentities returns all process identities.
func (c *Catalog) ListProcessIdentities(ctx context.Context) ([]domain.ProcessIdentity, error) {
	identities, err := c.Processes.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("catalog list process identities: %w", err)
	}
	return identities, nil
}

// SaveProxy upserts a proxy endpoint through NewProxyEndpoint validation at the caller.
func (c *Catalog) SaveProxy(ctx context.Context, endpoint domain.ProxyEndpoint) error {
	if err := c.Proxies.Save(ctx, endpoint); err != nil {
		return fmt.Errorf("catalog save proxy: %w", err)
	}
	return nil
}

// ListProxies returns all proxy endpoints.
func (c *Catalog) ListProxies(ctx context.Context) ([]domain.ProxyEndpoint, error) {
	proxies, err := c.Proxies.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("catalog list proxies: %w", err)
	}
	return proxies, nil
}

// SaveProfile rejects unset actions and unresolved proxy refs before any write.
func (c *Catalog) SaveProfile(ctx context.Context, profile domain.Profile) error {
	for _, policy := range profile.Policies {
		if !policy.Action.IsSet() {
			return fmt.Errorf("catalog save profile: %w", domain.ErrUnsetAction)
		}
	}
	proxies, err := c.Proxies.List(ctx)
	if err != nil {
		return fmt.Errorf("catalog save profile: %w", err)
	}
	known := make([]domain.ProxyID, 0, len(proxies))
	for _, proxy := range proxies {
		known = append(known, proxy.ID)
	}
	if err := domain.ValidateProxyRefs(profile, known); err != nil {
		return fmt.Errorf("catalog save profile: %w", err)
	}
	if err := c.Profiles.Save(ctx, profile); err != nil {
		return fmt.Errorf("catalog save profile: %w", err)
	}
	return nil
}

// ListProfiles returns all profiles.
func (c *Catalog) ListProfiles(ctx context.Context) ([]domain.Profile, error) {
	profiles, err := c.Profiles.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("catalog list profiles: %w", err)
	}
	return profiles, nil
}
