//go:build keychain_integration && darwin

package keychain_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/keychain"
	"github.com/macrox-pro/route-deck/internal/domain"
)

func TestKeychainRoundTrip(t *testing.T) {
	store := &keychain.Store{}
	ctx := context.Background()
	ref := domain.KeychainRef("rd-test-proxy-integration")
	secret := []byte("test-password")

	t.Cleanup(func() {
		_ = store.Delete(ctx, ref)
	})

	if err := store.Save(ctx, ref, secret); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load(ctx, ref)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !bytes.Equal(got, secret) {
		t.Fatalf("Load = %q, want %q", got, secret)
	}
	if err := store.Delete(ctx, ref); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
