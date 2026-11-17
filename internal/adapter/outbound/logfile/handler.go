package logfile

import (
	"log/slog"
	"strings"
)

var redactedKeys = map[string]struct{}{
	"password": {},
	"secret":   {},
	"token":    {},
	"profile":  {},
}

// RedactAttr replaces sensitive attribute values with [REDACTED].
func RedactAttr(_ []string, a slog.Attr) slog.Attr {
	key := strings.ToLower(a.Key)
	if _, ok := redactedKeys[key]; ok {
		return slog.String(a.Key, "[REDACTED]")
	}
	return a
}
