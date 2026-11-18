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
| `ApplicationsService` | `List`, `ListInstalled`, `Icon(path)`, `AddFromPath(path)`, `ListProcessIdentities`, `ListHelperGroups(applicationID)`, `PickAndAdd`, `InspectHelpers(applicationID)`, `ConfirmHelpers`, `Delete` |
| `ProfilesService` | `List`, `Save(ProfileSaveDTO)`, `Copy`, `Rename`, `Delete`, `SetDesired` |
| `EditorDialogsService` | `OpenProfile(id)`, `OpenProxy(id, guide)`, `OpenHelpers(applicationID)`, `SetDirty(kind, dirty)`, `Close(kind)` — ADR-0014 custom dialog windows |
| `SettingsService` | `Get`, `SetProxyBridgePath`, `SetDebugLogging`, `SetLaunchAtLogin`, `SetAppearance`, `SetAutoApply`, `OpenLogsFolder` |
| `ActivityService` | `Recent(before)` — page of up to 100 JSON log records; empty `before` is newest; exclusive RFC3339 cursor loads older; returns `hasMore` and `attrs` for non-secret string fields (ADR-0011) |
| `SetupGuidesService` | `List(refresh)`, `OpenDocs(guideID)`, `AddProxy(guideID)`, `Suggestions`, `DismissSuggestion(guideID, port)`, `Checklist` — Setup guides detection, one-click proxy add, suggestion banner, first-run checklist (ADR-0013). `List(true)` refreshes the installed-app cache. Bindings: `wails3 generate bindings -ts` |

`Apply` returns an `operation_id` immediately and runs confirm asynchronously. Progress and failure use product events. The client must reload the snapshot after events; events are not proof of `Active`. Snapshot `appliedConfirmed` means the CLI confirmed a profile; it is not the same as Protection `Active`. Snapshot proxy rows include `kind`, `host`, and `port` for UI matching (icons/labels); health wire tokens are unchanged.

`ApplicationsService.InspectHelpers` re-inspects a catalog application path and returns helper candidates plus `pendingCount`. It does not save them; callers must use `ConfirmHelpers` (ADR-0014). `ConfirmHelpers` is a sync: checked groups are upserted (one identity per bundle id; path refresh on re-confirm); unchecked groups delete every identity with that bundle id, or return `ErrReferenced` if any profile references any of them. `ListHelperGroups` returns Go-built groups for Profiles and the helpers dialog. Product event `catalog.changed` fires after catalog mutations so other windows refetch; events are not the source of truth.

`EditorDialogsService` opens reused native-frame windows (`profile`, `proxy`, `helpers`) with `AttachModal` + parent `SetEnabled(false)`. Bindings return immediately. Dirty close uses `Dialog.Question` ("Discard changes?"). See ADR-0014 and `docs/engineering/desktop-ux.md`.

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

`SettingsDTO.autoApply` is SQLite (default `true`). `SetAutoApply` persists the flag only. When true, process start evaluates ADR-0015 gates and may call `ProtectionService.Apply` once (independent of Launch at Login). Skip or failure does not rewrite policies to Direct.

`SettingsDTO.appearance` is `system` | `light` | `dark` (default `system`) stored in SQLite. `SetAppearance(appearance, mainHash)` persists first, then recreates the main window with verified create-time `Mac.Appearance` / `Windows.Theme` (Wails v3.0.0 has no public runtime SetAppearance). Dirty editor dialogs prompt Discard; Cancel aborts without save. Editors are closed and not reopened. `mainHash` is the caller's `location.hash` so the new main window can reopen on the same route. `localStorage` is only a FOUC write-through cache of the SQLite value.

`SettingsDTO.appVersion` / `appVersionLabel` / `appCommit` come from link-time `internal/buildmeta` (not SQLite). Untagged builds use `appVersion` `dev` and label `dev`. Tagged CI builds strip the leading `v` from `vMAJOR.MINOR.PATCH` into `appVersion` and set `appVersionLabel` to `vMAJOR.MINOR.PATCH`. `appCommit` is the short SHA when injected. The sidebar badge shows only `appVersionLabel`; Settings → About shows label and commit.

## Errors

Return typed errors the UI can map. MUST NOT return a bare string that means “ok” when the CLI did not confirm.
