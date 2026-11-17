// Package wails is the inbound Wails adapter (DTO and services).
package wails

import (
	"context"
	"fmt"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProxyRowDTO is one proxy in the snapshot payload.
type ProxyRowDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Health string `json:"health"`
}

// SnapshotDTO is the frontend snapshot. Domain types are not exposed.
type SnapshotDTO struct {
	Bridge           string        `json:"bridge"`
	Protection       string        `json:"protection"`
	VersionKnown     bool          `json:"versionKnown"`
	DesiredSet       bool          `json:"desiredSet"`
	DesiredID        string        `json:"desiredId"`
	AppliedID        string        `json:"appliedId"`
	AppliedConfirmed bool          `json:"appliedConfirmed"`
	DatabaseError    string        `json:"databaseError"`
	Proxies          []ProxyRowDTO `json:"proxies"`
}

// SnapshotReader is the application snapshot use case.
type SnapshotReader interface {
	Snapshot(ctx context.Context) (service.SnapshotResult, error)
}

// SnapshotService is the Wails-bound snapshot service.
type SnapshotService struct {
	Reader SnapshotReader
}

// Get returns the current snapshot DTO.
func (s *SnapshotService) Get(ctx context.Context) (SnapshotDTO, error) {
	if s == nil || s.Reader == nil {
		return SnapshotDTO{}, fmt.Errorf("snapshot service: reader not configured")
	}
	result, err := s.Reader.Snapshot(ctx)
	if err != nil {
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
		Bridge:           string(result.Bridge),
		Protection:       string(result.Protection),
		VersionKnown:     result.VersionKnown,
		DesiredSet:       result.DesiredSet,
		DesiredID:        desiredID,
		AppliedID:        appliedID,
		AppliedConfirmed: result.AppliedConfirmed,
		DatabaseError:    result.DatabaseError,
		Proxies:          proxies,
	}
}
