package sqlite

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps a GORM connection. *gorm.DB must not leave this package.
type DB struct {
	gormDB *gorm.DB
}

// DefaultPath returns Application Support/route-deck/route-deck.db on macOS.
func DefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("sqlite default path: %w", err)
	}
	return filepath.Join(configDir, "route-deck", "route-deck.db"), nil
}

// Open opens path, sets SQLite pragmas, and AutoMigrates the catalog models.
func Open(ctx context.Context, path string) (*DB, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("sqlite open mkdir: %w", err)
	}

	gormDB, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("sqlite open sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)

	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000",
		"PRAGMA synchronous=NORMAL",
	}
	for _, pragma := range pragmas {
		if err := gormDB.WithContext(ctx).Exec(pragma).Error; err != nil {
			_ = sqlDB.Close()
			return nil, fmt.Errorf("sqlite open %s: %w", pragma, err)
		}
	}

	// AutoMigrate known catalog models only:
	// applicationModel, processIdentityModel, proxyModel, profileModel, applicationPolicyModel.
	if err := gormDB.WithContext(ctx).AutoMigrate(
		&applicationModel{},
		&processIdentityModel{},
		&proxyModel{},
		&profileModel{},
		&applicationPolicyModel{},
	); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("sqlite automigrate: %w", err)
	}

	return &DB{gormDB: gormDB}, nil
}

// Close closes the underlying SQL connection.
func (d *DB) Close() error {
	if d == nil || d.gormDB == nil {
		return nil
	}
	sqlDB, err := d.gormDB.DB()
	if err != nil {
		return fmt.Errorf("sqlite close: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("sqlite close: %w", err)
	}
	return nil
}

// Store returns a Store bound to this database.
func (d *DB) Store() *Store {
	return &Store{db: d.gormDB}
}
