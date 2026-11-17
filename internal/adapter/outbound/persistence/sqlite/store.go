package sqlite

import (
	"context"
	"errors"
	"fmt"

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
		KeychainRef: string(endpoint.KeychainRef),
		SendDomain:  endpoint.SendDomain,
	}
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("proxy save: %w", err)
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
	row := applicationModel{
		ID:          string(app.ID),
		DisplayName: app.DisplayName,
		BundleID:    app.BundleID,
		IconRef:     app.IconRef,
	}
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("application save: %w", err)
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
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := profileModel{ID: string(profile.ID), Name: profile.Name}
		if err := tx.Save(&row).Error; err != nil {
			return fmt.Errorf("profile save: %w", err)
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
	endpoint, err := domain.NewProxyEndpoint(
		domain.ProxyID(row.ID),
		row.Name,
		kind,
		row.Host,
		row.Port,
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
	row := applicationPolicyModel{
		ProfileID: string(profileID),
		ProcessID: string(policy.ProcessID),
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
	return domain.ApplicationPolicy{
		ProcessID: domain.ProcessIdentityID(row.ProcessID),
		Action:    action,
	}, nil
}
