//go:build darwin

// Package keychain implements SecretStore with macOS Keychain via go-keychain.
package keychain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/keybase/go-keychain"
	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ServiceName is the Keychain service for proxy passwords (ADR-0008).
const ServiceName = "pro.macrox.routedeck.proxy"

// Store implements port.SecretStore.
type Store struct {
	Logger *slog.Logger
}

var _ port.SecretStore = (*Store)(nil)

func (s *Store) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// Load returns secret bytes for account = string(ref).
func (s *Store) Load(ctx context.Context, ref domain.KeychainRef) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("keychain load: %w", err)
	}
	account := string(ref)
	if account == "" {
		return nil, fmt.Errorf("keychain load: empty ref")
	}
	data, err := keychain.GetGenericPassword(ServiceName, account, "", "")
	if err != nil {
		s.logger().Error("keychain load failed", "component", "keychain", "proxy_id", account, "err", err)
		return nil, fmt.Errorf("keychain load: %w", err)
	}
	if data == nil {
		return nil, fmt.Errorf("keychain load: not found")
	}
	s.logger().Debug("keychain load ok", "component", "keychain", "proxy_id", account)
	return data, nil
}

// Save stores secret bytes for account = string(ref). Overwrites if present.
func (s *Store) Save(ctx context.Context, ref domain.KeychainRef, secret []byte) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("keychain save: %w", err)
	}
	account := string(ref)
	if account == "" {
		return fmt.Errorf("keychain save: empty ref")
	}
	item := keychain.NewGenericPassword(ServiceName, account, "", secret, "")
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)
	err := keychain.AddItem(item)
	if errors.Is(err, keychain.ErrorDuplicateItem) {
		query := keychain.NewItem()
		query.SetSecClass(keychain.SecClassGenericPassword)
		query.SetService(ServiceName)
		query.SetAccount(account)
		update := keychain.NewItem()
		update.SetData(secret)
		err = keychain.UpdateItem(query, update)
	}
	if err != nil {
		s.logger().Error("keychain save failed", "component", "keychain", "proxy_id", account, "err", err)
		return fmt.Errorf("keychain save: %w", err)
	}
	s.logger().Info("keychain save ok", "component", "keychain", "proxy_id", account)
	return nil
}

// Delete removes the Keychain item for account = string(ref).
func (s *Store) Delete(ctx context.Context, ref domain.KeychainRef) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("keychain delete: %w", err)
	}
	account := string(ref)
	if account == "" {
		return fmt.Errorf("keychain delete: empty ref")
	}
	item := keychain.NewGenericPassword(ServiceName, account, "", nil, "")
	err := keychain.DeleteItem(item)
	if err != nil && !errors.Is(err, keychain.ErrorItemNotFound) {
		s.logger().Error("keychain delete failed", "component", "keychain", "proxy_id", account, "err", err)
		return fmt.Errorf("keychain delete: %w", err)
	}
	s.logger().Info("keychain delete ok", "component", "keychain", "proxy_id", account)
	return nil
}
