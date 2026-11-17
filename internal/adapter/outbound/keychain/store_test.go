package keychain_test

import (
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/keychain"
	"github.com/macrox-pro/route-deck/internal/application/port"
)

func TestStoreImplementsSecretStore(t *testing.T) {
	var _ port.SecretStore = (*keychain.Store)(nil)
	if keychain.ServiceName != "pro.macrox.routedeck.proxy" {
		t.Fatalf("ServiceName = %q", keychain.ServiceName)
	}
}
