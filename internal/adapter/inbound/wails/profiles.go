package wails

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProfilesService exposes profile CRUD and desired selection.
type ProfilesService struct {
	Catalog    *service.Catalog
	Settings   port.SettingsStore
	Protection *service.Protection
	Logger     *slog.Logger
}

// PolicyDTO is one application policy row.
type PolicyDTO struct {
	ProcessID     string `json:"processId"`
	Action        string `json:"action"`
	ProxyID       string `json:"proxyId"`
	ProtocolScope string `json:"protocolScope"`
	UDPLimitation string `json:"udpLimitation"`
}

// ProfileDTO is a profile with policies.
type ProfileDTO struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Policies []PolicyDTO `json:"policies"`
}

// ProfileSaveDTO is create/update input.
type ProfileSaveDTO struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Policies []PolicyDTO `json:"policies"`
}

func (s *ProfilesService) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// List returns all profiles.
func (s *ProfilesService) List(ctx context.Context) ([]ProfileDTO, error) {
	if s == nil || s.Catalog == nil {
		return nil, fmt.Errorf("profiles service: not configured")
	}
	profiles, err := s.Catalog.ListProfiles(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]ProfileDTO, 0, len(profiles))
	for _, profile := range profiles {
		out = append(out, mapProfile(profile))
	}
	return out, nil
}

// Save creates or updates a profile.
func (s *ProfilesService) Save(ctx context.Context, input ProfileSaveDTO) (ProfileDTO, error) {
	if s == nil || s.Catalog == nil {
		return ProfileDTO{}, fmt.Errorf("profiles service: not configured")
	}
	id := input.ID
	if id == "" {
		generated, err := service.NewID()
		if err != nil {
			return ProfileDTO{}, err
		}
		id = generated
	}
	policies, err := mapPoliciesFromDTO(input.Policies)
	if err != nil {
		return ProfileDTO{}, fmt.Errorf("profiles save: %w", err)
	}
	profile, err := domain.NewProfile(domain.ProfileID(id), input.Name, policies)
	if err != nil {
		return ProfileDTO{}, fmt.Errorf("profiles save: %w", err)
	}
	profiles, listErr := s.Catalog.ListProfiles(ctx)
	if listErr != nil {
		return ProfileDTO{}, listErr
	}
	for _, existing := range profiles {
		if existing.ID == profile.ID {
			continue
		}
		if existing.Name == profile.Name {
			return ProfileDTO{}, fmt.Errorf("profiles save: %w", domain.ErrDuplicateProfileName)
		}
	}
	if err := s.Catalog.SaveProfile(ctx, profile); err != nil {
		s.logger().Info("profile save failed", "component", "wails.profiles", "err", err)
		return ProfileDTO{}, err
	}
	if s.Protection != nil {
		appliedID, confirmed := s.Protection.Applied()
		if confirmed && appliedID == profile.ID {
			s.Protection.MarkDesiredDiffers()
		}
	}
	return mapProfile(profile), nil
}

// Copy duplicates a profile.
func (s *ProfilesService) Copy(ctx context.Context, sourceID, newName string) (ProfileDTO, error) {
	if s == nil || s.Catalog == nil {
		return ProfileDTO{}, fmt.Errorf("profiles service: not configured")
	}
	id, err := domain.ParseProfileID(sourceID)
	if err != nil {
		return ProfileDTO{}, fmt.Errorf("profiles copy: %w", err)
	}
	copied, err := s.Catalog.CopyProfile(ctx, id, newName)
	if err != nil {
		return ProfileDTO{}, err
	}
	return mapProfile(copied), nil
}

// Rename renames a profile.
func (s *ProfilesService) Rename(ctx context.Context, id, newName string) error {
	if s == nil || s.Catalog == nil {
		return fmt.Errorf("profiles service: not configured")
	}
	profileID, err := domain.ParseProfileID(id)
	if err != nil {
		return fmt.Errorf("profiles rename: %w", err)
	}
	return s.Catalog.RenameProfile(ctx, profileID, newName)
}

