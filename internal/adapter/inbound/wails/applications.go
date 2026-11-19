package wails

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ApplicationsService exposes application catalog operations.
type ApplicationsService struct {
	Catalog *service.Catalog
	Logger  *slog.Logger
	Emit    Emitter
}

// ApplicationDTO is one catalog application.
type ApplicationDTO struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	BundleID    string `json:"bundleId"`
}

// HelperDTO is a helper process candidate.
type HelperDTO struct {
	BundleID       string `json:"bundleId"`
	ExecutableName string `json:"executableName"`
	Path           string `json:"path"`
	DisplayName    string `json:"displayName"`
	IsUpdater      bool   `json:"isUpdater"`
}

// HelperGroupDTO is one bundle-id group from Go (ADR-0014).
type HelperGroupDTO struct {
	BundleID        string    `json:"bundleId"`
	DisplayName     string    `json:"displayName"`
	ExecutableNames []string  `json:"executableNames"`
	ProcessIDs      []string  `json:"processIds"`
	Confirmed       bool      `json:"confirmed"`
	Preselected     bool      `json:"preselected"`
	IsUpdater       bool      `json:"isUpdater"`
	SameDeveloper   bool      `json:"sameDeveloper"`
	Representative  HelperDTO `json:"representative"`
}

// AddApplicationResultDTO is returned after picking an app.
type AddApplicationResultDTO struct {
	Application  ApplicationDTO   `json:"application"`
	Helpers      []HelperDTO      `json:"helpers"`
	Groups       []HelperGroupDTO `json:"groups"`
	PendingCount int              `json:"pendingCount"`
}

// ProcessIdentityDTO is a saved process identity.
type ProcessIdentityDTO struct {
	ID             string `json:"id"`
	ApplicationID  string `json:"applicationId"`
	BundleID       string `json:"bundleId"`
	ExecutableName string `json:"executableName"`
	Path           string `json:"path"`
}

// InstalledApplicationDTO is one browsable installed application.
type InstalledApplicationDTO struct {
	DisplayName string `json:"displayName"`
	BundleID    string `json:"bundleId"`
	Path        string `json:"path"`
	InCatalog   bool   `json:"inCatalog"`
	IsSystem    bool   `json:"isSystem"`
}

func (s *ApplicationsService) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// List returns applications.
func (s *ApplicationsService) List(ctx context.Context) ([]ApplicationDTO, error) {
	if s == nil || s.Catalog == nil {
		return nil, fmt.Errorf("applications service: not configured")
	}
	apps, err := s.Catalog.ListApplications(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]ApplicationDTO, 0, len(apps))
	for _, app := range apps {
		out = append(out, ApplicationDTO{
			ID:          string(app.ID),
			DisplayName: app.DisplayName,
			BundleID:    app.BundleID,
		})
	}
	return out, nil
}

// ListProcessIdentities returns all process identities.
func (s *ApplicationsService) ListProcessIdentities(ctx context.Context) ([]ProcessIdentityDTO, error) {
	if s == nil || s.Catalog == nil {
		return nil, fmt.Errorf("applications service: not configured")
	}
	identities, err := s.Catalog.ListProcessIdentities(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]ProcessIdentityDTO, 0, len(identities))
	for _, identity := range identities {
		out = append(out, ProcessIdentityDTO{
			ID:             string(identity.ID),
			ApplicationID:  string(identity.ApplicationID),
			BundleID:       identity.BundleID,
			ExecutableName: identity.ExecutableName,
			Path:           identity.Path,
		})
	}
	return out, nil
}

// ListInstalled returns platform-installed applications without writing the catalog.
func (s *ApplicationsService) ListInstalled(ctx context.Context) ([]InstalledApplicationDTO, error) {
	if s == nil || s.Catalog == nil {
		return nil, fmt.Errorf("applications service: not configured")
	}
	rows, err := s.Catalog.ListInstalled(ctx)
	if err != nil {
		s.logger().Warn("list installed failed", "component", "wails.applications", "err", err)
		return nil, err
	}
	out := make([]InstalledApplicationDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, InstalledApplicationDTO{
			DisplayName: row.DisplayName,
			BundleID:    row.BundleID,
			Path:        row.Path,
			InCatalog:   row.InCatalog,
			IsSystem:    row.IsSystem,
		})
	}
	return out, nil
}

// Icon returns a PNG icon for an application path.
func (s *ApplicationsService) Icon(ctx context.Context, path string) ([]byte, error) {
	if s == nil || s.Catalog == nil {
		return nil, fmt.Errorf("applications service: not configured")
	}
	png, err := s.Catalog.ApplicationIconPNG(ctx, path)
	if err != nil {
		s.logger().Debug("application icon failed", "component", "wails.applications", "err", err)
		return nil, err
	}
	return png, nil
}

