package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ErrOperationInProgress is returned when Apply or Stop is already running.
var ErrOperationInProgress = errors.New("protection: operation in progress")

// Protection coordinates Apply/Stop for one managed bridge process.
type Protection struct {
	Profiles  port.ProfileStore
	Proxies   port.ProxyStore
	Processes port.ProcessIdentityStore
	Settings  port.SettingsStore
	Secrets   port.SecretStore
	Locator   port.BridgeLocator
	Compiler  port.ProfileCompiler
	Process   port.BridgeProcess
	Health    *HealthMonitor
	Logger    *slog.Logger

	mu               sync.Mutex
	busy             bool
	appliedID        domain.ProfileID
	appliedConfirmed bool
	lastAppliedJSON  []byte
	bridgeState      domain.BridgeState
	desiredDiffers   bool
}

func (p *Protection) logger() *slog.Logger {
	if p != nil && p.Logger != nil {
		return p.Logger
	}
	return slog.Default()
}

// Apply compiles and starts the desired profile. operationID is for logging/events.
func (p *Protection) Apply(ctx context.Context, profileID domain.ProfileID, operationID string) error {
	p.mu.Lock()
	if p.busy {
		p.mu.Unlock()
		return ErrOperationInProgress
	}
	p.busy = true
	previousID := p.appliedID
	previousJSON := append([]byte(nil), p.lastAppliedJSON...)
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		p.busy = false
		p.mu.Unlock()
	}()

	p.logger().Info("apply start", "component", "coordinator", "operation_id", operationID, "profile_id", string(profileID))
	if p.Health != nil {
		_ = p.Health.Recheck(ctx)
	}
	settings, err := p.Settings.Get(ctx)
	if err != nil {
		return fmt.Errorf("protection apply settings: %w", err)
	}
	disc, err := p.Locator.Locate(ctx, settings.ProxyBridgePath)
	if err != nil {
		return fmt.Errorf("protection apply locate: %w", err)
	}
	if disc.State == domain.BridgeStateNotInstalled {
		p.setBridge(domain.BridgeStateNotInstalled)
		return fmt.Errorf("protection apply: proxybridge not installed")
	}
	profile, err := p.Profiles.Get(ctx, profileID)
	if err != nil {
		return fmt.Errorf("protection apply profile: %w", err)
	}
	proxies, err := p.Proxies.List(ctx)
	if err != nil {
		return fmt.Errorf("protection apply proxies: %w", err)
	}
	identities, err := p.Processes.List(ctx)
	if err != nil {
		return fmt.Errorf("protection apply identities: %w", err)
	}
	compiled, err := p.Compiler.Compile(ctx, profile, proxies, identities)
	if err != nil {
		return fmt.Errorf("protection apply compile: %w", err)
	}
	if p.Process.Running() {
		_ = p.Process.Stop(ctx)
	}
	p.setBridge(domain.BridgeStateStarting)
	if err := p.Process.Start(ctx, disc.Path, compiled.JSON); err != nil {
		p.setBridge(domain.BridgeStateFailed)
		msg := err.Error()
		if strings.Contains(msg, string(port.BridgeLineAnotherCLI)) || strings.Contains(msg, string(port.BridgeLineGUIConflict)) {
			p.setBridge(domain.BridgeStateConflict)
		}
		return fmt.Errorf("protection apply start: %w", err)
	}
	if err := p.Process.WaitConfirm(ctx); err != nil {
		_ = p.Process.Stop(ctx)
		p.mu.Lock()
		p.appliedID = ""
		p.appliedConfirmed = false
		p.mu.Unlock()
		confirmMsg := err.Error()
		if strings.Contains(confirmMsg, string(port.BridgeLineAnotherCLI)) || strings.Contains(confirmMsg, string(port.BridgeLineGUIConflict)) {
			p.setBridge(domain.BridgeStateConflict)
		} else {
			p.setBridge(domain.BridgeStateStopped)
		}
		if len(previousJSON) > 0 && previousID != "" {
			p.logger().Info("apply rollback", "component", "coordinator", "operation_id", operationID, "profile_id", string(previousID))
			if startErr := p.Process.Start(ctx, disc.Path, previousJSON); startErr == nil {
				if waitErr := p.Process.WaitConfirm(ctx); waitErr == nil {
					p.mu.Lock()
					p.appliedID = previousID
					p.appliedConfirmed = true
					p.lastAppliedJSON = previousJSON
					p.mu.Unlock()
					p.setBridge(domain.BridgeStateReady)
					return fmt.Errorf("protection apply confirm failed after rollback ok: %w", err)
				}
			}
			return fmt.Errorf("protection apply confirm and rollback failed: %w", err)
		}
		return fmt.Errorf("protection apply confirm: %w", err)
	}
	p.mu.Lock()
	p.appliedID = profileID
	p.appliedConfirmed = true
	p.lastAppliedJSON = compiled.JSON
	p.desiredDiffers = false
	p.mu.Unlock()
	p.setBridge(domain.BridgeStateReady)
	settings.DesiredProfileID = profileID
	settings.DesiredSet = true
	_ = p.Settings.Save(ctx, settings)
	p.logger().Info("apply confirmed", "component", "coordinator", "operation_id", operationID, "profile_id", string(profileID))
	return nil
}

// Stop stops the managed CLI.
func (p *Protection) Stop(ctx context.Context, operationID string) error {
	p.mu.Lock()
	if p.busy {
		p.mu.Unlock()
		return ErrOperationInProgress
	}
	p.busy = true
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		p.busy = false
		p.mu.Unlock()
	}()
	p.logger().Info("stop start", "component", "coordinator", "operation_id", operationID)
	p.setBridge(domain.BridgeStateStopping)
	if err := p.Process.Stop(ctx); err != nil {
		p.setBridge(domain.BridgeStateFailed)
		return fmt.Errorf("protection stop: %w", err)
	}
	p.mu.Lock()
	p.appliedID = ""
	p.appliedConfirmed = false
	p.lastAppliedJSON = nil
	p.mu.Unlock()
	p.setBridge(domain.BridgeStateStopped)
	p.logger().Info("stop done", "component", "coordinator", "operation_id", operationID)
	return nil
}

// Applied returns applied bookkeeping for snapshots.
func (p *Protection) Applied() (domain.ProfileID, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.appliedID, p.appliedConfirmed
}

// BridgeState returns the last known bridge state.
func (p *Protection) BridgeState() domain.BridgeState {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.bridgeState == "" {
		return domain.BridgeStateUnknown
	}
	return p.bridgeState
}

// MarkDesiredDiffers records that desired profile content or selection differs from applied.
func (p *Protection) MarkDesiredDiffers() {
	p.mu.Lock()
	p.desiredDiffers = true
	p.mu.Unlock()
}

// DesiredDiffers reports whether desired differs from applied.
func (p *Protection) DesiredDiffers() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.desiredDiffers {
		return true
	}
	if p.appliedConfirmed && p.appliedID != "" {
		return false
	}
	return false
}

func (p *Protection) setBridge(state domain.BridgeState) {
	p.mu.Lock()
	p.bridgeState = state
	p.mu.Unlock()
	p.logger().Info("bridge state", "component", "coordinator", "bridge_state", string(state))
}
