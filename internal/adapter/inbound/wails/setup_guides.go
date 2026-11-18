package wails

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// SetupGuidesService exposes ADR-0013 setup guides over Wails.
type SetupGuidesService struct {
	Guides *service.SetupGuides
	Logger *slog.Logger
}

// GuideStepDTO is one instruction step.
type GuideStepDTO struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	Warning bool   `json:"warning"`
}

// GuideEndpointDTO is a documented local listener.
type GuideEndpointDTO struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Kind string `json:"kind"`
	Note string `json:"note"`
}

// GuideDTO is one setup guide with detection.
type GuideDTO struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	Kind          string             `json:"kind"`
	Detection     string             `json:"detection"`
	InstalledPath string             `json:"installedPath"`
	PortOpen      bool               `json:"portOpen"`
	PortChecked   int                `json:"portChecked"`
	AlreadyAdded  bool               `json:"alreadyAdded"`
	AddMode       string             `json:"addMode"`
	DocsURL       string             `json:"docsUrl"`
	ConnectHint   bool               `json:"connectHint"`
	Steps         []GuideStepDTO     `json:"steps"`
	Endpoints     []GuideEndpointDTO `json:"endpoints"`
}

// SuggestionDTO is a detected open local proxy banner row.
type SuggestionDTO struct {
	GuideID string `json:"guideId"`
	Title   string `json:"title"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Kind    string `json:"kind"`
}

// ChecklistDTO is the first-run checklist.
// ProxyStep is a closed token (missing|available|udp_limited|checking); English UI copy maps by message ID on the frontend for future i18n.
type ChecklistDTO struct {
	BridgeFound      bool   `json:"bridgeFound"`
	ProxyAvailable   bool   `json:"proxyAvailable"`
	ProxyPartial     bool   `json:"proxyPartial"`
	ProxyStep        string `json:"proxyStep"`
	ApplicationAdded bool   `json:"applicationAdded"`
	ProfileApplied   bool   `json:"profileApplied"`
	Complete         bool   `json:"complete"`
	Dismissed        bool   `json:"dismissed"`
}

// ProxyDraftDTO is returned after one-click add (no secrets).
type ProxyDraftDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (s *SetupGuidesService) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

func guideDTO(v service.GuideView) GuideDTO {
	steps := make([]GuideStepDTO, 0, len(v.Steps))
	for _, step := range v.Steps {
		steps = append(steps, GuideStepDTO{Title: step.Title, Body: step.Body, Warning: step.Warning})
	}
	endpoints := make([]GuideEndpointDTO, 0, len(v.Endpoints))
	for _, ep := range v.Endpoints {
		endpoints = append(endpoints, GuideEndpointDTO{
			Host: ep.Host,
			Port: ep.Port,
			Kind: string(ep.Kind),
			Note: ep.Note,
		})
	}
	return GuideDTO{
		ID:            v.ID,
		Title:         v.Title,
		Kind:          v.Kind,
		Detection:     v.Detection,
		InstalledPath: v.InstalledPath,
		PortOpen:      v.PortOpen,
		PortChecked:   v.PortChecked,
		AlreadyAdded:  v.AlreadyAdded,
		AddMode:       v.AddMode,
		DocsURL:       v.DocsURL,
		ConnectHint:   v.ConnectHint,
		Steps:         steps,
		Endpoints:     endpoints,
	}
}

// List returns detected guides. Pass refresh=true to reload installed apps.
func (s *SetupGuidesService) List(ctx context.Context, refresh bool) ([]GuideDTO, error) {
	if s == nil || s.Guides == nil {
		return nil, fmt.Errorf("setup guides service: not configured")
	}
	views, err := s.Guides.List(ctx, refresh)
	if err != nil {
		return nil, err
	}
	out := make([]GuideDTO, 0, len(views))
	for _, v := range views {
		out = append(out, guideDTO(v))
	}
	return out, nil
}

// OpenDocs opens the guide's official docs URL in the system browser.
func (s *SetupGuidesService) OpenDocs(ctx context.Context, guideID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s == nil || s.Guides == nil {
		return fmt.Errorf("setup guides service: not configured")
	}
	url, err := s.Guides.DocsURL(guideID)
	if err != nil {
		return mapGuideErr(err)
	}
	app := application.Get()
	if app == nil {
		return fmt.Errorf("open docs: wails application unavailable")
	}
	if err := app.Browser.OpenURL(url); err != nil {
		s.logger().Warn("open docs failed", "component", "wails.setup_guides", "err", err)
		return fmt.Errorf("open docs: %w", err)
	}
	return nil
}

// AddProxy one-click saves the guide default endpoint and rechecks health.
func (s *SetupGuidesService) AddProxy(ctx context.Context, guideID string) (ProxyDraftDTO, error) {
	if s == nil || s.Guides == nil {
		return ProxyDraftDTO{}, fmt.Errorf("setup guides service: not configured")
	}
	ep, err := s.Guides.AddProxy(ctx, guideID)
	if err != nil {
		return ProxyDraftDTO{}, mapGuideErr(err)
	}
	return ProxyDraftDTO{
		ID:   string(ep.ID),
		Name: ep.Name,
		Kind: string(ep.Kind),
		Host: ep.Host,
		Port: ep.Port,
	}, nil
}

// Suggestions returns banner candidates for this session.
func (s *SetupGuidesService) Suggestions(ctx context.Context) ([]SuggestionDTO, error) {
	if s == nil || s.Guides == nil {
		return nil, fmt.Errorf("setup guides service: not configured")
	}
	views, err := s.Guides.Suggestions(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]SuggestionDTO, 0, len(views))
	for _, v := range views {
		out = append(out, SuggestionDTO{
			GuideID: v.GuideID,
			Title:   v.Title,
			Host:    v.Host,
			Port:    v.Port,
			Kind:    v.Kind,
		})
	}
	return out, nil
}

// DismissSuggestion hides a suggestion until RouteDeck restarts.
func (s *SetupGuidesService) DismissSuggestion(_ context.Context, guideID string, port int) error {
	if s == nil || s.Guides == nil {
		return fmt.Errorf("setup guides service: not configured")
	}
	return mapGuideErr(s.Guides.DismissSuggestion(guideID, port))
}

// Checklist returns the first-run checklist.
func (s *SetupGuidesService) Checklist(ctx context.Context) (ChecklistDTO, error) {
	if s == nil || s.Guides == nil {
		return ChecklistDTO{}, fmt.Errorf("setup guides service: not configured")
	}
	view, err := s.Guides.Checklist(ctx)
	if err != nil {
		return ChecklistDTO{}, err
	}
	return ChecklistDTO{
		BridgeFound:      view.BridgeFound,
		ProxyAvailable:   view.ProxyAvailable,
		ProxyPartial:     view.ProxyPartial,
		ProxyStep:        string(view.ProxyStep),
		ApplicationAdded: view.ApplicationAdded,
		ProfileApplied:   view.ProfileApplied,
		Complete:         view.Complete,
		Dismissed:        view.Dismissed,
	}, nil
}

// DismissChecklist hides the first-run checklist until restart or step regression.
func (s *SetupGuidesService) DismissChecklist(ctx context.Context) error {
	if s == nil || s.Guides == nil {
		return fmt.Errorf("setup guides service: not configured")
	}
	return s.Guides.DismissChecklist(ctx)
}

func mapGuideErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, service.ErrUnknownGuideID):
		return err
	case errors.Is(err, service.ErrGuideHasNoProxy):
		return err
	case errors.Is(err, service.ErrGuideNeedsPort):
		return err
	case errors.Is(err, service.ErrEndpointAlreadyAdded):
		return err
	default:
		return err
	}
}
