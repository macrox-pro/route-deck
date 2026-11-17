package logfile

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

var primaryActivityKeys = map[string]struct{}{
	"time":  {},
	"level": {},
	"msg":   {},
}

var droppedActivityKeys = map[string]struct{}{
	"password": {},
	"secret":   {},
	"token":    {},
	"profile":  {},
}

// Recent reads path and rotations .1 .. DefaultKeep, returning up to limit newest
// complete JSON records. If before is non-empty, only records strictly older than
// that RFC3339 timestamp are included. A missing or unreadable active path is an
// error. Broken and non-JSON lines are skipped. limit <= 0 returns an empty page.
func Recent(ctx context.Context, path string, limit int, before string) (port.ActivityPage, error) {
	if err := ctx.Err(); err != nil {
		return port.ActivityPage{}, fmt.Errorf("logfile recent: %w", err)
	}
	if path == "" {
		return port.ActivityPage{}, fmt.Errorf("logfile recent: empty path")
	}
	if limit <= 0 {
		return port.ActivityPage{}, nil
	}
	var beforeTime time.Time
	if before != "" {
		parsed, ok := parseActivityTime(before)
		if !ok {
			return port.ActivityPage{}, fmt.Errorf("logfile recent: invalid before %q", before)
		}
		beforeTime = parsed
	}
	if _, err := os.Stat(path); err != nil {
		return port.ActivityPage{}, fmt.Errorf("logfile recent: %w", err)
	}

	var records []port.ActivityRecord
	for i := DefaultKeep; i >= 1; i-- {
		rotated := fmt.Sprintf("%s.%d", path, i)
		batch, err := readActivityFile(ctx, rotated)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return port.ActivityPage{}, err
		}
		records = append(records, batch...)
	}
	batch, err := readActivityFile(ctx, path)
	if err != nil {
		return port.ActivityPage{}, err
	}
	records = append(records, batch...)

	if !beforeTime.IsZero() {
		filtered := make([]port.ActivityRecord, 0, len(records))
		for _, rec := range records {
			t, ok := parseActivityTime(rec.Time)
			if !ok {
				continue
			}
			if t.Before(beforeTime) {
				filtered = append(filtered, rec)
			}
		}
		records = filtered
	}

	sort.SliceStable(records, func(i, j int) bool {
		ti, okI := parseActivityTime(records[i].Time)
		tj, okJ := parseActivityTime(records[j].Time)
		switch {
		case okI && okJ:
			return ti.After(tj)
		case okI:
			return true
		case okJ:
			return false
		default:
			return i > j
		}
	})

	hasMore := len(records) > limit
	if hasMore {
		records = records[:limit]
	}
	return port.ActivityPage{Records: records, HasMore: hasMore}, nil
}

// Reader implements port.ActivityLogReader.
type Reader struct{}

var _ port.ActivityLogReader = (*Reader)(nil)

// Recent implements port.ActivityLogReader.
func (r *Reader) Recent(ctx context.Context, path string, limit int, before string) (port.ActivityPage, error) {
	_ = r
	return Recent(ctx, path, limit, before)
}

func readActivityFile(ctx context.Context, path string) ([]port.ActivityRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("logfile recent: %w", err)
	}
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var out []port.ActivityRecord
	scanner := bufio.NewScanner(file)
	// slog JSON lines can be long; raise the default 64 KiB token size.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("logfile recent: %w", err)
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		rec, ok := parseActivityLine(line)
		if !ok {
			continue
		}
		out = append(out, rec)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("logfile recent scan %s: %w", path, err)
	}
	return out, nil
}

func parseActivityLine(line string) (port.ActivityRecord, bool) {
	var raw map[string]any
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return port.ActivityRecord{}, false
	}
	for key := range droppedActivityKeys {
		delete(raw, key)
	}

	strings := make(map[string]string, len(raw))
	for key, v := range raw {
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		strings[key] = s
	}
	if strings["time"] == "" && strings["msg"] == "" && strings["level"] == "" {
		return port.ActivityRecord{}, false
	}

	attrs := make([]port.ActivityAttr, 0, len(strings))
	for key, value := range strings {
		if _, primary := primaryActivityKeys[key]; primary {
			continue
		}
		attrs = append(attrs, port.ActivityAttr{Key: key, Value: value})
	}
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})

	return port.ActivityRecord{
		Time:            strings["time"],
		Level:           strings["level"],
		Msg:             strings["msg"],
		Component:       strings["component"],
		ProxyID:         strings["proxy_id"],
		ProfileID:       strings["profile_id"],
		BridgeState:     strings["bridge_state"],
		ProtectionState: strings["protection_state"],
		Health:          strings["health"],
		LineClass:       strings["line_class"],
		OperationID:     strings["operation_id"],
		Err:             strings["err"],
		Attrs:           attrs,
	}, true
}

func parseActivityTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}
	return time.Time{}, false
}
