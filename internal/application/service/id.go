package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func newID() (string, error) {
	return NewID()
}

// NewID returns a 128-bit hex id for entities and operation ids.
func NewID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("new id: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}