// Delete removes a profile that is not desired or applied.
func (s *ProfilesService) Delete(ctx context.Context, id string) error {
	if s == nil || s.Catalog == nil || s.Settings == nil {
		return fmt.Errorf("profiles service: not configured")
	}
	profileID, err := domain.ParseProfileID(id)
	if err != nil {
		return fmt.Errorf("profiles delete: %w", err)
	}
	settings, err := s.Settings.Get(ctx)
	if err != nil {
		return fmt.Errorf("profiles delete: %w", err)
	}
	var appliedID domain.ProfileID
	if s.Protection != nil {
		appliedID, _ = s.Protection.Applied()
	}
	desiredID := domain.ProfileID("")
	if settings.DesiredSet {
		desiredID = settings.DesiredProfileID
	}
	if err := s.Catalog.DeleteProfile(ctx, profileID, desiredID, appliedID); err != nil {
		s.logger().Info("profile delete failed", "component", "wails.profiles", "err", err)
		return err
	}
	return nil
}

// SetDesired marks a profile as desired.
func (s *ProfilesService) SetDesired(ctx context.Context, id string) error {
	if s == nil || s.Settings == nil {
		return fmt.Errorf("profiles set desired: not configured")
	}
	profileID, err := domain.ParseProfileID(id)
	if err != nil {
		return fmt.Errorf("profiles set desired: %w", err)
	}
	if s.Catalog != nil {
		if _, getErr := s.Catalog.ListProfiles(ctx); getErr != nil {
			return getErr
		}
		profiles, listErr := s.Catalog.ListProfiles(ctx)
		if listErr != nil {
			return listErr
		}
		found := false
		for _, profile := range profiles {
			if profile.ID == profileID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("profiles set desired: profile not found")
		}
	}
	settings, err := s.Settings.Get(ctx)
	if err != nil {
		return fmt.Errorf("profiles set desired: %w", err)
	}
	settings.DesiredProfileID = profileID
	settings.DesiredSet = true
	if err := s.Settings.Save(ctx, settings); err != nil {
		return fmt.Errorf("profiles set desired: %w", err)
	}
	if s.Protection != nil {
		appliedID, confirmed := s.Protection.Applied()
		if !confirmed || appliedID != profileID {
			s.Protection.MarkDesiredDiffers()
		}
	}
	return nil
}

func mapProfile(profile domain.Profile) ProfileDTO {
	policies := make([]PolicyDTO, 0, len(profile.Policies))
	for _, policy := range profile.Policies {
		dto := PolicyDTO{
			ProcessID:     string(policy.ProcessID),
			ProtocolScope: string(policy.ProtocolScope),
			UDPLimitation: string(policy.UDPLimitation),
		}
		if policy.Action.IsDirect() {
			dto.Action = "direct"
		} else if policy.Action.IsBlock() {
			dto.Action = "block"
		} else if proxyID, ok := policy.Action.ProxyID(); ok {
			dto.Action = "proxy"
			dto.ProxyID = string(proxyID)
		}
		policies = append(policies, dto)
	}
	return ProfileDTO{
		ID:       string(profile.ID),
		Name:     profile.Name,
		Policies: policies,
	}
}

func mapPoliciesFromDTO(rows []PolicyDTO) ([]domain.ApplicationPolicy, error) {
	out := make([]domain.ApplicationPolicy, 0, len(rows))
	for _, row := range rows {
		processID, err := domain.ParseProcessIdentityID(row.ProcessID)
		if err != nil {
			return nil, err
		}
		var action domain.RoutingAction
		switch row.Action {
		case "direct":
			action = domain.DirectAction()
		case "block":
			action = domain.BlockAction()
		case "proxy":
			proxyID, parseErr := domain.ParseProxyID(row.ProxyID)
			if parseErr != nil {
				return nil, parseErr
			}
			action, err = domain.ProxyAction(proxyID)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown action %q: %w", row.Action, domain.ErrUnsetAction)
		}
		scope := domain.ProtocolScope("")
		if row.ProtocolScope != "" {
			scope, err = domain.ParseProtocolScope(row.ProtocolScope)
			if err != nil {
				return nil, err
			}
		}
		limitation := domain.UDPLimitationAbsent
		if row.UDPLimitation != "" {
			limitation, err = domain.ParseUDPLimitation(row.UDPLimitation)
			if err != nil {
				return nil, err
			}
		}
		policy, err := domain.NewApplicationPolicy(processID, action, scope, limitation)
		if err != nil {
			return nil, err
		}
		out = append(out, policy)
	}
	return out, nil
}
