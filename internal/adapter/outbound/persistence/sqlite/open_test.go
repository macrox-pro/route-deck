package sqlite

import (
	"context"
	"path/filepath"
	"testing"
)

func TestOpenPragmasAndPool(t *testing.T) {
	path := filepath.Join(t.TempDir(), "route-deck.db")
	db, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})

	sqlDB, err := db.gormDB.DB()
	if err != nil {
		t.Fatalf("DB(): %v", err)
	}
	if sqlDB.Stats().MaxOpenConnections != 1 {
		t.Fatalf("MaxOpenConnections = %d, want 1", sqlDB.Stats().MaxOpenConnections)
	}

	var journal string
	if err := db.gormDB.Raw("PRAGMA journal_mode").Scan(&journal).Error; err != nil {
		t.Fatalf("journal_mode: %v", err)
	}
	if journal != "wal" {
		t.Fatalf("journal_mode = %q, want wal", journal)
	}

	var foreignKeys int
	if err := db.gormDB.Raw("PRAGMA foreign_keys").Scan(&foreignKeys).Error; err != nil {
		t.Fatalf("foreign_keys: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("foreign_keys = %d, want 1", foreignKeys)
	}
}

func TestNoPasswordColumn(t *testing.T) {
	path := filepath.Join(t.TempDir(), "route-deck.db")
	db, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})

	type column struct {
		Name string `gorm:"column:name"`
	}
	var cols []column
	if err := db.gormDB.Raw("PRAGMA table_info(proxies)").Scan(&cols).Error; err != nil {
		t.Fatalf("table_info: %v", err)
	}
	for _, col := range cols {
		if col.Name == "password" {
			t.Fatal("proxies must not have a password column")
		}
	}
}

func TestAutoMigrateAddsAppearanceWithDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "route-deck.db")
	db, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := db.gormDB.Exec(`CREATE TABLE IF NOT EXISTS settings_legacy (
		id INTEGER PRIMARY KEY,
		desired_profile_id TEXT NOT NULL,
		desired_set INTEGER NOT NULL,
		proxybridge_path TEXT NOT NULL,
		debug_logging INTEGER NOT NULL
	)`).Error; err != nil {
		_ = db.Close()
		t.Fatalf("create legacy: %v", err)
	}
	// Replace settings with a pre-appearance schema to simulate an upgraded install.
	if err := db.gormDB.Exec(`DROP TABLE settings`).Error; err != nil {
		_ = db.Close()
		t.Fatalf("drop settings: %v", err)
	}
	if err := db.gormDB.Exec(`ALTER TABLE settings_legacy RENAME TO settings`).Error; err != nil {
		_ = db.Close()
		t.Fatalf("rename: %v", err)
	}
	if err := db.gormDB.Exec(
		`INSERT INTO settings (id, desired_profile_id, desired_set, proxybridge_path, debug_logging) VALUES (1, '', 0, '', 0)`,
	).Error; err != nil {
		_ = db.Close()
		t.Fatalf("seed: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	db2, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open after legacy schema: %v", err)
	}
	t.Cleanup(func() {
		_ = db2.Close()
	})

	settings, err := db2.Store().Settings().Get(context.Background())
	if err != nil {
		t.Fatalf("Get settings: %v", err)
	}
	if settings.Appearance != "system" {
		t.Fatalf("appearance = %q, want system", settings.Appearance)
	}
	if !settings.AutoApply {
		t.Fatal("auto_apply = false after migrate, want true")
	}
}

func TestAutoMigrateAddsAutoApplyWithDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "route-deck.db")
	db, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := db.gormDB.Exec(`DROP TABLE settings`).Error; err != nil {
		_ = db.Close()
		t.Fatalf("drop settings: %v", err)
	}
	if err := db.gormDB.Exec(`CREATE TABLE settings (
		id INTEGER PRIMARY KEY,
		desired_profile_id TEXT NOT NULL,
		desired_set INTEGER NOT NULL,
		proxybridge_path TEXT NOT NULL,
		debug_logging INTEGER NOT NULL,
		appearance TEXT NOT NULL DEFAULT 'system'
	)`).Error; err != nil {
		_ = db.Close()
		t.Fatalf("create pre-auto_apply schema: %v", err)
	}
	if err := db.gormDB.Exec(
		`INSERT INTO settings (id, desired_profile_id, desired_set, proxybridge_path, debug_logging, appearance) VALUES (1, '', 0, '', 0, 'system')`,
	).Error; err != nil {
		_ = db.Close()
		t.Fatalf("seed: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	db2, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("Open after pre-auto_apply schema: %v", err)
	}
	t.Cleanup(func() {
		_ = db2.Close()
	})

	settings, err := db2.Store().Settings().Get(context.Background())
	if err != nil {
		t.Fatalf("Get settings: %v", err)
	}
	if !settings.AutoApply {
		t.Fatal("auto_apply = false after AutoMigrate, want true")
	}
}