// AddFromPath inspects path and adds the main application to the catalog.
func (s *ApplicationsService) AddFromPath(ctx context.Context, path string) (AddApplicationResultDTO, error) {
	if s == nil || s.Catalog == nil {
		return AddApplicationResultDTO{}, fmt.Errorf("applications service: not configured")
	}
	if path == "" {
		return AddApplicationResultDTO{}, fmt.Errorf("applications add from path: empty path")
	}
	saved, inspected, err := s.Catalog.AddApplicationFromPath(ctx, path)
	if err != nil {
		s.logger().Warn("add application from path failed", "component", "wails.applications", "err", err)
		return AddApplicationResultDTO{}, err
	}
	result := mapAddResult(saved, inspected)
	if groups, pending, groupErr := s.Catalog.HelperGroups(ctx, saved.ID); groupErr == nil {
		result.Groups = mapHelperGroups(groups)
		result.PendingCount = pending
	}
	emit(s.Emit, EventCatalogChanged, map[string]any{"reason": "add_application", "applicationId": string(saved.ID)})
	return result, nil
}

// PickAndAdd opens a native file dialog, inspects the .app, and saves the main bundle.
func (s *ApplicationsService) PickAndAdd(ctx context.Context) (AddApplicationResultDTO, error) {
	if s == nil || s.Catalog == nil {
		return AddApplicationResultDTO{}, fmt.Errorf("applications service: not configured")
	}
	app := application.Get()
	if app == nil {
		return AddApplicationResultDTO{}, fmt.Errorf("applications pick: wails application unavailable")
	}
	path, err := app.Dialog.OpenFile().
		SetTitle("Select Application").
		CanChooseFiles(true).
		CanChooseDirectories(false).
		TreatsFilePackagesAsDirectories(false).
		AddFilter("Applications", "*.app").
		PromptForSingleSelection()
	if err != nil {
		return AddApplicationResultDTO{}, fmt.Errorf("applications pick: %w", err)
	}
	if path == "" {
		// User dismissed the dialog; not a failure for the UI.
		return AddApplicationResultDTO{}, nil
	}
	saved, inspected, err := s.Catalog.AddApplicationFromPath(ctx, path)
	if err != nil {
		s.logger().Warn("add application failed", "component", "wails.applications", "err", err)
		return AddApplicationResultDTO{}, err
	}
	result := mapAddResult(saved, inspected)
	if groups, pending, groupErr := s.Catalog.HelperGroups(ctx, saved.ID); groupErr == nil {
		result.Groups = mapHelperGroups(groups)
		result.PendingCount = pending
	}
	emit(s.Emit, EventCatalogChanged, map[string]any{"reason": "add_application", "applicationId": string(saved.ID)})
	return result, nil
}

func mapAddResult(saved domain.Application, inspected port.InspectedApplication) AddApplicationResultDTO {
	helpers := make([]HelperDTO, 0, len(inspected.Helpers))
	for _, helper := range inspected.Helpers {
		helpers = append(helpers, HelperDTO{
			BundleID:       helper.BundleID,
			ExecutableName: helper.ExecutableName,
			Path:           helper.Path,
			DisplayName:    helper.DisplayName,
			IsUpdater:      helper.IsUpdater,
		})
	}
	return AddApplicationResultDTO{
		Application: ApplicationDTO{
			ID:          string(saved.ID),
			DisplayName: saved.DisplayName,
			BundleID:    saved.BundleID,
		},
		Helpers: helpers,
	}
}

func mapHelperGroups(groups []service.HelperGroup) []HelperGroupDTO {
	out := make([]HelperGroupDTO, 0, len(groups))
	for _, group := range groups {
		ids := make([]string, 0, len(group.ProcessIDs))
		for _, id := range group.ProcessIDs {
			ids = append(ids, string(id))
		}
		out = append(out, HelperGroupDTO{
			BundleID:        group.BundleID,
			DisplayName:     group.DisplayName,
			ExecutableNames: group.ExecutableNames,
			ProcessIDs:      ids,
			Confirmed:       group.Confirmed,
			Preselected:     group.Preselected,
			IsUpdater:       group.IsUpdater,
			SameDeveloper:   group.SameDeveloper,
			Representative: HelperDTO{
				BundleID:       group.Representative.BundleID,
				ExecutableName: group.Representative.ExecutableName,
				Path:           group.Representative.Path,
				DisplayName:    group.Representative.DisplayName,
				IsUpdater:      group.Representative.IsUpdater,
			},
		})
	}
	return out
}

