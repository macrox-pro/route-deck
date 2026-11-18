package sqlite

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
	"gorm.io/gorm"
)

// Store implements catalog ports against SQLite.
type Store struct {
	db *gorm.DB
}

var (
	_ port.ProxyStore           = (*Store)(nil)
	_ port.ApplicationStore     = applicationStore{}
	_ port.ProcessIdentityStore = processStore{}
	_ port.ProfileStore         = profileStore{}
)

// Get loads one proxy endpoint by id.
func (s *Store) Get(ctx context.Context, id domain.ProxyID) (domain.ProxyEndpoint, error) {
	var row proxyModel
	err := s.db.WithContext(ctx).Where("id = ?", string(id)).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ProxyEndpoint{}, fmt.Errorf("proxy get: %w", ErrNotFound)
	}
	if err != nil {
		return domain.ProxyEndpoint{}, fmt.Errorf("proxy get: %w", err)
	}
	return proxyToDomain(row)
}

// List returns all proxy endpoints.
func (s *Store) List(ctx context.Context) ([]domain.ProxyEndpoint, error) {
	var rows []proxyModel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("proxy list: %w", err)
	}
	out := make([]domain.ProxyEndpoint, 0, len(rows))
	for _, row := range rows {
		endpoint, err := proxyToDomain(row)
		if err != nil {
			return nil, err
		}
		out = append(out, endpoint)
	}
	return out, nil
}

// Save upserts a proxy endpoint. Passwords are never written.
func (s *Store) Save(ctx context.Context, endpoint domain.ProxyEndpoint) error {
	row := proxyModel{
		ID:          string(endpoint.ID),
		Name:        endpoint.Name,
		Kind:        string(endpoint.Kind),
		Host:        endpoint.Host,
		Port:        endpoint.Port,
		Username:    endpoint.Username,
		KeychainRef: string(endpoint.KeychainRef),
		SendDomain:  endpoint.SendDomain,
	}
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("proxy save: %w", err)
	}
	return nil
}

// Delete removes a proxy by id when no policy references it.
func (s *Store) Delete(ctx context.Context, id domain.ProxyID) error {
	var count int64
	if err := s.db.WithContext(ctx).Model(&applicationPolicyModel{}).
		Where("proxy_id = ?", string(id)).Count(&count).Error; err != nil {
		return fmt.Errorf("proxy delete: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("proxy delete: %w", domain.ErrReferenced)
	}
	res := s.db.WithContext(ctx).Where("id = ?", string(id)).Delete(&proxyModel{})
	if res.Error != nil {
		return fmt.Errorf("proxy delete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("proxy delete: %w", ErrNotFound)
	}
	return nil
}

// GetApplication loads one application by id.
func (s *Store) GetApplication(ctx context.Context, id domain.ApplicationID) (domain.Application, error) {
	var row applicationModel
	err := s.db.WithContext(ctx).Where("id = ?", string(id)).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Application{}, fmt.Errorf("application get: %w", ErrNotFound)
	}
	if err != nil {
		return domain.Application{}, fmt.Errorf("application get: %w", err)
	}
	return applicationToDomain(row), nil
}

// ListApplications lists all applications.
func (s *Store) ListApplications(ctx context.Context) ([]domain.Application, error) {
	var rows []applicationModel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("application list: %w", err)
	}
	out := make([]domain.Application, 0, len(rows))
	for _, row := range rows {
		out = append(out, applicationToDomain(row))
	}
	return out, nil
}

// SaveApplication upserts an application.
func (s *Store) SaveApplication(ctx context.Context, app domain.Application) error {
	var existing applicationModel
	err := s.db.WithContext(ctx).Where("bundle_id = ? AND id <> ?", app.BundleID, string(app.ID)).First(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// No other application owns this bundle id.
	case err != nil:
		return fmt.Errorf("application save: %w", err)
	default:
		return fmt.Errorf("application save: %w", domain.ErrDuplicateBundleID)
	}
	row := applicationModel{
		ID:          string(app.ID),
		DisplayName: app.DisplayName,
		BundleID:    app.BundleID,
		IconRef:     app.IconRef,
	}
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("application save: %w", mapUniqueError(err, domain.ErrDuplicateBundleID))
	}
	return nil
}

