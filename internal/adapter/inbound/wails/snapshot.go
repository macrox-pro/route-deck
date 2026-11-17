// Package wails is the inbound Wails adapter (DTO and services).
package wails

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProxyRowDTO is one proxy in the snapshot payload.
type ProxyRowDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Health string `json:"health"`
}

// SnapshotDTO is the frontend snapshot. Domain types are not exposed.
type SnapshotDTO struct {
	Bridge                    string        `json:"bridge"`
	Protection                string        `json:"protection"`
	VersionKnown              bool          `json:"versionKnown"`
	Version                   string        `json:"version"`
	BridgePath                string        `json:"bridgePath"`
	CompatibilityWarning      string        `json:"compatibilityWarning"`
	DesiredSet                bool          `json:"desiredSet"`
	DesiredID                 string        `json:"desiredId"`
	AppliedID                 string        `json:"appliedId"`
	AppliedConfirmed          bool          `json:"appliedConfirmed"`
	DesiredDiffersFromApplied bool          `json:"desiredDiffersFromApplied"`
	RecheckPending            bool          `json:"recheckPending"`
	DatabaseError             string        `json:"databaseError"`
	Proxies                   []ProxyRowDTO `json:"proxies"`
}

// SnapshotReader is the application snapshot use case.
type SnapshotReader interface {
	Snapshot(ctx context.Context) (service.SnapshotResult, error)
}

// SnapshotService is the Wails-bound snapshot service.
type SnapshotService struct {
	Reader SnapshotReader
	Logger *slog.Logger
}

func (s *SnapshotService) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// Get returns the current snapshot DTO.
func (s *SnapshotService) Get(ctx context.Context) (SnapshotDTO, error) {
	if s == nil || s.Reader == nil {
		return SnapshotDTO{}, fmt.Errorf("snapshot service: reader not configured")
	}
	s.logger().Debug("snapshot get", "component", "wails.snapshot")
	result, err := s.Reader.Snapshot(ctx)
	if err != nil {
		s.logger().Error("snapshot get failed", "component", "wails.snapshot", "err", err)
		return SnapshotDTO{}, fmt.Errorf("snapshot get: %w", err)
	}
	return mapSnapshot(result), nil
}

func mapSnapshot(result service.SnapshotResult) SnapshotDTO {
	proxies := make([]ProxyRowDTO, 0, len(result.Proxies))
	for _, row := range result.Proxies {
		proxies = append(proxies, ProxyRowDTO{
			ID:     string(row.ID),
			Name:   row.Name,
			Kind:   string(row.Kind),
			Host:   row.Host,
			Port:   row.Port,
			Health: string(row.Health),
		})
	}
	desiredID := ""
	if result.DesiredID != domain.ProfileID("") {
		desiredID = string(result.DesiredID)
	}
	appliedID := ""
	if result.AppliedID != domain.ProfileID("") {
		appliedID = string(result.AppliedID)
	}
	return SnapshotDTO{
		Bridge:                    string(result.Bridge),
		Protection:                string(result.Protection),
		VersionKnown:              result.VersionKnown,
		Version:                   result.Version,
		BridgePath:                result.BridgePath,
		CompatibilityWarning:      result.CompatibilityWarning,
		DesiredSet:                result.DesiredSet,
		DesiredID:                 desiredID,
		AppliedID:                 appliedID,
		AppliedConfirmed:          result.AppliedConfirmed,
		DesiredDiffersFromApplied: result.DesiredDiffersFromApplied,
		RecheckPending:            result.RecheckPending,
		DatabaseError:             result.DatabaseError,
		Proxies:                   proxies,
	}
}
