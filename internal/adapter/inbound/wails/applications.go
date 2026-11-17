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
}

// AddApplicationResultDTO is returned after picking an app.
type AddApplicationResultDTO struct {
	Application ApplicationDTO `json:"application"`
	Helpers     []HelperDTO    `json:"helpers"`
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
	return mapAddResult(saved, inspected), nil
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
	return mapAddResult(saved, inspected), nil
}

func mapAddResult(saved domain.Application, inspected port.InspectedApplication) AddApplicationResultDTO {
	helpers := make([]HelperDTO, 0, len(inspected.Helpers))
	for _, helper := range inspected.Helpers {
		helpers = append(helpers, HelperDTO{
			BundleID:       helper.BundleID,
			ExecutableName: helper.ExecutableName,
			Path:           helper.Path,
			DisplayName:    helper.DisplayName,
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

// ConfirmHelpers saves selected helpers for an application.
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
		})
	}
	return s.Catalog.ConfirmHelpers(ctx, appID, inspected)
}

// Delete removes an application when unreferenced.
func (s *ApplicationsService) Delete(ctx context.Context, id string) error {
	if s == nil || s.Catalog == nil {
		return fmt.Errorf("applications service: not configured")
	}
	appID, err := domain.ParseApplicationID(id)
	if err != nil {
		return fmt.Errorf("applications delete: %w", err)
	}
	if err := s.Catalog.DeleteApplication(ctx, appID); err != nil {
		s.logger().Warn("application delete failed", "component", "wails.applications", "err", err)
		return err
	}
	return nil
}
