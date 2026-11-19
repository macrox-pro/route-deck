package wails

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Utility window kinds (ADR-0020). Non-modal; not editor dialogs.
const (
	UtilityKindGuides = "guides"
	UtilityKindAbout  = "about"
)

var (
	allowedGuideSteps = map[string]struct{}{
		"":            {},
		"bridge":      {},
		"proxy":       {},
		"application": {},
		"profile":     {},
	}
	allowedMainRoutes = map[string]struct{}{
		"/":             {},
		"/applications": {},
		"/proxies":      {},
		"/profiles":     {},
		"/activity":     {},
		"/settings":     {},
	}
)

// UtilityWindows manages reused non-modal native-frame utility windows.
type UtilityWindows struct {
	App     *application.App
	Parent  *application.WebviewWindow
	Editors *EditorDialogs
	Logger  *slog.Logger
	Emit    Emitter

	mu         sync.Mutex
	windows    map[string]*application.WebviewWindow
	appearance string
}

// SetAppearance stores the preference and closes open utility windows (reopen on next Open).
func (u *UtilityWindows) SetAppearance(appearance string) {
	if u == nil {
		return
	}
	u.mu.Lock()
	u.appearance = appearance
	u.mu.Unlock()
	u.CloseAllForced()
}

// CloseAllForced force-closes every utility window.
func (u *UtilityWindows) CloseAllForced() {
	if u == nil {
		return
	}
	u.mu.Lock()
	u.ensureMaps()
	kinds := make([]string, 0, len(u.windows))
	for kind := range u.windows {
		kinds = append(kinds, kind)
	}
	u.mu.Unlock()
	for _, kind := range kinds {
		_ = u.Close(kind)
	}
}

func (u *UtilityWindows) ensureMaps() {
	if u.windows == nil {
		u.windows = make(map[string]*application.WebviewWindow)
	}
}

// Open opens or focuses a utility window for kind at hashPath (without leading /#).
func (u *UtilityWindows) Open(kind, hashPath string) error {
	if u == nil || u.App == nil {
		return fmt.Errorf("utility windows: not configured")
	}
	title, width, height, resizable, minW, minH, err := utilitySpec(kind)
	if err != nil {
		return err
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	u.ensureMaps()

	if w := u.windows[kind]; w != nil {
		w.SetTitle(title)
		w.SetURL("/#" + hashPath)
		w.Show().Focus()
		return nil
	}

	appearance := u.appearance
	opts := application.WebviewWindowOptions{
		Name:          "utility-" + kind,
		Title:         title,
		Width:         width,
		Height:        height,
		MinWidth:      minW,
		MinHeight:     minH,
		DisableResize: !resizable,
		Hidden:        true,
		URL:           "/#" + hashPath,
	}
	opts = DialogWindowOptions(appearance, opts)
	win := u.App.Window.NewWithOptions(opts)
	u.windows[kind] = win

	win.OnWindowEvent(events.Common.WindowClosing, func(_ *application.WindowEvent) {
		u.mu.Lock()
		delete(u.windows, kind)
		u.mu.Unlock()
	})

	win.Show().Focus()
	return nil
}

// Close closes one utility window by kind.
func (u *UtilityWindows) Close(kind string) error {
	if u == nil {
		return nil
	}
	if kind != UtilityKindGuides && kind != UtilityKindAbout {
		return fmt.Errorf("utility windows: unknown kind %q", kind)
	}
	u.mu.Lock()
	w := u.windows[kind]
	u.mu.Unlock()
	if w != nil {
		w.Close()
	}
	return nil
}

func utilitySpec(kind string) (title string, width, height int, resizable bool, minW, minH int, err error) {
	switch kind {
	case UtilityKindGuides:
		return "Setup Guides", 880, 640, true, 720, 520, nil
	case UtilityKindAbout:
		return "About RouteDeck", 420, 460, false, 420, 460, nil
	default:
		return "", 0, 0, false, 0, 0, fmt.Errorf("utility windows: unknown kind %q", kind)
	}
}

// UtilityWindowsService exposes utility windows to the frontend (ADR-0020).
type UtilityWindowsService struct {
	Windows *UtilityWindows
	Logger  *slog.Logger
}

// OpenGuides opens the Setup Guides utility window. step is "" or a checklist step id.
func (s *UtilityWindowsService) OpenGuides(step string) error {
	if s == nil || s.Windows == nil {
		return fmt.Errorf("utility windows: not configured")
	}
	if _, ok := allowedGuideSteps[step]; !ok {
		return fmt.Errorf("utility windows: invalid guides step %q", step)
	}
	hash := "/dialog/guides"
	if step != "" {
		hash = "/dialog/guides/" + step
	}
	return s.Windows.Open(UtilityKindGuides, hash)
}

// OpenAbout opens the About RouteDeck utility window.
func (s *UtilityWindowsService) OpenAbout() error {
	if s == nil || s.Windows == nil {
		return fmt.Errorf("utility windows: not configured")
	}
	return s.Windows.Open(UtilityKindAbout, "/dialog/about")
}

// Close closes a utility window by kind.
func (s *UtilityWindowsService) Close(kind string) error {
	if s == nil || s.Windows == nil {
		return fmt.Errorf("utility windows: not configured")
	}
	return s.Windows.Close(kind)
}

// ShowMain shows the main window and emits main.command navigate, or focuses an open editor.
func (s *UtilityWindowsService) ShowMain(route string) error {
	if s == nil || s.Windows == nil {
		return fmt.Errorf("utility windows: not configured")
	}
	if _, ok := allowedMainRoutes[route]; !ok {
		return fmt.Errorf("utility windows: invalid main route %q", route)
	}
	u := s.Windows
	if u.Editors != nil && u.Editors.FocusOpen() {
		return nil
	}
	emit(u.Emit, EventMainCommand, map[string]any{
		"command": "navigate",
		"route":   route,
	})
	if u.Parent != nil {
		u.Parent.Show().Focus()
	}
	return nil
}
