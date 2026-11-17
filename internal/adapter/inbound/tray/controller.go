package tray

import (
	"context"
	"embed"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/macrox-pro/route-deck/internal/domain"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed icon_template.png
var iconFS embed.FS

// SnapshotSource loads protection and bridge words for the Status line.
type SnapshotSource interface {
	ProtectionAndBridge(ctx context.Context) (domain.ProtectionState, domain.BridgeState, error)
}

// BridgeReader returns the last known bridge state for Stop enablement.
type BridgeReader interface {
	BridgeState() domain.BridgeState
}

// Stopper stops the managed ProxyBridge CLI.
type Stopper interface {
	Stop(ctx context.Context, operationID string) error
}

// Controller owns the system tray menu and shared Quit dialog (ADR-0010).
type Controller struct {
	App        *application.App
	Window     *application.WebviewWindow
	Snapshot   SnapshotSource
	Bridge     BridgeReader
	Protection Stopper
	Logger     *slog.Logger

	tray       *application.SystemTray
	statusItem *application.MenuItem
	openItem   *application.MenuItem
	stopItem   *application.MenuItem

	mu            sync.Mutex
	protection    domain.ProtectionState
	bridge        domain.BridgeState
	allowQuit     atomic.Bool
	quitDialoging atomic.Bool
}

func (c *Controller) logger() *slog.Logger {
	if c != nil && c.Logger != nil {
		return c.Logger
	}
	return slog.Default()
}

// Attach wires the tray, close-to-hide, and ShouldQuit after the Wails app exists.
func (c *Controller) Attach() error {
	if c == nil || c.App == nil || c.Window == nil {
		return nil
	}

	icon, err := iconFS.ReadFile("icon_template.png")
	if err != nil {
		return err
	}

	c.protection = domain.ProtectionStateUnknown
	c.bridge = domain.BridgeStateUnknown
	if c.Bridge != nil {
		c.bridge = c.Bridge.BridgeState()
	}

	c.Window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		e.Cancel()
		c.Window.Hide()
		c.syncOpenItemVisibility()
	})
	c.Window.OnWindowEvent(events.Common.WindowShow, func(_ *application.WindowEvent) {
		c.syncOpenItemVisibility()
	})
	c.Window.OnWindowEvent(events.Common.WindowHide, func(_ *application.WindowEvent) {
		c.syncOpenItemVisibility()
	})

	menu := application.NewMenu()
	c.statusItem = menu.Add(StatusLabel(c.protection, c.bridge)).SetEnabled(false)
	c.openItem = menu.Add("Open RouteDeck").OnClick(func(_ *application.Context) {
		c.ShowWindow()
	})
	c.stopItem = menu.Add("Stop protection").OnClick(func(_ *application.Context) {
		c.onStop()
	})
	c.stopItem.SetEnabled(CLIOwned(c.bridge))
	menu.AddSeparator()
	menu.Add("Quit RouteDeck").OnClick(func(_ *application.Context) {
		c.RequestQuit()
	})

	c.tray = c.App.SystemTray.New()
	c.tray.SetTemplateIcon(icon)
	c.tray.SetTooltip(StatusLabel(c.protection, c.bridge))
	c.tray.SetMenu(menu)
	// Left-click opens the menu (macOS menu-bar norm). Window restore is
	// "Open RouteDeck" (only when hidden) or Dock click — not the tray glyph click.
	c.tray.OnClick(func() {
		c.openTrayMenu()
	})
	c.tray.OnRightClick(func() {
		c.openTrayMenu()
	})

	c.syncOpenItemVisibility()
	c.Refresh()
	return nil
}

func (c *Controller) openTrayMenu() {
	c.syncOpenItemVisibility()
	if c.tray != nil {
		c.tray.ShowMenu()
	}
}

// syncOpenItemVisibility hides Open RouteDeck while the main window is visible.
func (c *Controller) syncOpenItemVisibility() {
	if c == nil || c.openItem == nil || c.Window == nil {
		return
	}
	c.openItem.SetHidden(c.Window.IsVisible())
}

