package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

const activityPageSize = 100

// Activity loads recent RouteDeck JSON log records for the Activity screen.
type Activity struct {
	Reader port.ActivityLogReader
	Path   string
	Logger *slog.Logger
}

func (a *Activity) logger() *slog.Logger {
	if a != nil && a.Logger != nil {
		return a.Logger
	}
	return slog.Default()
}

// Recent returns one page of up to 100 records. before is an exclusive RFC3339
// cursor; empty means the newest page.
func (a *Activity) Recent(ctx context.Context, before string) (port.ActivityPage, error) {
	if a == nil || a.Reader == nil {
		return port.ActivityPage{}, fmt.Errorf("activity: not configured")
	}
	if a.Path == "" {
		return port.ActivityPage{}, fmt.Errorf("activity: empty log path")
	}
	page, err := a.Reader.Recent(ctx, a.Path, activityPageSize, before)
	if err != nil {
		a.logger().Debug("activity recent failed", "component", "activity", "err", err)
		return port.ActivityPage{}, fmt.Errorf("activity recent: %w", err)
	}
	return page, nil
}
