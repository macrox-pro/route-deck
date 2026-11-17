# Wails API contract

Library: `github.com/wailsapp/wails/v3` v3.0.0. Docs: https://v3.wails.io/

MUST NOT invent methods. Read the installed module or the v3 docs before adding a binding.

## Current product surface

`main.go` registers these inbound services (RD-029 / RD-030 / RD-037). Tray lifecycle is ADR-0010 (`ApplicationShouldTerminateAfterLastWindowClosed: false`).

| Service | Methods |
| --- | --- |
| `SnapshotService` | `Get(ctx) (SnapshotDTO, error)` |
| `ProtectionService` | `Apply(ctx) (ApplyResultDTO, error)`, `Stop(ctx) (ApplyResultDTO, error)` |
| `ProxiesService` | `List`, `Save(ProxySaveDTO)`, `Delete(id)`, `ClearPassword(id)`, `Recheck` |
| `ApplicationsService` | `List`, `ListInstalled`, `Icon(path)`, `AddFromPath(path)`, `ListProcessIdentities`, `PickAndAdd`, `ConfirmHelpers`, `Delete` |
| `ProfilesService` | `List`, `Save(ProfileSaveDTO)`, `Copy`, `Rename`, `Delete`, `SetDesired` |
| `SettingsService` | `Get`, `SetProxyBridgePath`, `SetDebugLogging`, `SetLaunchAtLogin`, `OpenLogsFolder` |
| `ActivityService` | `Recent(before)` — page of up to 100 JSON log records; empty `before` is newest; exclusive RFC3339 cursor loads older; returns `hasMore` and `attrs` for non-secret string fields (ADR-0011) |
| `SetupGuidesService` | `List(refresh)`, `OpenDocs(guideID)`, `AddProxy(guideID)`, `Suggestions`, `DismissSuggestion(guideID, port)`, `Checklist` — Setup guides detection, one-click proxy add, suggestion banner, first-run checklist (ADR-0013). `List(true)` refreshes the installed-app cache. Bindings: `wails3 generate bindings -ts` |

`Apply` returns an `operation_id` immediately and runs confirm asynchronously. Progress and failure use product events. The client must reload the snapshot after events; events are not proof of `Active`.

Native dialogs (not Wails service methods): GUI conflict on Apply failure uses `Dialog.Warning`; first-run extension approval during confirm uses `Dialog.Info` once per CLI start. Copy lives in `internal/adapter/inbound/wails/usermsg.go` with stable IDs for future i18n.

Empty `password` on `ProxiesService.Save` for an existing id keeps the Keychain secret. Clearing a password is a separate `ClearPassword` call.

`ApplicationsService.PickAndAdd` uses Wails `Dialog.OpenFile()` (verified in v3.0.0). `SettingsService.OpenLogsFolder` uses `Env.OpenFileManager` (verified in v3.0.0). `SettingsService.SetLaunchAtLogin` uses Wails `app.Autostart` (verified in v3.0.0). OS registration is the source of truth; the flag is not stored in SQLite.

## Tray (not a Wails service)

`internal/adapter/inbound/tray` attaches a `SystemTray` from the composition root. Tray click opens the menu. Menu: Status (disabled), Open RouteDeck, Stop protection, Quit RouteDeck. No Apply from tray. Quit always uses `Dialog.Question`. See ADR-0010.

## Pattern

Wails services are inbound adapters:

1. Accept a DTO
2. Validate input
3. Call an application service
4. Return a DTO and a typed error

They MUST NOT open GORM and MUST NOT call `exec.Command`.

## Types

Generated bindings under `frontend/bindings/.../inbound/wails/` are the TypeScript source for these DTOs. Hand-written duplicates are forbidden. Domain types and GORM models are not the DTO.

`InstalledApplicationDTO` includes `isSystem` from the platform lister (ADR-0009). The UI filters with that flag; it MUST NOT hard-code OS paths.

`SettingsDTO.discoveredProxyBridgePath` is filled by the application Settings use case via `BridgeLocator.Locate("")` (platform `PlatformPaths` only). The UI uses it as the ProxyBridge path placeholder and MUST NOT hard-code install locations.

`SettingsDTO.launchAtLogin` comes from `Autostart.IsEnabled`.

## Errors

Return typed errors the UI can map. MUST NOT return a bare string that means “ok” when the CLI did not confirm.