// GetProcessIdentity loads one process identity by id.
func (s *Store) GetProcessIdentity(ctx context.Context, id domain.ProcessIdentityID) (domain.ProcessIdentity, error) {
	var row processIdentityModel
	err := s.db.WithContext(ctx).Where("id = ?", string(id)).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ProcessIdentity{}, fmt.Errorf("process identity get: %w", ErrNotFound)
	}
	if err != nil {
		return domain.ProcessIdentity{}, fmt.Errorf("process identity get: %w", err)
	}
	return processToDomain(row), nil
}

// ListProcessIdentities lists all process identities.
func (s *Store) ListProcessIdentities(ctx context.Context) ([]domain.ProcessIdentity, error) {
	var rows []processIdentityModel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("process identity list: %w", err)
	}
	out := make([]domain.ProcessIdentity, 0, len(rows))
	for _, row := range rows {
		out = append(out, processToDomain(row))
	}
	return out, nil
}

// SaveProcessIdentity upserts a process identity.
func (s *Store) SaveProcessIdentity(ctx context.Context, identity domain.ProcessIdentity) error {
	row := processIdentityModel{
		ID:             string(identity.ID),
		ApplicationID:  string(identity.ApplicationID),
		BundleID:       identity.BundleID,
		ExecutableName: identity.ExecutableName,
		Path:           identity.Path,
		SigningHint:    identity.SigningHint,
	}
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("process identity save: %w", err)
	}
	return nil
}

// GetProfile loads one profile with its policies.
func (s *Store) GetProfile(ctx context.Context, id domain.ProfileID) (domain.Profile, error) {
	var row profileModel
	err := s.db.WithContext(ctx).Where("id = ?", string(id)).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Profile{}, fmt.Errorf("profile get: %w", ErrNotFound)
	}
	if err != nil {
		return domain.Profile{}, fmt.Errorf("profile get: %w", err)
	}
	policies, err := s.loadPolicies(ctx, id)
	if err != nil {
		return domain.Profile{}, err
	}
	return domain.Profile{ID: domain.ProfileID(row.ID), Name: row.Name, Policies: policies}, nil
}

// ListProfiles lists all profiles with policies.
func (s *Store) ListProfiles(ctx context.Context) ([]domain.Profile, error) {
	var rows []profileModel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("profile list: %w", err)
	}
	out := make([]domain.Profile, 0, len(rows))
	for _, row := range rows {
		policies, err := s.loadPolicies(ctx, domain.ProfileID(row.ID))
		if err != nil {
			return nil, err
		}
		out = append(out, domain.Profile{ID: domain.ProfileID(row.ID), Name: row.Name, Policies: policies})
	}
	return out, nil
}

// SaveProfile upserts a profile and replaces its policies in one transaction.
func (s *Store) SaveProfile(ctx context.Context, profile domain.Profile) error {
	for _, policy := range profile.Policies {
		if !policy.Action.IsSet() {
			return fmt.Errorf("profile save: %w", domain.ErrUnsetAction)
		}
	}
	var existing profileModel
	err := s.db.WithContext(ctx).Where("name = ? AND id <> ?", profile.Name, string(profile.ID)).First(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// No other profile uses this name.
	case err != nil:
		return fmt.Errorf("profile save: %w", err)
	default:
		return fmt.Errorf("profile save: %w", domain.ErrDuplicateProfileName)
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := profileModel{ID: string(profile.ID), Name: profile.Name}
		if err := tx.Save(&row).Error; err != nil {
			return fmt.Errorf("profile save: %w", mapUniqueError(err, domain.ErrDuplicateProfileName))
		}
		if err := tx.Where("profile_id = ?", string(profile.ID)).Delete(&applicationPolicyModel{}).Error; err != nil {
			return fmt.Errorf("profile save delete policies: %w", err)
		}
		for _, policy := range profile.Policies {
			policyRow, err := policyToModel(profile.ID, policy)
			if err != nil {
				return err
			}
			if err := tx.Create(&policyRow).Error; err != nil {
				return fmt.Errorf("profile save policy: %w", err)
			}
		}
		return nil
	})
}

