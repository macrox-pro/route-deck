package wails

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
)

// ActivityService exposes recent JSON log records to the UI.
type ActivityService struct {
	Activity *service.Activity
	Logger   *slog.Logger
}

// ActivityAttrDTO is one non-secret string attribute.
type ActivityAttrDTO struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ActivityRecordDTO is one Activity row.
type ActivityRecordDTO struct {
	Time            string            `json:"time"`
	Level           string            `json:"level"`
	Msg             string            `json:"msg"`
	Component       string            `json:"component"`
	ProxyID         string            `json:"proxyId"`
	ProfileID       string            `json:"profileId"`
	BridgeState     string            `json:"bridgeState"`
	ProtectionState string            `json:"protectionState"`
	Health          string            `json:"health"`
	LineClass       string            `json:"lineClass"`
	OperationID     string            `json:"operationId"`
	Err             string            `json:"err"`
	Attrs           []ActivityAttrDTO `json:"attrs"`
}

// ActivityRecentDTO is one page of Activity records.
type ActivityRecentDTO struct {
	Records []ActivityRecordDTO `json:"records"`
	HasMore bool                `json:"hasMore"`
}

func (s *ActivityService) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// Recent returns up to 100 newest complete log records.
// before is an exclusive RFC3339 cursor; empty means the newest page.
func (s *ActivityService) Recent(ctx context.Context, before string) (ActivityRecentDTO, error) {
	if s == nil || s.Activity == nil {
		return ActivityRecentDTO{}, fmt.Errorf("activity service: not configured")
	}
	page, err := s.Activity.Recent(ctx, before)
	if err != nil {
		s.logger().Debug("activity recent failed", "component", "wails", "err", err)
		return ActivityRecentDTO{}, err
	}
	out := make([]ActivityRecordDTO, 0, len(page.Records))
	for _, rec := range page.Records {
		out = append(out, activityRecordDTO(rec))
	}
	return ActivityRecentDTO{Records: out, HasMore: page.HasMore}, nil
}

func activityRecordDTO(rec port.ActivityRecord) ActivityRecordDTO {
	attrs := make([]ActivityAttrDTO, 0, len(rec.Attrs))
	for _, attr := range rec.Attrs {
		attrs = append(attrs, ActivityAttrDTO{Key: attr.Key, Value: attr.Value})
	}
	return ActivityRecordDTO{
		Time:            rec.Time,
		Level:           rec.Level,
		Msg:             rec.Msg,
		Component:       rec.Component,
		ProxyID:         rec.ProxyID,
		ProfileID:       rec.ProfileID,
		BridgeState:     rec.BridgeState,
		ProtectionState: rec.ProtectionState,
		Health:          rec.Health,
		LineClass:       rec.LineClass,
		OperationID:     rec.OperationID,
		Err:             rec.Err,
		Attrs:           attrs,
	}
}
