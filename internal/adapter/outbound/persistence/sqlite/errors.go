// Package sqlite is the GORM SQLite persistence adapter for the catalog.
package sqlite

import "errors"

// ErrNotFound is returned when a row is missing.
var ErrNotFound = errors.New("not found")
