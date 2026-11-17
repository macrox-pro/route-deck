package wails

import (
	"time"
)

// Event names from docs/contracts/events.md.
const (
	EventProtectionStateChanged  = "protection.state.changed"
	EventProxyHealthChanged      = "proxy.health.changed"
	EventProxyBridgeStateChanged = "proxybridge.state.changed"
	EventOperationProgress       = "operation.progress"
	EventOperationFailed         = "operation.failed"
)

// EventEnvelope is the product event shape.
type EventEnvelope struct {
	Schema    int            `json:"schema"`
	Timestamp string         `json:"timestamp"`
	Payload   map[string]any `json:"payload"`
}

// Emitter sends a named product event to the UI.
type Emitter func(name string, envelope EventEnvelope)

func emit(emitter Emitter, name string, payload map[string]any) {
	if emitter == nil {
		return
	}
	emitter(name, EventEnvelope{
		Schema:    1,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Payload:   payload,
	})
}
