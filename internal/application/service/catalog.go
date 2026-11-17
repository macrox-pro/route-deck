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
	Inspector port.ApplicationInspector
	Installed port.InstalledApplicationLister
	AppIcons  port.ApplicationIconReader
}

// InstalledRow is a browsable installed application with catalog membership.
type InstalledRow struct {
	DisplayName string
	BundleID    string
	Path        string
	InCatalog   bool
	IsSystem    bool
}

// ListInstalled returns platform-installed applications without writing SQLite.
func (c *Catalog) ListInstalled(ctx context.Context) ([]InstalledRow, error) {
	if c.Installed == nil {
		return nil, fmt.Errorf("catalog list installed: lister not configured")
	}
	found, err := c.Installed.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("catalog list installed: %w", err)
	}
	catalogApps, err := c.Apps.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("catalog list installed: %w", err)
	}
	inCatalog := make(map[string]struct{}, len(catalogApps))
	for _, app := range catalogApps {
		inCatalog[app.BundleID] = struct{}{}
	}
	out := make([]InstalledRow, 0, len(found))
	for _, app := range found {
		_, ok := inCatalog[app.BundleID]
		out = append(out, InstalledRow{
			DisplayName: app.DisplayName,
			BundleID:    app.BundleID,
			Path:        app.Path,
			InCatalog:   ok,
			IsSystem:    app.IsSystem,
		})
	}
	return out, nil
}

// ApplicationIconPNG returns a PNG icon for path through the icon port.
func (c *Catalog) ApplicationIconPNG(ctx context.Context, path string) ([]byte, error) {
	if c.AppIcons == nil {
		return nil, fmt.Errorf("catalog application icon: reader not configured")
	}
	png, err := c.AppIcons.PNG(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("catalog application icon: %w", err)
	}
	return png, nil
}

// AddApplicationFromPath inspects path and saves the main application + process identity.
// Helper candidates are returned for ConfirmHelpers; they are not saved yet.
func (c *Catalog) AddApplicationFromPath(ctx context.Context, path string) (domain.Application, port.InspectedApplication, error) {
	if c.Inspector == nil {
		return domain.Application{}, port.InspectedApplication{}, fmt.Errorf("catalog add application: inspector not configured")
	}
	inspected, err := c.Inspector.Inspect(ctx, path)
	if err != nil {
		return domain.Application{}, port.InspectedApplication{}, fmt.Errorf("catalog add application: %w", err)
	}
	apps, err := c.Apps.List(ctx)
	if err != nil {
		return domain.Application{}, port.InspectedApplication{}, fmt.Errorf("catalog add application: %w", err)
	}
	for _, app := range apps {
		if app.BundleID == inspected.BundleID {
			return domain.Application{}, port.InspectedApplication{}, fmt.Errorf("catalog add application: %w", domain.ErrDuplicateBundleID)
		}
	}
	appID, err := newID()
	if err != nil {
		return domain.Application{}, port.InspectedApplication{}, err
	}
	procID, err := newID()
	if err != nil {
		return domain.Application{}, port.InspectedApplication{}, err
	}
	app := domain.Application{
		ID:          domain.ApplicationID(appID),
		DisplayName: inspected.DisplayName,
		BundleID:    inspected.BundleID,
	}
	identity := domain.ProcessIdentity{
		ID:             domain.ProcessIdentityID(procID),
		ApplicationID:  app.ID,
		BundleID:       inspected.Main.BundleID,
		ExecutableName: inspected.Main.ExecutableName,
		Path:           inspected.Main.Path,
	}
	if err := c.Apps.Save(ctx, app); err != nil {
		return domain.Application{}, port.InspectedApplication{}, fmt.Errorf("catalog add application: %w", err)
	}
	if err := c.Processes.Save(ctx, identity); err != nil {
		return domain.Application{}, port.InspectedApplication{}, fmt.Errorf("catalog add application: %w", err)
	}
	return app, inspected, nil
}

// InspectPath inspects an application bundle without writing the catalog.
func (c *Catalog) InspectPath(ctx context.Context, path string) (port.InspectedApplication, error) {
	if c.Inspector == nil {
		return port.InspectedApplication{}, fmt.Errorf("catalog inspect: inspector not configured")
	}
	if path == "" {
		return port.InspectedApplication{}, fmt.Errorf("catalog inspect: empty path")
	}
	inspected, err := c.Inspector.Inspect(ctx, path)
	if err != nil {
		return port.InspectedApplication{}, fmt.Errorf("catalog inspect: %w", err)
	}
	return inspected, nil
}

