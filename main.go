// Package main is the Wails composition root for MACROX RouteDeck.
package main

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/wails"
	"github.com/macrox-pro/route-deck/internal/adapter/outbound/persistence/sqlite"
	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	dbPath, err := sqlite.DefaultPath()
	if err != nil {
		return fmt.Errorf("route-deck: %w", err)
	}

	snapshotter := &service.Snapshotter{}
	var db *sqlite.DB
	opened, openErr := sqlite.Open(ctx, dbPath)
	if openErr != nil {
		log.Printf("route-deck: database open failed: %v", openErr)
		snapshotter.DatabaseError = openErr.Error()
		snapshotter.Proxies = unavailableProxyStore{}
	} else {
		db = opened
		snapshotter.Proxies = db.Store()
		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("route-deck: database close: %v", err)
			}
		}()
	}

	snapshotService := &wails.SnapshotService{Reader: snapshotter}

	app := application.New(application.Options{
		Name:        "route-deck",
		Description: "MACROX RouteDeck assigns each selected application one policy and applies that policy by driving the official ProxyBridge CLI.",
		Services: []application.Service{
			application.NewService(snapshotService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "MACROX RouteDeck",
		Width:     1000,
		Height:    618,
		MinWidth:  800,
		MinHeight: 500,
		Mac: application.MacWindow{
			// Default title bar: traffic lights stay in the system chrome.
			// HiddenInset caused a white toolbar strip on fullscreen hover.
			TitleBar: application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		return fmt.Errorf("route-deck: %w", err)
	}
	return nil
}

// unavailableProxyStore is used when SQLite did not open.
type unavailableProxyStore struct{}

var _ port.ProxyStore = unavailableProxyStore{}

func (unavailableProxyStore) Get(ctx context.Context, _ domain.ProxyID) (domain.ProxyEndpoint, error) {
	if err := ctx.Err(); err != nil {
		return domain.ProxyEndpoint{}, err
	}
	return domain.ProxyEndpoint{}, fmt.Errorf("database unavailable")
}

func (unavailableProxyStore) List(ctx context.Context) ([]domain.ProxyEndpoint, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (unavailableProxyStore) Save(ctx context.Context, endpoint domain.ProxyEndpoint) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_ = endpoint
	return fmt.Errorf("database unavailable")
}
