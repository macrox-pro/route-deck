package wails

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/tray"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ErrAppearanceCancelled means the user cancelled Discard before applying appearance.
var ErrAppearanceCancelled = errors.New("appearance change cancelled")

// WindowChrome recreates the main window so native frame options can change.
// Wails v3.0.0 exposes Mac.Appearance / Windows.Theme only at create time.
type WindowChrome struct {
	App     *application.App
	Editors *EditorDialogs
	Tray    *tray.Controller
	Dialog  *NativeUserDialog

	mu         sync.Mutex
	main       *application.WebviewWindow
	appearance string
}

// SetMain stores the current main window pointer.
func (c *WindowChrome) SetMain(win *application.WebviewWindow) {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.main = win
	c.mu.Unlock()
}

// Appearance returns the last applied appearance (default system).
func (c *WindowChrome) Appearance() string {
	if c == nil {
		return service.AppearanceSystem
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.appearance == "" {
		return service.AppearanceSystem
	}
	return c.appearance
}

// SetAppearanceMemory updates the in-memory appearance used for new editor windows.
func (c *WindowChrome) SetAppearanceMemory(appearance string) {
	if c == nil {
		return
	}
	normalized := service.CoerceAppearance(appearance)
	c.mu.Lock()
	c.appearance = normalized
	c.mu.Unlock()
	if c.Editors != nil {
		c.Editors.SetAppearance(normalized)
	}
}

// MainWindowOptions builds verified create-time options for the main window.
func MainWindowOptions(appearance, url string) application.WebviewWindowOptions {
	opts := application.WebviewWindowOptions{
		Name:      "main",
		Title:     "MACROX RouteDeck",
		Width:     1000,
		Height:    618,
		MinWidth:  800,
		MinHeight: 500,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBarDefault,
		},
		URL: mainURL(url),
	}
	applyNativeAppearance(&opts, appearance)
	return opts
}

// DialogWindowOptions merges dialog geometry with native appearance.
func DialogWindowOptions(appearance string, base application.WebviewWindowOptions) application.WebviewWindowOptions {
	applyNativeAppearance(&base, appearance)
	return base
}

func applyNativeAppearance(opts *application.WebviewWindowOptions, appearance string) {
	normalized := service.CoerceAppearance(appearance)
	switch normalized {
	case service.AppearanceLight:
		opts.Mac.Appearance = application.NSAppearanceNameAqua
		opts.Windows.Theme = application.Light
	case service.AppearanceDark:
		opts.Mac.Appearance = application.NSAppearanceNameDarkAqua
		opts.Windows.Theme = application.Dark
	case service.AppearanceSystem:
		opts.Mac.Appearance = application.DefaultAppearance
		opts.Windows.Theme = application.SystemDefault
	}
	// Do not set BackgroundColour (desktop-ux: no white override).
}

func mainURL(hashOrPath string) string {
	h := strings.TrimSpace(hashOrPath)
	if h == "" || h == "#" || h == "/" {
		return "/"
	}
	if strings.HasPrefix(h, "/#") {
		return h
	}
	if strings.HasPrefix(h, "#") {
		return "/" + h
	}
	if strings.HasPrefix(h, "/") {
		return "/#" + h
	}
	return "/#/" + h
}

// PrepareForAppearanceChange asks to discard dirty editors; Cancel returns ErrAppearanceCancelled.
// On success every editor window is force-closed.
func (c *WindowChrome) PrepareForAppearanceChange() error {
	if c == nil || c.Editors == nil {
		return nil
	}
	if c.Editors.HasDirty() {
		ok := c.confirmDiscard()
		if !ok {
			return ErrAppearanceCancelled
		}
	}
	c.Editors.CloseAllForced()
	return nil
}

func (c *WindowChrome) confirmDiscard() bool {
	if c == nil || c.App == nil {
		return false
	}
	done := make(chan bool, 1)
	dialog := c.App.Dialog.Question().
		SetTitle("Discard changes?").
		SetMessage("Discard changes?")
	cancel := dialog.AddButton("Cancel").SetAsCancel()
	cancel.OnClick(func() { done <- false })
	discard := dialog.AddButton("Discard").SetAsDefault()
	discard.OnClick(func() { done <- true })
	dialog.Show()
	return <-done
}

// RecreateMain closes the main window and opens a new one with appearance + URL.
func (c *WindowChrome) RecreateMain(appearance, mainHash string) error {
	if c == nil || c.App == nil {
		return fmt.Errorf("window chrome: not configured")
	}
	c.SetAppearanceMemory(appearance)

	c.mu.Lock()
	old := c.main
	c.mu.Unlock()
	if old != nil {
		old.Close()
	}

	win := c.App.Window.NewWithOptions(MainWindowOptions(appearance, mainHash))
	win.Show().Focus()
	c.rewire(win)
	return nil
}

func (c *WindowChrome) rewire(win *application.WebviewWindow) {
	c.SetMain(win)
	if c.Editors != nil {
		c.Editors.Parent = win
	}
	if c.Tray != nil {
		c.Tray.Window = win
	}
	if c.Dialog != nil {
		c.Dialog.Window = win
	}
}
