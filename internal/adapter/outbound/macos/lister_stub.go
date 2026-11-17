//go:build !darwin

package macos

import (
	"context"
	"fmt"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

// UnsupportedLister fails closed on platforms without a real adapter.
type UnsupportedLister struct{}

var _ port.InstalledApplicationLister = UnsupportedLister{}

func (UnsupportedLister) List(ctx context.Context) ([]port.InstalledApplication, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("installed application lister: unsupported platform")
}
