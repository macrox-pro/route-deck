package wails

import (
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/tray"
	"github.com/macrox-pro/route-deck/internal/domain"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// GitHubRepoURL is the Help → RouteDeck on GitHub target (Browser.OpenURL).
const GitHubRepoURL = "https://github.com/macrox-pro/route-deck"

// BridgeReader returns the last known bridge state for Stop enablement.
type BridgeReader interface {
	BridgeState() domain.BridgeState
}

// AppMenuDeps wires menu actions without inventing Wails surfaces.
type AppMenuDeps struct {
	App       *application.App
	Emit      Emitter
	Editors   *EditorDialogs
	EditorSvc *EditorDialogsService
	Utilities *UtilityWindowsService
	Bridge    BridgeReader
	Parent    *application.WebviewWindow
	Logger    *slog.Logger
}

// AppMenu owns the RouteDeck application menu (ADR-0020).
type AppMenu struct {
	deps AppMenuDeps

	applyItem     *application.MenuItem
	stopItem      *application.MenuItem
	recheckItem   *application.MenuItem
	addAppItem    *application.MenuItem
	viewItems     []*application.MenuItem
	toggleSidebar *application.MenuItem
}

// MenuItemEnabled reports whether a main-window command stays enabled (ADR-0020).
// Stop also requires CLIOwned; pass cliOwned for command "stop".
func MenuItemEnabled(command string, editorOpen, cliOwned bool) bool {
	switch command {
	case "about", "settings", "new-profile", "new-proxy", "setup-guides", "github", "quit",
		"hide", "hide-others", "show-all", "close-window", "enter-full-screen":
		return true
	case "stop":
		return !editorOpen && cliOwned
	case "apply", "recheck", "add-application", "navigate", "toggle-sidebar":
		return !editorOpen
	default:
		return !editorOpen
	}
}

// Install builds the menu, sets it as the application menu, and syncs enable rules.
func (m *AppMenu) Install() {
	if m == nil || m.deps.App == nil {
		return
	}
	menu := m.build()
	m.deps.App.Menu.SetApplicationMenu(menu)
	m.SyncEnabled()
	if m.deps.Editors != nil {
		m.deps.Editors.OnOpenChange = func(_ bool) {
			m.SyncEnabled()
		}
	}
}

// SyncEnabled updates Protection / View / Add Application enablement.
func (m *AppMenu) SyncEnabled() {
	if m == nil {
		return
	}
	editorOpen := m.deps.Editors != nil && m.deps.Editors.AnyOpen()
	cliOwned := false
	if m.deps.Bridge != nil {
		cliOwned = tray.CLIOwned(m.deps.Bridge.BridgeState())
	}
	set := func(item *application.MenuItem, command string) {
		if item == nil {
			return
		}
		item.SetEnabled(MenuItemEnabled(command, editorOpen, cliOwned))
	}
	set(m.applyItem, "apply")
	set(m.stopItem, "stop")
	set(m.recheckItem, "recheck")
	set(m.addAppItem, "add-application")
	set(m.toggleSidebar, "toggle-sidebar")
	for _, item := range m.viewItems {
		set(item, "navigate")
	}
}

func (m *AppMenu) logger() *slog.Logger {
	if m != nil && m.deps.Logger != nil {
		return m.deps.Logger
	}
	return slog.Default()
}

func (m *AppMenu) emitCommand(command string, extra map[string]any) {
	payload := map[string]any{"command": command}
	for k, v := range extra {
		payload[k] = v
	}
	emit(m.deps.Emit, EventMainCommand, payload)
	if m.deps.Parent != nil {
		m.deps.Parent.Show().Focus()
	}
}

func (m *AppMenu) navigate(route string) {
	m.emitCommand("navigate", map[string]any{"route": route})
}

func (m *AppMenu) build() *application.Menu {
	root := application.NewMenu()

	appMenu := root.AddSubmenu("RouteDeck")
	appMenu.Add("About RouteDeck").OnClick(func(_ *application.Context) {
		if m.deps.Utilities != nil {
			_ = m.deps.Utilities.OpenAbout()
		}
	})
	appMenu.AddSeparator()
	appMenu.Add("Settings…").SetAccelerator("CmdOrCtrl+,").OnClick(func(_ *application.Context) {
		m.navigate("/settings")
	})
	appMenu.AddSeparator()
	appMenu.Add("Hide RouteDeck").SetAccelerator("CmdOrCtrl+h").SetRole(application.Hide)
	appMenu.Add("Hide Others").SetAccelerator("CmdOrCtrl+OptionOrAlt+h").SetRole(application.HideOthers)
	appMenu.Add("Show All").SetRole(application.UnHide)
	appMenu.AddSeparator()
	appMenu.Add("Quit RouteDeck").SetAccelerator("CmdOrCtrl+q").OnClick(func(_ *application.Context) {
		if m.deps.App != nil {
			m.deps.App.Quit()
		}
	})

	fileMenu := root.AddSubmenu("File")
	fileMenu.Add("New Profile").SetAccelerator("CmdOrCtrl+n").OnClick(func(_ *application.Context) {
		if m.deps.EditorSvc != nil {
			_ = m.deps.EditorSvc.OpenProfile("")
		}
	})
	fileMenu.Add("New Proxy").SetAccelerator("CmdOrCtrl+Shift+n").OnClick(func(_ *application.Context) {
		if m.deps.EditorSvc != nil {
			_ = m.deps.EditorSvc.OpenProxy("", "")
		}
	})
	m.addAppItem = fileMenu.Add("Add Application…").SetAccelerator("CmdOrCtrl+o").OnClick(func(_ *application.Context) {
		m.emitCommand("add-application", nil)
	})
	fileMenu.AddSeparator()
	fileMenu.AddRole(application.CloseWindow)

	root.AddRole(application.EditMenu)

	prot := root.AddSubmenu("Protection")
	m.applyItem = prot.Add("Apply").SetAccelerator("CmdOrCtrl+Return").OnClick(func(_ *application.Context) {
		m.emitCommand("apply", nil)
	})
	m.stopItem = prot.Add("Stop").SetAccelerator("CmdOrCtrl+.").OnClick(func(_ *application.Context) {
		m.emitCommand("stop", nil)
	})
	m.recheckItem = prot.Add("Recheck Health").SetAccelerator("CmdOrCtrl+r").OnClick(func(_ *application.Context) {
		m.emitCommand("recheck", nil)
	})

	view := root.AddSubmenu("View")
	viewRoutes := []struct {
		label string
		route string
		accel string
	}{
		{"Dashboard", "/", "CmdOrCtrl+1"},
		{"Applications", "/applications", "CmdOrCtrl+2"},
		{"Proxies", "/proxies", "CmdOrCtrl+3"},
		{"Profiles", "/profiles", "CmdOrCtrl+4"},
		{"Activity", "/activity", "CmdOrCtrl+5"},
	}
	m.viewItems = m.viewItems[:0]
	for _, row := range viewRoutes {
		route := row.route
		item := view.Add(row.label).SetAccelerator(row.accel).OnClick(func(_ *application.Context) {
			m.navigate(route)
		})
		m.viewItems = append(m.viewItems, item)
	}
	view.AddSeparator()
	m.toggleSidebar = view.Add("Toggle Sidebar").SetAccelerator("Ctrl+Command+S").OnClick(func(_ *application.Context) {
		m.emitCommand("toggle-sidebar", nil)
	})
	view.AddSeparator()
	view.Add("Enter Full Screen").SetAccelerator("Ctrl+Command+F").OnClick(func(_ *application.Context) {
		if m.deps.App == nil {
			return
		}
		current := m.deps.App.Window.Current()
		if current != nil {
			current.ToggleFullscreen()
		}
	})

	root.AddRole(application.WindowMenu)

	help := root.AddSubmenu("Help")
	help.Add("Setup Guides").OnClick(func(_ *application.Context) {
		if m.deps.Utilities != nil {
			_ = m.deps.Utilities.OpenGuides("")
		}
	})
	help.Add("RouteDeck on GitHub").OnClick(func(_ *application.Context) {
		if m.deps.App == nil {
			return
		}
		if err := m.deps.App.Browser.OpenURL(GitHubRepoURL); err != nil {
			m.logger().Warn("open github url failed", "component", "wails.app_menu", "err", err)
		}
	})

	return root
}

// NewAppMenu constructs an AppMenu from deps.
func NewAppMenu(deps AppMenuDeps) *AppMenu {
	return &AppMenu{deps: deps}
}