func (s *Store) loadPolicies(ctx context.Context, profileID domain.ProfileID) ([]domain.ApplicationPolicy, error) {
	var rows []applicationPolicyModel
	if err := s.db.WithContext(ctx).Where("profile_id = ?", string(profileID)).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("profile policies: %w", err)
	}
	out := make([]domain.ApplicationPolicy, 0, len(rows))
	for _, row := range rows {
		policy, err := policyToDomain(row)
		if err != nil {
			return nil, err
		}
		out = append(out, policy)
	}
	return out, nil
}

// Port-facing method names required by ApplicationStore / ProcessIdentityStore / ProfileStore.
// Go cannot overload Get/List/Save across interfaces on one type, so adapters wrap below.

// applicationStore adapts Store to port.ApplicationStore.
type applicationStore struct{ *Store }

func (a applicationStore) Get(ctx context.Context, id domain.ApplicationID) (domain.Application, error) {
	return a.GetApplication(ctx, id)
}
func (a applicationStore) List(ctx context.Context) ([]domain.Application, error) {
	return a.ListApplications(ctx)
}
func (a applicationStore) Save(ctx context.Context, app domain.Application) error {
	return a.SaveApplication(ctx, app)
}
func (a applicationStore) Delete(ctx context.Context, id domain.ApplicationID) error {
	var processIDs []string
	if err := a.db.WithContext(ctx).Model(&processIdentityModel{}).
		Where("application_id = ?", string(id)).Pluck("id", &processIDs).Error; err != nil {
		return fmt.Errorf("application delete: %w", err)
	}
	if len(processIDs) > 0 {
		var count int64
		if err := a.db.WithContext(ctx).Model(&applicationPolicyModel{}).
			Where("process_id IN ?", processIDs).Count(&count).Error; err != nil {
			return fmt.Errorf("application delete: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("application delete: %w", domain.ErrReferenced)
		}
	}
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("application_id = ?", string(id)).Delete(&processIdentityModel{}).Error; err != nil {
			return fmt.Errorf("application delete processes: %w", err)
		}
		res := tx.Where("id = ?", string(id)).Delete(&applicationModel{})
		if res.Error != nil {
			return fmt.Errorf("application delete: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("application delete: %w", ErrNotFound)
		}
		return nil
	})
}

// processStore adapts Store to port.ProcessIdentityStore.
type processStore struct{ *Store }

func (p processStore) Get(ctx context.Context, id domain.ProcessIdentityID) (domain.ProcessIdentity, error) {
	return p.GetProcessIdentity(ctx, id)
}
func (p processStore) List(ctx context.Context) ([]domain.ProcessIdentity, error) {
	return p.ListProcessIdentities(ctx)
}
func (p processStore) Save(ctx context.Context, identity domain.ProcessIdentity) error {
	return p.SaveProcessIdentity(ctx, identity)
}
func (p processStore) Delete(ctx context.Context, id domain.ProcessIdentityID) error {
	res := p.db.WithContext(ctx).Where("id = ?", string(id)).Delete(&processIdentityModel{})
	if res.Error != nil {
		return fmt.Errorf("process identity delete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("process identity delete: %w", ErrNotFound)
	}
	return nil
}

// profileStore adapts Store to port.ProfileStore.
type profileStore struct{ *Store }

func (p profileStore) Get(ctx context.Context, id domain.ProfileID) (domain.Profile, error) {
	return p.GetProfile(ctx, id)
}
func (p profileStore) List(ctx context.Context) ([]domain.Profile, error) {
	return p.ListProfiles(ctx)
}
func (p profileStore) Save(ctx context.Context, profile domain.Profile) error {
	return p.SaveProfile(ctx, profile)
}
func (p profileStore) Delete(ctx context.Context, id domain.ProfileID) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("profile_id = ?", string(id)).Delete(&applicationPolicyModel{}).Error; err != nil {
			return fmt.Errorf("profile delete policies: %w", err)
		}
		res := tx.Where("id = ?", string(id)).Delete(&profileModel{})
		if res.Error != nil {
			return fmt.Errorf("profile delete: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("profile delete: %w", ErrNotFound)
		}
		return nil
	})
}