// ShouldQuit is the Wails Options.ShouldQuit hook. It always refuses until the
// shared Quit dialog confirms (then allowQuit is set and Quit is re-entered).
func (c *Controller) ShouldQuit() bool {
	if c == nil {
		return true
	}
	if c.allowQuit.Load() {
		return true
	}
	c.RequestQuit()
	return false
}

// ShowWindow shows and focuses the main window.
func (c *Controller) ShowWindow() {
	if c == nil || c.Window == nil {
		return
	}
	c.Window.Show().Focus()
	c.syncOpenItemVisibility()
}

// Refresh reloads Status / tooltip / Stop enablement from snapshot and bridge.
func (c *Controller) Refresh() {
	if c == nil {
		return
	}
	protection := domain.ProtectionStateUnknown
	bridge := domain.BridgeStateUnknown
	if c.Bridge != nil {
		bridge = c.Bridge.BridgeState()
	}
	if c.Snapshot != nil {
		p, b, err := c.Snapshot.ProtectionAndBridge(context.Background())
		if err != nil {
			c.logger().Debug("tray snapshot refresh failed", "component", "tray", "err", err)
		} else {
			protection = p
			if b != domain.BridgeStateUnknown {
				bridge = b
			}
		}
	}
	c.applyState(protection, bridge)
}

// OnProtectionEvent refreshes tray state after product events.
func (c *Controller) OnProtectionEvent() {
	c.Refresh()
}

func (c *Controller) applyState(protection domain.ProtectionState, bridge domain.BridgeState) {
	c.mu.Lock()
	c.protection = protection
	c.bridge = bridge
	label := StatusLabel(protection, bridge)
	owned := CLIOwned(bridge)
	c.mu.Unlock()

	if c.statusItem != nil {
		c.statusItem.SetLabel(label)
	}
	if c.stopItem != nil {
		c.stopItem.SetEnabled(owned)
	}
	if c.tray != nil {
		c.tray.SetTooltip(label)
	}
}

func (c *Controller) onStop() {
	c.mu.Lock()
	bridge := c.bridge
	c.mu.Unlock()
	if c.Bridge != nil {
		bridge = c.Bridge.BridgeState()
	}
	if !CLIOwned(bridge) {
		return
	}
	if c.Protection == nil {
		return
	}
	if err := c.Protection.Stop(context.Background(), "tray-stop"); err != nil {
		c.logger().Warn("tray stop failed", "component", "tray", "err", err)
	}
	c.Refresh()
}

// RequestQuit shows the shared Quit Question dialog.
func (c *Controller) RequestQuit() {
	if c == nil || c.App == nil {
		return
	}
	if c.allowQuit.Load() {
		c.App.Quit()
		return
	}
	if !c.quitDialoging.CompareAndSwap(false, true) {
		return
	}

	if c.Window != nil && !c.Window.IsVisible() {
		c.ShowWindow()
	}

	c.mu.Lock()
	bridge := c.bridge
	c.mu.Unlock()
	if c.Bridge != nil {
		bridge = c.Bridge.BridgeState()
	}
	owned := CLIOwned(bridge)

	dialog := c.App.Dialog.Question().
		SetTitle("Quit RouteDeck").
		SetMessage(QuitMessage(owned))
	if c.Window != nil {
		dialog.AttachToWindow(c.Window)
	}
	cancel := dialog.AddButton("Cancel").SetAsCancel()
	cancel.OnClick(func() {
		c.quitDialoging.Store(false)
	})
	quit := dialog.AddButton("Quit").SetAsDefault()
	quit.OnClick(func() {
		if owned && c.Protection != nil {
			if err := c.Protection.Stop(context.Background(), "tray-quit"); err != nil {
				c.logger().Error("tray quit stop failed", "component", "tray", "err", err)
			}
		}
		c.allowQuit.Store(true)
		c.quitDialoging.Store(false)
		c.App.Quit()
	})
	dialog.Show()
}
