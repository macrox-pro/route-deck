package wails

// UserMsg is user-visible dialog copy with a stable ID for future i18n.
type UserMsg struct {
	ID    string
	Title string
	Body  string
}

// English dialog copy until a locale pack lands. Keep IDs stable.
var (
	UserMsgGUIConflict = UserMsg{
		ID:    "protection.gui_conflict",
		Title: "ProxyBridge app is running",
		Body:  "Quit the ProxyBridge app, then Apply again. RouteDeck cannot drive the CLI while the GUI holds the network extension.",
	}
	UserMsgExtensionApproval = UserMsg{
		ID:    "protection.extension_approval",
		Title: "Approve the network extension",
		Body:  "Open System Settings → Privacy & Security and click Allow, then enter your password. This is needed once. RouteDeck will keep waiting until you approve or the request times out.",
	}
)

// DialogOK is the dismiss button label (English until i18n).
const DialogOK = "OK"