// Applications returns the ApplicationStore view of this Store.
func (s *Store) Applications() port.ApplicationStore { return applicationStore{s} }

// ProcessIdentities returns the ProcessIdentityStore view of this Store.
func (s *Store) ProcessIdentities() port.ProcessIdentityStore { return processStore{s} }

// Profiles returns the ProfileStore view of this Store.
func (s *Store) Profiles() port.ProfileStore { return profileStore{s} }

func applicationToDomain(row applicationModel) domain.Application {
	return domain.Application{
		ID:          domain.ApplicationID(row.ID),
		DisplayName: row.DisplayName,
		BundleID:    row.BundleID,
		IconRef:     row.IconRef,
	}
}

func processToDomain(row processIdentityModel) domain.ProcessIdentity {
	return domain.ProcessIdentity{
		ID:             domain.ProcessIdentityID(row.ID),
		ApplicationID:  domain.ApplicationID(row.ApplicationID),
		BundleID:       row.BundleID,
		ExecutableName: row.ExecutableName,
		Path:           row.Path,
		SigningHint:    row.SigningHint,
	}
}

func proxyToDomain(row proxyModel) (domain.ProxyEndpoint, error) {
	kind, err := domain.ParseProxyKind(row.Kind)
	if err != nil {
		return domain.ProxyEndpoint{}, fmt.Errorf("proxy map: %w", err)
	}
	endpoint, err := domain.NewProxyEndpointWithUsername(
		domain.ProxyID(row.ID),
		row.Name,
		kind,
		row.Host,
		row.Port,
		row.Username,
		domain.KeychainRef(row.KeychainRef),
	)
	if err != nil {
		return domain.ProxyEndpoint{}, fmt.Errorf("proxy map: %w", err)
	}
	endpoint.SendDomain = row.SendDomain
	return endpoint, nil
}

func policyToModel(profileID domain.ProfileID, policy domain.ApplicationPolicy) (applicationPolicyModel, error) {
	if !policy.Action.IsSet() {
		return applicationPolicyModel{}, fmt.Errorf("policy map: %w", domain.ErrUnsetAction)
	}
	limitation := policy.UDPLimitation
	if limitation == "" {
		limitation = domain.UDPLimitationAbsent
	}
	row := applicationPolicyModel{
		ProfileID:     string(profileID),
		ProcessID:     string(policy.ProcessID),
		ProtocolScope: string(policy.ProtocolScope),
		UDPLimitation: string(limitation),
	}
	if policy.Action.IsDirect() {
		row.ActionKind = "direct"
		return row, nil
	}
	if policy.Action.IsBlock() {
		row.ActionKind = "block"
		return row, nil
	}
	proxyID, ok := policy.Action.ProxyID()
	if !ok {
		return applicationPolicyModel{}, fmt.Errorf("policy map: %w", domain.ErrUnsetAction)
	}
	id := string(proxyID)
	row.ActionKind = "proxy"
	row.ProxyID = &id
	return row, nil
}

