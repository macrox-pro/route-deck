package domain

// Application is a user-facing app entry (display name, bundle id, icon ref).
type Application struct {
	ID          ApplicationID
	DisplayName string
	BundleID    string
	IconRef     string
}
