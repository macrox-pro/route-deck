// Package main is the Wails composition root for MACROX RouteDeck.
package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync/atomic"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/tray"
	"github.com/macrox-pro/route-deck/internal/adapter/inbound/wails"
	"github.com/macrox-pro/route-deck/internal/adapter/outbound/autostart"
	"github.com/macrox-pro/route-deck/internal/adapter/outbound/keychain"
	"github.com/macrox-pro/route-deck/internal/adapter/outbound/logfile"
	"github.com/macrox-pro/route-deck/internal/adapter/outbound/macos"
	"github.com/macrox-pro/route-deck/internal/adapter/outbound/persistence/sqlite"
	"github.com/macrox-pro/route-deck/internal/adapter/outbound/proxybridge"
	"github.com/macrox-pro/route-deck/internal/adapter/outbound/proxyhealth"
	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/application/service"
	"github.com/macrox-pro/route-deck/internal/domain"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "route-deck: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	level := new(slog.LevelVar)
	level.Set(slog.LevelInfo)
	logger, logCloser, err := newLogger(level)
	if err != nil {
		return fmt.Errorf("route-deck: %w", err)
	}
	if logCloser != nil {
		defer func() {
			if closeErr := logCloser.Close(); closeErr != nil {
				logger.Info("log file close failed", "component", "main", "err", closeErr)
			}
		}()
	}
	slog.SetDefault(logger)

	dbPath, err := sqlite.DefaultPath()
	if err != nil {
		return fmt.Errorf("route-deck: %w", err)
	}

	logPath, _ := logfile.DefaultPath()
	logger.Info("route-deck starting",
		"component", "main",
		"db_path", dbPath,
		"log_path", logPath,
	)

	locator := &proxybridge.Locator{Logger: logger}
	secrets := &keychain.Store{Logger: logger}
	compiler := &proxybridge.Compiler{Secrets: secrets, Logger: logger}
	bridgeProcess := &proxybridge.Process{Logger: logger}
	healthChecker := &proxyhealth.Checker{}

	snapshotter := &service.Snapshotter{Logger: logger, Locator: locator}
	var db *sqlite.DB
	var catalog *service.Catalog
	var protection *service.Protection
	var healthMon *service.HealthMonitor
	var settingsStore port.SettingsStore
	recheckPending := &atomic.Bool{}

	opened, openErr := sqlite.Open(ctx, dbPath)
	if openErr != nil {
		logger.Info("database open failed", "component", "main", "err", openErr)
		snapshotter.DatabaseError = openErr.Error()
		snapshotter.Proxies = unavailableProxyStore{}
	} else {
		db = opened
		store := db.Store()
		snapshotter.Proxies = store
		snapshotter.Settings = store.Settings()
		settingsStore = store.Settings()
		settings, settingsErr := settingsStore.Get(ctx)
		if settingsErr != nil {
			logger.Info("settings load failed", "component", "main", "err", settingsErr)
		} else if settings.DebugLogging {
			level.Set(slog.LevelDebug)
		}

		catalog = &service.Catalog{
			Proxies:   store,
			Apps:      store.Applications(),
			Processes: store.ProcessIdentities(),
			Profiles:  store.Profiles(),
			Inspector: macos.Inspector{},
			Installed: &macos.Lister{Logger: logger},
			AppIcons:  &macos.IconReader{},
		}
		healthMon = &service.HealthMonitor{
			Proxies: store,
			Secrets: secrets,
			Checker: healthChecker,
			Logger:  logger,
		}
		healthMon.Start(ctx)
		protection = &service.Protection{
			Profiles:  store.Profiles(),
			Proxies:   store,
			Processes: store.ProcessIdentities(),
			Settings:  settingsStore,
			Secrets:   secrets,
			Locator:   locator,
			Compiler:  compiler,
			Process:   bridgeProcess,
			Health:    healthMon,
			Logger:    logger,
		}
		defer func() {
			healthMon.Stop()
			if err := db.Close(); err != nil {
				logger.Info("database close failed", "component", "main", "err", err)
			}
		}()
	}

	live := &wails.LiveSnapshotter{
		Base:           snapshotter,
		Protection:     protection,
		Health:         healthMon,
		RecheckPending: recheckPending,
	}
	snapshotService := &wails.SnapshotService{Reader: live, Logger: logger}

	var settingsReader wails.SettingsStoreReader
	if settingsStore != nil {
		settingsReader = wails.PortSettingsReader{Store: settingsStore}
	}
	var emit wails.Emitter
	protectionService := &wails.ProtectionService{
		Coordinator: protection,
		Settings:    settingsReader,
		Emit:        nil, // set after app exists
		Logger:      logger,
	}
	proxiesService := &wails.ProxiesService{
		Catalog: catalog,
		Secrets: secrets,
		Health:  healthMon,
		Logger:  logger,
	}
	applicationsService := &wails.ApplicationsService{Catalog: catalog, Logger: logger}
	profilesService := &wails.ProfilesService{
		Catalog:    catalog,
		Settings:   settingsStore,
		Protection: protection,
		Logger:     logger,
	}
	appSettings := &service.Settings{
		Store:   settingsStore,
		Locator: locator,
		Level:   level,
		Logger:  logger,
	}
	settingsService := &wails.SettingsService{
		Settings: appSettings,
		Logger:   logger,
	}

	trayCtrl := &tray.Controller{
		Snapshot: tray.Snapshotter{Snapshot: live.Snapshot},
		Logger:   logger,
	}
	if protection != nil {
		trayCtrl.Bridge = protection
		trayCtrl.Protection = protection
	}

	application.RegisterEvent[wails.EventEnvelope](wails.EventProtectionStateChanged)
	application.RegisterEvent[wails.EventEnvelope](wails.EventProxyHealthChanged)
	application.RegisterEvent[wails.EventEnvelope](wails.EventProxyBridgeStateChanged)
	application.RegisterEvent[wails.EventEnvelope](wails.EventOperationProgress)
	application.RegisterEvent[wails.EventEnvelope](wails.EventOperationFailed)

	app := application.New(application.Options{
		Name:        "route-deck",
		Description: "MACROX RouteDeck assigns each selected application one policy and applies that policy by driving the official ProxyBridge CLI.",
		Logger:      logger,
		LogLevel:    level.Level(),
		Services: []application.Service{
			application.NewService(snapshotService),
			application.NewService(protectionService),
			application.NewService(proxiesService),
			application.NewService(applicationsService),
			application.NewService(profilesService),
			application.NewService(settingsService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
			ActivationPolicy: application.ActivationPolicyRegular,
		},
		ShouldQuit: trayCtrl.ShouldQuit,
		OnShutdown: func() {
			logger.Info("shutdown begin", "component", "main")
			if protection != nil {
				_ = protection.Stop(context.Background(), "shutdown")
			}
		},
	})

	appSettings.LaunchAtLogin = &autostart.Manager{Autostart: app.Autostart}

	emit = func(name string, envelope wails.EventEnvelope) {
		app.Event.Emit(name, envelope)
	}
	protectionService.Emit = emit
	proxiesService.Emit = emit

	app.Event.OnApplicationEvent(events.Common.SystemDidWake, func(_ *application.ApplicationEvent) {
		logger.Info("system did wake", "component", "main")
		live.MarkRecheckPending()
		emit(wails.EventProtectionStateChanged, wails.EventEnvelope{
			Schema:    1,
			Timestamp: "",
			Payload:   map[string]any{"reason": "wake"},
		})
		go func() {
			if healthMon != nil {
				_ = healthMon.Recheck(context.Background())
			}
			live.ClearRecheckPending()
			emit(wails.EventProxyHealthChanged, wails.EventEnvelope{
				Schema:  1,
				Payload: map[string]any{"reason": "wake"},
			})
			trayCtrl.OnProtectionEvent()
		}()
	})
	app.Event.OnApplicationEvent(events.Common.NetworkChanged, func(_ *application.ApplicationEvent) {
		logger.Info("network changed", "component", "main")
		live.MarkRecheckPending()
		go func() {
			if healthMon != nil {
				_ = healthMon.Recheck(context.Background())
			}
			live.ClearRecheckPending()
			emit(wails.EventProxyHealthChanged, wails.EventEnvelope{
				Schema:  1,
				Payload: map[string]any{"reason": "network"},
			})
			trayCtrl.OnProtectionEvent()
		}()
	})

	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "MACROX RouteDeck",
		Width:     1000,
		Height:    618,
		MinWidth:  800,
		MinHeight: 500,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
	})

	trayCtrl.App = app
	trayCtrl.Window = mainWindow
	if err := trayCtrl.Attach(); err != nil {
		return fmt.Errorf("route-deck: tray: %w", err)
	}

	refreshTray := func() { trayCtrl.OnProtectionEvent() }
	app.Event.On(wails.EventProtectionStateChanged, func(_ *application.CustomEvent) {
		refreshTray()
	})
	app.Event.On(wails.EventProxyBridgeStateChanged, func(_ *application.CustomEvent) {
		refreshTray()
	})

	if err := app.Run(); err != nil {
		return fmt.Errorf("route-deck: %w", err)
	}
	return nil
}

func newLogger(level *slog.LevelVar) (*slog.Logger, io.Closer, error) {
	opts := &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: logfile.RedactAttr,
	}
	stderrHandler := slog.NewTextHandler(os.Stderr, opts)

	logPath, err := logfile.DefaultPath()
	if err != nil {
		fallback := slog.New(stderrHandler)
		fallback.Info("log path unavailable; stderr only", "component", "main", "err", err)
		return fallback, nil, nil
	}

	fileWriter, err := logfile.Open(logPath, logfile.Options{})
	if err != nil {
		fallback := slog.New(stderrHandler)
		fallback.Info("log file unavailable; stderr only", "component", "main", "err", err)
		return fallback, nil, nil
	}

	fileHandler := slog.NewJSONHandler(fileWriter, opts)
	logger := slog.New(slog.NewMultiHandler(fileHandler, stderrHandler))
	return logger, fileWriter, nil
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

func (unavailableProxyStore) Delete(ctx context.Context, id domain.ProxyID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_ = id
	return fmt.Errorf("database unavailable")
}