func policyToDomain(row applicationPolicyModel) (domain.ApplicationPolicy, error) {
	var action domain.RoutingAction
	switch row.ActionKind {
	case "direct":
		action = domain.DirectAction()
	case "block":
		action = domain.BlockAction()
	case "proxy":
		if row.ProxyID == nil || *row.ProxyID == "" {
			return domain.ApplicationPolicy{}, fmt.Errorf("policy map: proxy action missing proxy id")
		}
		var err error
		action, err = domain.ProxyAction(domain.ProxyID(*row.ProxyID))
		if err != nil {
			return domain.ApplicationPolicy{}, fmt.Errorf("policy map: %w", err)
		}
	default:
		return domain.ApplicationPolicy{}, fmt.Errorf("policy map: unknown action kind %q", row.ActionKind)
	}
	limitation := domain.UDPLimitation(row.UDPLimitation)
	if limitation == "" {
		limitation = domain.UDPLimitationAbsent
	}
	return domain.ApplicationPolicy{
		ProcessID:     domain.ProcessIdentityID(row.ProcessID),
		Action:        action,
		ProtocolScope: domain.ProtocolScope(row.ProtocolScope),
		UDPLimitation: limitation,
	}, nil
}

func mapUniqueError(err error, sentinel error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(strings.ToUpper(err.Error()), "UNIQUE") {
		return sentinel
	}
	return err
}

// coerceSettingsAppearance maps empty/unknown appearance to system.
func coerceSettingsAppearance(raw string) string {
	switch raw {
	case "light":
		return "light"
	case "dark":
		return "dark"
	default:
		return "system"
	}
}

// GetSettings loads the single settings row (defaults when missing).
// Invalid appearance is coerced to system and rewritten into SQLite.
func (s *Store) GetSettings(ctx context.Context) (port.Settings, error) {
	var row settingsModel
	err := s.db.WithContext(ctx).Where("id = ?", 1).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		defaults := port.Settings{Appearance: "system", AutoApply: true}
		if saveErr := s.SaveSettings(ctx, defaults); saveErr != nil {
			return defaults, nil
		}
		return defaults, nil
	}
	if err != nil {
		return port.Settings{}, fmt.Errorf("settings get: %w", err)
	}
	settings := port.Settings{
		DesiredProfileID: domain.ProfileID(row.DesiredProfileID),
		DesiredSet:       row.DesiredSet,
		ProxyBridgePath:  row.ProxyBridgePath,
		DebugLogging:     row.DebugLogging,
		Appearance:       row.Appearance,
		AutoApply:        row.AutoApply,
	}
	if settings.DesiredSet && settings.DesiredProfileID == "" {
		settings.DesiredSet = false
	}
	appearance := coerceSettingsAppearance(settings.Appearance)
	needsRewrite := settings.Appearance != appearance || (row.DesiredSet && row.DesiredProfileID == "")
	settings.Appearance = appearance
	if needsRewrite {
		if saveErr := s.SaveSettings(ctx, settings); saveErr != nil {
			// Return coerced settings even when rewrite fails.
			return settings, nil
		}
	}
	return settings, nil
}

// SaveSettings upserts the single settings row.
func (s *Store) SaveSettings(ctx context.Context, settings port.Settings) error {
	if settings.DesiredSet && settings.DesiredProfileID == "" {
		settings.DesiredSet = false
	}
	settings.Appearance = coerceSettingsAppearance(settings.Appearance)
	row := settingsModel{
		ID:               1,
		DesiredProfileID: string(settings.DesiredProfileID),
		DesiredSet:       settings.DesiredSet,
		ProxyBridgePath:  settings.ProxyBridgePath,
		DebugLogging:     settings.DebugLogging,
		Appearance:       settings.Appearance,
		AutoApply:        settings.AutoApply,
	}
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("settings save: %w", err)
	}
	return nil
}

type settingsStore struct{ *Store }

func (s settingsStore) Get(ctx context.Context) (port.Settings, error) {
	return s.GetSettings(ctx)
}

func (s settingsStore) Save(ctx context.Context, settings port.Settings) error {
	return s.SaveSettings(ctx, settings)
}

// Settings returns the SettingsStore view of this Store.
func (s *Store) Settings() port.SettingsStore { return settingsStore{s} }

var _ port.SettingsStore = settingsStore{}
