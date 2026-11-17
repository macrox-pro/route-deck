package wails

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// UserDialog shows native message dialogs to the user.
type UserDialog interface {
	Info(msg UserMsg)
	Warning(msg UserMsg)
}

// NativeUserDialog presents Wails Info/Warning dialogs attached to the main window.
type NativeUserDialog struct {
	App    *application.App
	Window *application.WebviewWindow
}

// Info shows a native information dialog with an OK button.
func (d *NativeUserDialog) Info(msg UserMsg) {
	d.show(d.App.Dialog.Info(), msg)
}

// Warning shows a native warning dialog with an OK button.
func (d *NativeUserDialog) Warning(msg UserMsg) {
	d.show(d.App.Dialog.Warning(), msg)
}

func (d *NativeUserDialog) show(dialog *application.MessageDialog, msg UserMsg) {
	if d == nil || d.App == nil || dialog == nil {
		return
	}
	dialog.SetTitle(msg.Title).SetMessage(msg.Body)
	if d.Window != nil {
		dialog.AttachToWindow(d.Window)
	}
	ok := dialog.AddButton(DialogOK).SetAsDefault()
	ok.OnClick(func() {})
	dialog.Show()
}
