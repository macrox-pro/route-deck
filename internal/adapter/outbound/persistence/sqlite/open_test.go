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
