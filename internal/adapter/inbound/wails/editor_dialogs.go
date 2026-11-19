package wails

import (
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Editor dialog kinds (ADR-0014).
const (
	DialogKindProfile = "profile"
	DialogKindProxy   = "proxy"
	DialogKindHelpers = "helpers"
)

// EditorDialogs manages reused native-frame editor windows.
type EditorDialogs struct {
	App    *application.App
	Parent *application.WebviewWindow
	Logger *slog.Logger
	Emit   Emitter

	mu              sync.Mutex
	windows         map[string]*application.WebviewWindow
	dirty           map[string]bool
	asking          map[string]bool // discard Question already open for kind
	appearance      string          // system|light|dark for create-time chrome
	saveForQuitWait map[string]struct{}
	saveForQuitDone chan error
}

// SetAppearance stores the preference used when creating editor windows.
func (d *EditorDialogs) SetAppearance(appearance string) {
	if d == nil {
		return
	}
	d.mu.Lock()
	d.appearance = appearance
	d.mu.Unlock()
}

// HasDirty reports whether any open editor has unsaved edits.
func (d *EditorDialogs) HasDirty() bool {
	if d == nil {
		return false
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ensureMaps()
	for _, dirty := range d.dirty {
		if dirty {
			return true
		}
	}
	return false
}

// DirtySaveableKinds returns dirty profile/proxy kinds (ADR-0018 Save and Quit).
func (d *EditorDialogs) DirtySaveableKinds() []string {
	if d == nil {
		return nil
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ensureMaps()
	out := make([]string, 0, 2)
	for _, kind := range []string{DialogKindProfile, DialogKindProxy} {
		if d.dirty[kind] {
			out = append(out, kind)
		}
	}
	return out
}

// RunSaveForQuit emits editors.save_for_quit and waits for AckSaveForQuit from each dirty saveable kind.
func (d *EditorDialogs) RunSaveForQuit() error {
	if d == nil {
		return nil
	}
	kinds := d.DirtySaveableKinds()
	if len(kinds) == 0 {
		return nil
	}
	d.mu.Lock()
	if d.saveForQuitDone != nil {
		d.mu.Unlock()
		return fmt.Errorf("editor dialogs: save for quit already in progress")
	}
	d.saveForQuitWait = make(map[string]struct{}, len(kinds))
	for _, kind := range kinds {
		d.saveForQuitWait[kind] = struct{}{}
	}
	done := make(chan error, 1)
	d.saveForQuitDone = done
	d.mu.Unlock()

	emit(d.Emit, EventEditorsSaveForQuit, map[string]any{"kinds": kinds})

	timer := time.NewTimer(2 * time.Minute)
	defer timer.Stop()
	select {
	case err := <-done:
		return err
	case <-timer.C:
		d.mu.Lock()
		d.saveForQuitWait = nil
		d.saveForQuitDone = nil
		d.mu.Unlock()
		return fmt.Errorf("editor dialogs: save for quit timed out")
	}
}

// AckSaveForQuit records one kind's save result during Save and Quit (ADR-0018).
func (d *EditorDialogs) AckSaveForQuit(kind, errMsg string) {
	if d == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.saveForQuitDone == nil || d.saveForQuitWait == nil {
		return
	}
	if _, ok := d.saveForQuitWait[kind]; !ok {
		return
	}
	delete(d.saveForQuitWait, kind)
	if errMsg != "" {
		done := d.saveForQuitDone
		d.saveForQuitWait = nil
		d.saveForQuitDone = nil
		done <- fmt.Errorf("%s", errMsg)
		return
	}
	if len(d.saveForQuitWait) > 0 {
		return
	}
	done := d.saveForQuitDone
	d.saveForQuitDone = nil
	done <- nil
}

// CloseAllForced force-closes every editor window without a discard prompt.
func (d *EditorDialogs) CloseAllForced() {
	if d == nil {
		return
	}
	d.mu.Lock()
	d.ensureMaps()
	kinds := make([]string, 0, len(d.windows))
	for kind := range d.windows {
		kinds = append(kinds, kind)
	}
	d.mu.Unlock()
	for _, kind := range kinds {
		d.close(kind, true)
	}
}

func (d *EditorDialogs) ensureMaps() {
	if d.windows == nil {
		d.windows = make(map[string]*application.WebviewWindow)
	}
	if d.dirty == nil {
		d.dirty = make(map[string]bool)
	}
	if d.asking == nil {
		d.asking = make(map[string]bool)
	}
}

// FocusOpen focuses any open editor window. Returns true when one was focused.
func (d *EditorDialogs) FocusOpen() bool {
	if d == nil {
		return false
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, kind := range []string{DialogKindProfile, DialogKindProxy, DialogKindHelpers} {
		w := d.windows[kind]
		if w == nil {
			continue
		}
		w.Show().Focus()
		return true
	}
	return false
}

func (d *EditorDialogs) open(kind, title, hashPath string, width, height int, resizable bool) error {
	if d == nil || d.App == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ensureMaps()

	if w := d.windows[kind]; w != nil {
		w.SetTitle(title)
		w.SetURL("/#" + hashPath)
		w.Show().Focus()
		d.disableParentLocked()
		return nil
	}

	appearance := d.appearance

	opts := application.WebviewWindowOptions{
		Name:          "dialog-" + kind,
		Title:         title,
		Width:         width,
		Height:        height,
		MinWidth:      width / 2,
		MinHeight:     height / 2,
		DisableResize: !resizable,
		Hidden:        true,
		URL:           "/#" + hashPath,
	}
	if resizable {
		opts.MinWidth = 600
		opts.MinHeight = 480
	}
	opts = DialogWindowOptions(appearance, opts)
	win := d.App.Window.NewWithOptions(opts)
	d.windows[kind] = win
	d.dirty[kind] = false

	win.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		d.mu.Lock()
		isDirty := d.dirty[kind]
		alreadyAsking := d.asking[kind]
		d.mu.Unlock()
		if !isDirty {
			return
		}
		// Cancel first so listeners do not tear down while the question is open.
		e.Cancel()
		if alreadyAsking {
			return
		}
		app := application.Get()
		if app == nil {
			return
		}
		d.mu.Lock()
		d.asking[kind] = true
		d.mu.Unlock()

		// App-modal (no AttachToWindow): a sheet attached to a window is async on
		// macOS, so Show() returns before OnClick and a post-Show Close never runs.
		dialog := app.Dialog.Question().
			SetTitle("Discard changes?").
			SetMessage("Discard changes?")
		cancel := dialog.AddButton("Cancel").SetAsCancel()
		cancel.OnClick(func() {
			d.mu.Lock()
			d.asking[kind] = false
			d.mu.Unlock()
		})
		discard := dialog.AddButton("Discard").SetAsDefault()
		discard.OnClick(func() {
			d.mu.Lock()
			d.dirty[kind] = false
			d.asking[kind] = false
			d.mu.Unlock()
			win.Close()
		})
		dialog.Show()
	})

	win.OnWindowEvent(events.Common.WindowClosing, func(_ *application.WindowEvent) {
		d.mu.Lock()
		delete(d.windows, kind)
		delete(d.dirty, kind)
		delete(d.asking, kind)
		d.enableParentLocked()
		d.mu.Unlock()
	})

	if d.Parent != nil {
		d.Parent.AttachModal(win)
	}
	d.disableParentLocked()
	win.Show().Focus()
	return nil
}

func (d *EditorDialogs) disableParentLocked() {
	if d.Parent != nil {
		d.Parent.SetEnabled(false)
	}
}

func (d *EditorDialogs) enableParentLocked() {
	if d.Parent == nil {
		return
	}
	for _, w := range d.windows {
		if w != nil {
			return
		}
	}
	d.Parent.SetEnabled(true)
	d.Parent.Focus()
}

func (d *EditorDialogs) setDirty(kind string, dirty bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ensureMaps()
	d.dirty[kind] = dirty
}

func (d *EditorDialogs) close(kind string, force bool) {
	d.mu.Lock()
	w := d.windows[kind]
	if force {
		d.dirty[kind] = false
	}
	d.mu.Unlock()
	if w != nil {
		w.Close()
	}
}

func (d *EditorDialogs) setTitle(kind, title string) {
	d.mu.Lock()
	w := d.windows[kind]
	d.mu.Unlock()
	if w != nil && title != "" {
		w.SetTitle(title)
	}
}

// EditorDialogsService exposes dialog open/close to the frontend.
type EditorDialogsService struct {
	Dialogs *EditorDialogs
	Logger  *slog.Logger
}

// OpenProfile opens or focuses the profile editor (empty id = create).
func (s *EditorDialogsService) OpenProfile(id string) error {
	if s == nil || s.Dialogs == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	title := "New Profile"
	path := "/dialog/profile"
	if id != "" {
		title = "Edit Profile"
		path = "/dialog/profile/" + url.PathEscape(id)
	}
	return s.Dialogs.open(DialogKindProfile, title, path, 760, 640, true)
}

// OpenProxy opens or focuses the proxy editor.
func (s *EditorDialogsService) OpenProxy(id, guide string) error {
	if s == nil || s.Dialogs == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	title := "New Proxy"
	path := "/dialog/proxy"
	q := url.Values{}
	if id != "" {
		title = "Edit Proxy"
		path = "/dialog/proxy/" + url.PathEscape(id)
	}
	if guide != "" {
		q.Set("guide", guide)
	}
	if enc := q.Encode(); enc != "" {
		path += "?" + enc
	}
	return s.Dialogs.open(DialogKindProxy, title, path, 480, 560, false)
}

// OpenHelpers opens the helpers review dialog for an application.
func (s *EditorDialogsService) OpenHelpers(applicationID string) error {
	if s == nil || s.Dialogs == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	if applicationID == "" {
		return fmt.Errorf("editor dialogs: empty application id")
	}
	path := "/dialog/helpers/" + url.PathEscape(applicationID)
	return s.Dialogs.open(DialogKindHelpers, "Review Helpers", path, 600, 560, true)
}

// SetTitle updates the native title of an open dialog kind (ADR-0014 name suffixes).
func (s *EditorDialogsService) SetTitle(kind, title string) error {
	if s == nil || s.Dialogs == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	switch kind {
	case DialogKindProfile, DialogKindProxy, DialogKindHelpers:
	default:
		return fmt.Errorf("editor dialogs: unknown kind %q", kind)
	}
	s.Dialogs.setTitle(kind, title)
	return nil
}

// SetDirty records unsaved edits for a dialog kind.
func (s *EditorDialogsService) SetDirty(kind string, dirty bool) error {
	if s == nil || s.Dialogs == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	switch kind {
	case DialogKindProfile, DialogKindProxy, DialogKindHelpers:
	default:
		return fmt.Errorf("editor dialogs: unknown kind %q", kind)
	}
	s.Dialogs.setDirty(kind, dirty)
	return nil
}

// Close force-closes a dialog kind without a discard prompt (after save).
func (s *EditorDialogsService) Close(kind string) error {
	if s == nil || s.Dialogs == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	s.Dialogs.close(kind, true)
	return nil
}

// RequestClose closes via the native close path so a dirty editor asks "Discard changes?".
func (s *EditorDialogsService) RequestClose(kind string) error {
	if s == nil || s.Dialogs == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	switch kind {
	case DialogKindProfile, DialogKindProxy, DialogKindHelpers:
	default:
		return fmt.Errorf("editor dialogs: unknown kind %q", kind)
	}
	s.Dialogs.close(kind, false)
	return nil
}

// AckSaveForQuit acknowledges one editor's save during Quit Save and Quit (ADR-0018).
func (s *EditorDialogsService) AckSaveForQuit(kind, errMsg string) error {
	if s == nil || s.Dialogs == nil {
		return fmt.Errorf("editor dialogs: not configured")
	}
	switch kind {
	case DialogKindProfile, DialogKindProxy:
	default:
		return fmt.Errorf("editor dialogs: unknown kind %q", kind)
	}
	s.Dialogs.AckSaveForQuit(kind, errMsg)
	return nil
}
