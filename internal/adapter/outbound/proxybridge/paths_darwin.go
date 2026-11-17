//go:build darwin

package proxybridge

// DefaultDocumentedPaths are the official ProxyBridge install locations on macOS
// (ADR-0008 / RD-005). Other platforms supply their own via paths_*.go.
func DefaultDocumentedPaths() []string {
	return []string{
		"/usr/local/bin/proxybridge",
		"/Applications/ProxyBridge.app/Contents/MacOS/ProxyBridge",
	}
}
