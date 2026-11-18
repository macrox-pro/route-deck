# Handoff

Timestamp: 2026-11-18 21:25 UTC
Agent: agent (RD-059 Appearance + native frame sync)
Task ID: RD-059
Branch / commit: working tree (agent MUST NOT commit)

## Done

- SQLite `appearance` on settings row; port + AutoMigrate + store tests
- `service.SetAppearance` / NormalizeAppearance (system|light|dark); unit tests
- Wails `SettingsDTO.appearance` + `SetAppearance(appearance, mainHash)`
- `WindowChrome`: dirty Discard gate, persist-first, recreate main, rewire Parent/tray/NativeUserDialog
- Main + editor windows: verified `Mac.Appearance` / `Windows.Theme`; no white BackgroundColour
- Frontend Select + ThemeProvider SoT via Settings; localStorage FOUC write-through only
- Bindings regenerated; wails-api + desktop-ux updated

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/application/service/ ./internal/adapter/outbound/persistence/sqlite/ ./internal/adapter/inbound/wails/` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `cd frontend && npm run build` | PASS |
| `wails3 generate bindings -ts` | PASS |
| macOS Appearance recreate smoke | NOT RUN (owner) |

## Exact next step

Owner: `wails3 dev` → Settings → Appearance; confirm title bar + content; with a dirty editor Cancel must abort without save.

## Suggested commit message

```text
persist appearance in sqlite and recreate windows for native chrome (rd-059)
```

## Context to read

- `docs/project/tasks/RD-059.md`
- `internal/adapter/inbound/wails/window_chrome.go`
- `internal/adapter/inbound/wails/settings.go`
- `frontend/src/providers/theme-provider.tsx`
- `frontend/src/pages/settings.tsx`

## Notes

- Calling webview may die during recreate; DB already saved.
- Editors are closed on successful change and not reopened.
- Prior localStorage-only theme is not auto-imported into SQLite.
