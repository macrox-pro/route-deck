//go:build !darwin

package macos

import (
	"context"
	"fmt"

	"github.com/macrox-pro/route-deck/internal/application/port"
)

// IconReader is a non-darwin stub.
type IconReader struct{}

var _ port.ApplicationIconReader = (*IconReader)(nil)

func (r *IconReader) PNG(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	_ = path
	return nil, fmt.Errorf("application icon: unsupported platform")
}
