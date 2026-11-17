package port

import "context"

// InstalledApplication is one browsable installed application (not yet in the catalog).
type InstalledApplication struct {
	DisplayName string
	BundleID    string
	Path        string
	// IsSystem is set by the platform adapter. On macOS: true when the path is
	// outside /Applications and ~/Applications. Other OS adapters define their
	// own user-facing roots; the UI must not hard-code OS paths.
	IsSystem bool
}

// InstalledApplicationLister lists applications registered on the platform.
type InstalledApplicationLister interface {
	List(ctx context.Context) ([]InstalledApplication, error)
}

// ApplicationIconReader returns a PNG icon for an application path.
type ApplicationIconReader interface {
	PNG(ctx context.Context, path string) ([]byte, error)
}
