package port

import "context"

// LocalPortProber reports whether a TCP port accepts connections on 127.0.0.1.
// Callers MUST pass only catalog default ports (ADR-0013).
type LocalPortProber interface {
	Open(ctx context.Context, port int) (bool, error)
}
