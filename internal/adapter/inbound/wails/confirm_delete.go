package wails

import "github.com/wailsapp/wails/v3/pkg/application"

// confirmDelete shows a native Question. Cancel is the default. Returns true when the delete button is chosen.
func confirmDelete(title, message, deleteLabel string) bool {
	app := application.Get()
	if app == nil {
		return false
	}
	if deleteLabel == "" {
		deleteLabel = "Delete"
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
