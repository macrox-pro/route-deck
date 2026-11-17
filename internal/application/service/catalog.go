// Package service holds application use cases over domain and ports.
package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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

// HelperGroup is one bundle-id group for Review helpers and Profiles (ADR-0014).
type HelperGroup struct {
	BundleID        string
	DisplayName     string
	ExecutableNames []string
	ProcessIDs      []domain.ProcessIdentityID
	Confirmed       bool
	Preselected     bool
	IsUpdater       bool
	SameDeveloper   bool
	Representative  port.InspectedProcess
}

// ConfirmHelpers syncs selected helper groups under applicationID (ADR-0014).
// Checked groups are upserted (one identity per bundle id; path/exec refreshed).
// Unchecked groups delete every identity with that bundle id, or return ErrReferenced
// when any profile references any of them. Policies are never deleted as a side effect.
func (c *Catalog) ConfirmHelpers(ctx context.Context, applicationID domain.ApplicationID, helpers []port.InspectedProcess) error {
	app, err := c.Apps.Get(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("catalog confirm helpers: %w", err)
	}
	existing, err := c.Processes.List(ctx)
	if err != nil {
		return fmt.Errorf("catalog confirm helpers: %w", err)
	}
	profiles, err := c.Profiles.List(ctx)
	if err != nil {
		return fmt.Errorf("catalog confirm helpers: %w", err)
	}
	referenced := make(map[domain.ProcessIdentityID]struct{})
	for _, profile := range profiles {
		for _, policy := range profile.Policies {
			referenced[policy.ProcessID] = struct{}{}
		}
	}

	selected := make(map[string]port.InspectedProcess)
	for _, helper := range helpers {
		if helper.BundleID == "" || helper.BundleID == app.BundleID {
			continue
		}
		if _, ok := selected[helper.BundleID]; ok {
			continue
		}
		selected[helper.BundleID] = helper
	}

	byBundle := make(map[string][]domain.ProcessIdentity)
	for _, identity := range existing {
		if identity.ApplicationID != applicationID || identity.BundleID == "" || identity.BundleID == app.BundleID {
			continue
		}
		byBundle[identity.BundleID] = append(byBundle[identity.BundleID], identity)
	}

	for bundleID, rows := range byBundle {
		if _, keep := selected[bundleID]; keep {
			continue
		}
		for _, row := range rows {
			if _, ok := referenced[row.ID]; ok {
				return fmt.Errorf("catalog confirm helpers: %w", domain.ErrReferenced)
			}
		}
		for _, row := range rows {
			if err := c.Processes.Delete(ctx, row.ID); err != nil {
				return fmt.Errorf("catalog confirm helpers: %w", err)
			}
		}
	}

	for bundleID, helper := range selected {
		rows := byBundle[bundleID]
		if len(rows) == 0 {
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
			continue
		}
		// Refresh the first identity; leave sibling duplicates for compile-time collapse.
		rep := rows[0]
		rep.ExecutableName = helper.ExecutableName
		rep.Path = helper.Path
		rep.BundleID = helper.BundleID
		if err := c.Processes.Save(ctx, rep); err != nil {
			return fmt.Errorf("catalog confirm helpers: %w", err)
		}
	}
	return nil
}

// HelperGroups returns Go-built helper groups for an application (ADR-0014).
// ExecutableNames come from inspect candidates. Preselected follows Closed rules.
func (c *Catalog) HelperGroups(ctx context.Context, applicationID domain.ApplicationID) ([]HelperGroup, int, error) {
	app, err := c.Apps.Get(ctx, applicationID)
	if err != nil {
		return nil, 0, fmt.Errorf("catalog helper groups: %w", err)
	}
	identities, err := c.Processes.List(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("catalog helper groups: %w", err)
	}
	bundlePath := ""
	confirmed := make(map[string][]domain.ProcessIdentityID)
	for _, identity := range identities {
		if identity.ApplicationID != applicationID {
			continue
		}
		if identity.BundleID == app.BundleID {
			if bundlePath == "" {
				bundlePath = appBundlePathFromProcess(identity.Path)
			}
			continue
		}
		if identity.BundleID != "" {
			confirmed[identity.BundleID] = append(confirmed[identity.BundleID], identity.ID)
		}
		if bundlePath == "" {
			bundlePath = appBundlePathFromProcess(identity.Path)
		}
	}
	if bundlePath == "" || c.Inspector == nil {
		return nil, 0, fmt.Errorf("catalog helper groups: no application bundle path")
	}
	inspected, err := c.Inspector.Inspect(ctx, bundlePath)
	if err != nil {
		return nil, 0, fmt.Errorf("catalog helper groups: %w", err)
	}

	type bucket struct {
		members []port.InspectedProcess
	}
	order := make([]string, 0)
	byBundle := make(map[string]*bucket)
	for _, helper := range inspected.Helpers {
		if helper.BundleID == "" {
			continue
		}
		b, ok := byBundle[helper.BundleID]
		if !ok {
			b = &bucket{}
			byBundle[helper.BundleID] = b
			order = append(order, helper.BundleID)
		}
		b.members = append(b.members, helper)
	}

	out := make([]HelperGroup, 0, len(order))
	pending := 0
	for _, bundleID := range order {
		members := byBundle[bundleID].members
		rep := members[0]
		names := make([]string, 0, len(members))
		for _, m := range members {
			if m.ExecutableName != "" {
				names = append(names, m.ExecutableName)
			}
		}
		ids := confirmed[bundleID]
		isConfirmed := len(ids) > 0
		sameDev := sameDeveloperPrefix(app.BundleID, bundleID)
		preselected := isConfirmed || (!rep.IsUpdater && sameDev)
		if !isConfirmed {
			pending++
		}
		display := rep.DisplayName
		if display == "" {
			display = rep.ExecutableName
		}
		out = append(out, HelperGroup{
			BundleID:        bundleID,
			DisplayName:     display,
			ExecutableNames: names,
			ProcessIDs:      ids,
			Confirmed:       isConfirmed,
			Preselected:     preselected,
			IsUpdater:       rep.IsUpdater,
			SameDeveloper:   sameDev,
			Representative:  rep,
		})
	}
	return out, pending, nil
}

func sameDeveloperPrefix(mainID, helperID string) bool {
	mainParts := splitBundleLabels(mainID)
	helperParts := splitBundleLabels(helperID)
	if len(mainParts) < 2 || len(helperParts) < 2 {
		return false
	}
	return mainParts[0] == helperParts[0] && mainParts[1] == helperParts[1]
}

func splitBundleLabels(id string) []string {
	if id == "" {
		return nil
	}
	parts := make([]string, 0, 4)
	start := 0
	for i := 0; i <= len(id); i++ {
		if i == len(id) || id[i] == '.' {
			if i > start {
				parts = append(parts, id[start:i])
			}
			start = i + 1
		}
	}
	return parts
}

func appBundlePathFromProcess(processPath string) string {
	p := filepath.Clean(processPath)
	for p != "." && p != "/" {
		if strings.HasSuffix(strings.ToLower(p), ".app") {
			return p
		}
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
		p = parent
	}
	return ""
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

// SaveProfile rejects unset actions, unresolved proxy refs, and match conflicts before any write.
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
	identities, err := c.Processes.List(ctx)
	if err != nil {
		return fmt.Errorf("catalog save profile: %w", err)
	}
	if err := domain.ValidateMatchConflicts(profile, identities); err != nil {
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