// ConfirmHelpers saves selected helper process identities under applicationID.
// ProxyBridge matches bundle id / executable name (not path), so a helper whose
// bundle id is already catalogued for this application is skipped.
func (c *Catalog) ConfirmHelpers(ctx context.Context, applicationID domain.ApplicationID, helpers []port.InspectedProcess) error {
	existing, err := c.Processes.List(ctx)
	if err != nil {
		return fmt.Errorf("catalog confirm helpers: %w", err)
	}
	knownBundles := make(map[string]struct{})
	for _, identity := range existing {
		if identity.ApplicationID != applicationID || identity.BundleID == "" {
			continue
		}
		knownBundles[identity.BundleID] = struct{}{}
	}
	seenInBatch := make(map[string]struct{})
	for _, helper := range helpers {
		if helper.BundleID == "" {
			continue
		}
		if _, ok := knownBundles[helper.BundleID]; ok {
			continue
		}
		if _, ok := seenInBatch[helper.BundleID]; ok {
			continue
		}
		id, err := newID()
		if err != nil {
			return err
		}
		identity := domain.ProcessIdentity{
			ID:             domain.ProcessIdentityID(id),
			ApplicationID:  applicationID,
			BundleID:       helper.BundleID,
			ExecutableName: helper.ExecutableName,
			Path:           helper.Path,
		}
		if err := c.Processes.Save(ctx, identity); err != nil {
			return fmt.Errorf("catalog confirm helpers: %w", err)
		}
		seenInBatch[helper.BundleID] = struct{}{}
		knownBundles[helper.BundleID] = struct{}{}
	}
	return nil
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

// DeleteProxy removes a proxy when it is not referenced.
func (c *Catalog) DeleteProxy(ctx context.Context, id domain.ProxyID) error {
	if err := c.Proxies.Delete(ctx, id); err != nil {
		return fmt.Errorf("catalog delete proxy: %w", err)
	}
	return nil
}

// DeleteApplication removes an application when it is not referenced.
func (c *Catalog) DeleteApplication(ctx context.Context, id domain.ApplicationID) error {
	if err := c.Apps.Delete(ctx, id); err != nil {
		return fmt.Errorf("catalog delete application: %w", err)
	}
	return nil
}

// DeleteProfile removes a profile unless it is desired or applied.
func (c *Catalog) DeleteProfile(ctx context.Context, id domain.ProfileID, desiredID, appliedID domain.ProfileID) error {
	if id == desiredID || id == appliedID {
		return fmt.Errorf("catalog delete profile: %w", domain.ErrReferenced)
	}
	if err := c.Profiles.Delete(ctx, id); err != nil {
		return fmt.Errorf("catalog delete profile: %w", err)
	}
	return nil
}

// CopyProfile duplicates a profile under a new unique name.
func (c *Catalog) CopyProfile(ctx context.Context, sourceID domain.ProfileID, newName string) (domain.Profile, error) {
	source, err := c.Profiles.Get(ctx, sourceID)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("catalog copy profile: %w", err)
	}
	if err := c.ensureUniqueProfileName(ctx, newName, ""); err != nil {
		return domain.Profile{}, err
	}
	id, err := newID()
	if err != nil {
		return domain.Profile{}, err
	}
	copyProfile, err := domain.NewProfile(domain.ProfileID(id), newName, append([]domain.ApplicationPolicy(nil), source.Policies...))
	if err != nil {
		return domain.Profile{}, fmt.Errorf("catalog copy profile: %w", err)
	}
	if err := c.SaveProfile(ctx, copyProfile); err != nil {
		return domain.Profile{}, err
	}
	return copyProfile, nil
}

// RenameProfile changes a profile name with uniqueness check.
func (c *Catalog) RenameProfile(ctx context.Context, id domain.ProfileID, newName string) error {
	profile, err := c.Profiles.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("catalog rename profile: %w", err)
	}
	if err := c.ensureUniqueProfileName(ctx, newName, id); err != nil {
		return err
	}
	profile.Name = newName
	if err := c.SaveProfile(ctx, profile); err != nil {
		return err
	}
	return nil
}

func (c *Catalog) ensureUniqueProfileName(ctx context.Context, name string, except domain.ProfileID) error {
	if name == "" {
		return fmt.Errorf("catalog profile name: %w", domain.ErrEmptyName)
	}
	profiles, err := c.Profiles.List(ctx)
	if err != nil {
		return fmt.Errorf("catalog profile name: %w", err)
	}
	for _, profile := range profiles {
		if profile.ID == except {
			continue
		}
		if profile.Name == name {
			return fmt.Errorf("catalog profile name: %w", domain.ErrDuplicateProfileName)
		}
	}
	return nil
}