// ListHelperGroups returns Go-built helper groups for an application (ADR-0014).
func (s *ApplicationsService) ListHelperGroups(ctx context.Context, applicationID string) ([]HelperGroupDTO, error) {
	if s == nil || s.Catalog == nil {
		return nil, fmt.Errorf("applications service: not configured")
	}
	appID, err := domain.ParseApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("list helper groups: %w", err)
	}
	groups, _, err := s.Catalog.HelperGroups(ctx, appID)
	if err != nil {
		return nil, err
	}
	return mapHelperGroups(groups), nil
}

// InspectHelpers re-inspects helpers and returns groups with pendingCount (ADR-0014).
func (s *ApplicationsService) InspectHelpers(ctx context.Context, applicationID string) (AddApplicationResultDTO, error) {
	if s == nil || s.Catalog == nil {
		return AddApplicationResultDTO{}, fmt.Errorf("applications service: not configured")
	}
	appID, err := domain.ParseApplicationID(applicationID)
	if err != nil {
		return AddApplicationResultDTO{}, fmt.Errorf("inspect helpers: %w", err)
	}
	apps, err := s.Catalog.ListApplications(ctx)
	if err != nil {
		return AddApplicationResultDTO{}, fmt.Errorf("inspect helpers: %w", err)
	}
	var app domain.Application
	found := false
	for _, row := range apps {
		if row.ID == appID {
			app = row
			found = true
			break
		}
	}
	if !found {
		return AddApplicationResultDTO{}, fmt.Errorf("inspect helpers: application not found")
	}
	groups, pending, err := s.Catalog.HelperGroups(ctx, appID)
	if err != nil {
		s.logger().Warn("inspect helpers failed", "component", "wails.applications", "err", err)
		return AddApplicationResultDTO{}, err
	}
	helpers := make([]HelperDTO, 0, len(groups))
	for _, group := range groups {
		helpers = append(helpers, HelperDTO{
			BundleID:       group.Representative.BundleID,
			ExecutableName: group.Representative.ExecutableName,
			Path:           group.Representative.Path,
			DisplayName:    group.DisplayName,
			IsUpdater:      group.IsUpdater,
		})
	}
	return AddApplicationResultDTO{
		Application: ApplicationDTO{
			ID:          string(app.ID),
			DisplayName: app.DisplayName,
			BundleID:    app.BundleID,
		},
		Helpers:      helpers,
		Groups:       mapHelperGroups(groups),
		PendingCount: pending,
	}, nil
}

// ConfirmHelpers syncs selected helpers for an application (ADR-0014).
func (s *ApplicationsService) ConfirmHelpers(ctx context.Context, applicationID string, helpers []HelperDTO) error {
	if s == nil || s.Catalog == nil {
		return fmt.Errorf("applications service: not configured")
	}
	appID, err := domain.ParseApplicationID(applicationID)
	if err != nil {
		return fmt.Errorf("confirm helpers: %w", err)
	}
	inspected := make([]port.InspectedProcess, 0, len(helpers))
	for _, helper := range helpers {
		inspected = append(inspected, port.InspectedProcess{
			BundleID:       helper.BundleID,
			ExecutableName: helper.ExecutableName,
			Path:           helper.Path,
			DisplayName:    helper.DisplayName,
			IsHelper:       true,
			IsUpdater:      helper.IsUpdater,
		})
	}
	if err := s.Catalog.ConfirmHelpers(ctx, appID, inspected); err != nil {
		return err
	}
	emit(s.Emit, EventCatalogChanged, map[string]any{"reason": "confirm_helpers", "applicationId": applicationID})
	return nil
}

// Delete removes an application when unreferenced after a native confirm (Cancel default).
func (s *ApplicationsService) Delete(ctx context.Context, id string) error {
	if s == nil || s.Catalog == nil {
		return fmt.Errorf("applications service: not configured")
	}
	appID, err := domain.ParseApplicationID(id)
	if err != nil {
		return fmt.Errorf("applications delete: %w", err)
	}
	name := id
	if apps, listErr := s.Catalog.ListApplications(ctx); listErr == nil {
		for _, app := range apps {
			if app.ID == appID {
				if app.DisplayName != "" {
					name = app.DisplayName
				}
				break
			}
		}
	}
	if !confirmDelete(
		"Delete application?",
		fmt.Sprintf("Delete this application and its confirmed helpers from the catalog?\n\n%s", name),
		"Delete",
	) {
		return nil
	}
	if err := s.Catalog.DeleteApplication(ctx, appID); err != nil {
		s.logger().Warn("application delete failed", "component", "wails.applications", "err", err)
		return err
	}
	emit(s.Emit, EventCatalogChanged, map[string]any{"reason": "delete_application", "applicationId": id})
	return nil
}
