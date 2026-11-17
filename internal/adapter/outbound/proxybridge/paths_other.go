//go:build !darwin

package proxybridge

// DefaultDocumentedPaths is empty until a non-macOS adapter ships documented
// install locations for that OS. Locate still honors a custom Settings path.
func DefaultDocumentedPaths() []string {
	return nil
}
