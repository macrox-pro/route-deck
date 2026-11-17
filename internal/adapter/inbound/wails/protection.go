package wails

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ProtectionCoordinator is the application Apply/Stop surface.
type ProtectionCoordinator interface {
	Apply(ctx context.Context, profileID domain.ProfileID, operationID string) error
	Stop(ctx context.Context, operationID string) error
}

// SettingsStoreReader loads desired profile bookkeeping for Apply.
type SettingsStoreReader interface {
	Get(ctx context.Context) (desiredProfileID domain.ProfileID, desiredSet bool, err error)
}

// ProtectionService exposes Apply/Stop over Wails.
type ProtectionService struct {
	Coordinator ProtectionCoordinator
	Settings    SettingsStoreReader
	NewID       func() (string, error)
	Emit        Emitter
	Logger      *slog.Logger
}

func (s *ProtectionService) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// ApplyResultDTO is returned immediately; confirmation arrives via events + snapshot.
type ApplyResultDTO struct {
	OperationID string `json:"operationId"`
}

// Apply starts Apply for the desired profile and returns an operation id.
func (s *ProtectionService) Apply(ctx context.Context) (ApplyResultDTO, error) {
	if s == nil || s.Coordinator == nil || s.Settings == nil {
		return ApplyResultDTO{}, fmt.Errorf("protection service: not configured")
	}
	desiredID, desiredSet, err := s.Settings.Get(ctx)
	if err != nil {
		s.logger().Info("apply settings failed", "component", "wails.protection", "err", err)
		return ApplyResultDTO{}, fmt.Errorf("protection apply: %w", err)
	}
	if !desiredSet || desiredID == "" {
		return ApplyResultDTO{}, fmt.Errorf("protection apply: desired profile not set")
	}
	opID, err := s.newOperationID()
	if err != nil {
		return ApplyResultDTO{}, err
	}
	s.logger().Info("apply requested", "component", "wails.protection", "operation_id", opID, "profile_id", string(desiredID))
	emit(s.Emit, EventOperationProgress, map[string]any{
		"operation_id": opID,
		"phase":        "start",
	})
	go func() {
		applyCtx := context.Background()
		if err := s.Coordinator.Apply(applyCtx, desiredID, opID); err != nil {
			s.logger().Info("apply failed", "component", "wails.protection", "operation_id", opID, "err", err)
			emit(s.Emit, EventOperationFailed, map[string]any{
				"operation_id": opID,
				"error":        err.Error(),
			})
			return
		}
		emit(s.Emit, EventOperationProgress, map[string]any{
			"operation_id": opID,
			"phase":        "confirmed",
		})
		emit(s.Emit, EventProtectionStateChanged, map[string]any{
			"operation_id": opID,
		})
	}()
	return ApplyResultDTO{OperationID: opID}, nil
}

// Stop stops the managed CLI.
func (s *ProtectionService) Stop(ctx context.Context) (ApplyResultDTO, error) {
	if s == nil || s.Coordinator == nil {
		return ApplyResultDTO{}, fmt.Errorf("protection service: not configured")
	}
	opID, err := s.newOperationID()
	if err != nil {
		return ApplyResultDTO{}, err
	}
	s.logger().Info("stop requested", "component", "wails.protection", "operation_id", opID)
	emit(s.Emit, EventOperationProgress, map[string]any{
		"operation_id": opID,
		"phase":        "stop_start",
	})
	if err := s.Coordinator.Stop(ctx, opID); err != nil {
		s.logger().Info("stop failed", "component", "wails.protection", "operation_id", opID, "err", err)
		emit(s.Emit, EventOperationFailed, map[string]any{
			"operation_id": opID,
			"error":        err.Error(),
		})
		return ApplyResultDTO{}, fmt.Errorf("protection stop: %w", err)
	}
	emit(s.Emit, EventOperationProgress, map[string]any{
		"operation_id": opID,
		"phase":        "stopped",
	})
	emit(s.Emit, EventProtectionStateChanged, map[string]any{
		"operation_id": opID,
	})
	emit(s.Emit, EventProxyBridgeStateChanged, map[string]any{
		"operation_id": opID,
	})
	return ApplyResultDTO{OperationID: opID}, nil
}

func (s *ProtectionService) newOperationID() (string, error) {
	if s.NewID != nil {
		return s.NewID()
	}
	return service.NewID()
}
