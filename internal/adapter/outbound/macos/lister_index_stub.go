//go:build !darwin

package macos

import (
	"context"
	"fmt"
)

func queryApplicationIndex(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("macos application index: unsupported platform")
}
