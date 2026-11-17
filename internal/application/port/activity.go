package port

import "context"

// ActivityAttr is one non-secret string attribute from a JSON log line.
type ActivityAttr struct {
	Key   string
	Value string
}

// ActivityRecord is one redacted JSON log line for the Activity screen.
type ActivityRecord struct {
	Time            string
	Level           string
	Msg             string
	Component       string
	ProxyID         string
	ProfileID       string
	BridgeState     string
	ProtectionState string
	Health          string
	LineClass       string
	OperationID     string
	Err             string
	// Attrs holds every non-secret string attribute except time, level, and msg.
	Attrs []ActivityAttr
}

// ActivityPage is one page of Activity records (newest first).
type ActivityPage struct {
	Records []ActivityRecord
	HasMore bool
}

// ActivityLogReader reads recent complete records from a RouteDeck JSON log path.
// before is an exclusive RFC3339 cursor: empty means the newest page; otherwise
// only records strictly older than before are returned.
type ActivityLogReader interface {
	Recent(ctx context.Context, path string, limit int, before string) (ActivityPage, error)
}
