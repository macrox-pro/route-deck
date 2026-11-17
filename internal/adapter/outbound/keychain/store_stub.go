//go:build !darwin

package keychain

import (
	"context"
	"fmt"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// ServiceName matches the darwin adapter (ADR-0008).
const ServiceName = "pro.macrox.routedeck.proxy"

// Store is a non-darwin stub.
type Store struct {
	Logger any
}

var _ port.SecretStore = (*Store)(nil)

func (s *Store) Load(ctx context.Context, ref domain.KeychainRef) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	_ = ref
	return nil, fmt.Errorf("keychain: unsupported platform")
}

func (s *Store) Save(ctx context.Context, ref domain.KeychainRef, secret []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_ = ref
	_ = secret
	return fmt.Errorf("keychain: unsupported platform")
}

func (s *Store) Delete(ctx context.Context, ref domain.KeychainRef) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_ = ref
	return fmt.Errorf("keychain: unsupported platform")
}
