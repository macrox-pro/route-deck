package wails

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ProfilesService exposes profile CRUD and desired selection.
type ProfilesService struct {
	Catalog    *service.Catalog
	Settings   port.SettingsStore
	Protection *service.Protection
	Emit       Emitter
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
		s.logger().Warn("profile save failed", "component", "wails.profiles", "err", err)
		return ProfileDTO{}, err
	}
	if s.Protection != nil {
		appliedID, confirmed := s.Protection.Applied()
		if confirmed && appliedID == profile.ID {
			s.Protection.MarkDesiredDiffers()
		}
	}
	emit(s.Emit, EventCatalogChanged, map[string]any{"reason": "save_profile", "profileId": string(profile.ID)})
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
	emit(s.Emit, EventCatalogChanged, map[string]any{"reason": "copy_profile", "profileId": string(copied.ID)})
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

// Delete confirms with a native Question dialog, then removes the profile.
// Desired is cleared in the same path. Applied profiles are Stopped first (ADR-0019 UX).
// Cancel returns nil without deleting. MUST NOT use window.confirm (unreliable in Wails).
func (s *ProfilesService) Delete(ctx context.Context, id string) error {
	if s == nil || s.Catalog == nil || s.Settings == nil {
		return fmt.Errorf("profiles service: not configured")
	}
	profileID, err := domain.ParseProfileID(id)
	if err != nil {
		return fmt.Errorf("profiles delete: %w", err)
	}
	name := id
	profiles, listErr := s.Catalog.ListProfiles(ctx)
	if listErr != nil {
		return fmt.Errorf("profiles delete: %w", listErr)
	}
	found := false
	for _, profile := range profiles {
		if profile.ID == profileID {
			found = true
			if profile.Name != "" {
				name = profile.Name
			}
			break
		}
	}
	if !found {
		return fmt.Errorf("profiles delete: profile not found")
	}

	settings, err := s.Settings.Get(ctx)
	if err != nil {
		return fmt.Errorf("profiles delete: %w", err)
	}
	isDesired := settings.DesiredSet && settings.DesiredProfileID == profileID
	var appliedID domain.ProfileID
	var isApplied bool
	if s.Protection != nil {
		appliedID, isApplied = s.Protection.Applied()
		isApplied = isApplied && appliedID == profileID
	}

	if !confirmProfileDelete(name, isDesired, isApplied) {
		s.logger().Info("profile delete cancelled", "component", "wails.profiles", "profileId", id)
		return nil
	}

	if isApplied {
		if s.Protection == nil {
			return fmt.Errorf("stop protection before deleting the applied profile")
		}
		if stopErr := s.Protection.Stop(ctx, "delete-profile"); stopErr != nil {
			s.logger().Warn("profile delete stop failed", "component", "wails.profiles", "err", stopErr)
			return fmt.Errorf("stop protection before deleting the applied profile: %w", stopErr)
		}
		emit(s.Emit, EventProtectionStateChanged, map[string]any{"reason": "delete_profile_stop", "profileId": id})
		emit(s.Emit, EventProxyBridgeStateChanged, map[string]any{"reason": "delete_profile_stop", "profileId": id})
	}

	if isDesired {
		settings.DesiredProfileID = ""
		settings.DesiredSet = false
		if err := s.Settings.Save(ctx, settings); err != nil {
			return fmt.Errorf("profiles delete clear desired: %w", err)
		}
	}

	if err := s.Catalog.DeleteProfile(ctx, profileID, ""); err != nil {
		s.logger().Warn("profile delete failed", "component", "wails.profiles", "err", err)
		if errors.Is(err, service.ErrProfileApplied) {
			return fmt.Errorf("stop protection before deleting the applied profile")
		}
		return err
	}
	s.logger().Info("profile deleted", "component", "wails.profiles", "profileId", id)
	emit(s.Emit, EventCatalogChanged, map[string]any{"reason": "delete_profile", "profileId": id})
	return nil
}

// confirmProfileDelete shows a native Question. Cancel is the default (desktop-ux).
func confirmProfileDelete(name string, desired, applied bool) bool {
	app := application.Get()
	if app == nil {
		return false
	}
	title := "Delete profile?"
	message := fmt.Sprintf("Delete this profile?\n\n%s", name)
	deleteLabel := "Delete"
	switch {
	case applied && desired:
		message = fmt.Sprintf("This profile is Desired and Applied. Stop protection, clear Desired, and delete?\n\n%s", name)
		deleteLabel = "Stop and Delete"
	case applied:
		message = fmt.Sprintf("This profile is Applied. Stop protection and delete?\n\n%s", name)
		deleteLabel = "Stop and Delete"
	case desired:
		message = fmt.Sprintf("This is your Desired profile. Clear Desired and delete?\n\n%s", name)
	}

	result := make(chan bool, 1)
	dialog := app.Dialog.Question().SetTitle(title).SetMessage(message)
	cancel := dialog.AddButton("Cancel").SetAsCancel().SetAsDefault()
	cancel.OnClick(func() { result <- false })
	del := dialog.AddButton(deleteLabel)
	del.OnClick(func() { result <- true })
	dialog.Show()
	return <-result
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
