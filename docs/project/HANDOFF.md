# Handoff

Timestamp: 2026-11-18 22:10 UTC
Agent: agent (RD-060 Auto Apply at startup)
Task ID: RD-060
Branch / commit: working tree (agent MUST NOT commit)

## Done

- ADR-0015 Accepted; ADR-0008 / ADR-0010 / lifecycle supersede patches
- SQLite `auto_apply` default true + AutoMigrate tests
- `SetAutoApply` + SettingsDTO + bindings
- Gate helper `ShouldAutoApplyOnStartup` + `MaybeAutoApplyAtStartup` one-shot in main
- Settings UI toggle Auto Apply on launch (independent of Launch at Login)
- wails-api.md updated

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/application/service/ ./internal/adapter/outbound/persistence/sqlite/ ./internal/adapter/inbound/wails/` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `cd frontend && npm run build` | PASS |
| `wails3 generate bindings -ts` | PASS |
| macOS Auto Apply smoke | NOT RUN (owner) |

## Exact next step

Owner: `wails3 dev` → Settings → Auto Apply on launch ON with desired Proxy/Block profile → quit/relaunch and confirm Apply; OFF → no Apply; Direct-only desired → quiet skip (check logs).

## Suggested commit message

```text
add opt-out auto apply on launch with adr-0015 gates (rd-060)
```

## Context to read

- `docs/project/tasks/RD-060.md`
- `docs/decisions/ADR-0015-auto-apply-startup.md`
- `internal/application/service/auto_apply.go`
- `internal/adapter/inbound/wails/auto_apply.go`
- `frontend/src/pages/settings.tsx`

## Notes

- Launch at Login unchanged (OS Autostart only).
- Auto Apply default ON; never silent Direct.
- Login race: proxy not up → honest Apply failure (no retry in this task).
