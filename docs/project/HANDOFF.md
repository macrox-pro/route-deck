# Handoff

Timestamp: 2026-11-17 16:05 UTC
Agent: agent (Phase 7 tray + Launch at Login)
Task ID: RD-037 (current; RD-035 and RD-036 also completed this session)
Branch / commit: working tree (agent MUST NOT commit)

## Done

- ADR-0010 Accepted: tray menu, Quit dialog always, Dock stays, monochrome template icon, Launch at Login
- RD-035: `ShouldTerminate=false`, SystemTray, close→Hide, Open, Quit Question dialog
- RD-036: live Status/tooltip from snapshot + bridge; Stop enabled only when CLI owned; Quit copy mentions Stop when owned
- RD-037: `port.LaunchAtLogin` + Wails Autostart adapter; Settings toggle; bindings regenerated
- Docs: lifecycle, wails-api, ROADMAP Phase 7, ADR-0008 pointer, BACKLOG tasks

## Files changed (high level)

- `internal/adapter/inbound/tray/*` (controller, status helpers + tests, template icon)
- `internal/adapter/outbound/autostart/wails.go`
- `internal/application/port/launch.go`
- `internal/application/service/settings.go`
- `internal/adapter/inbound/wails/settings.go`
- `main.go`
- `frontend/src/pages/settings.tsx` + bindings
- `docs/decisions/ADR-0010-tray-autostart.md`, ADR-0008 pointer
- `docs/architecture/lifecycle.md`, `docs/contracts/wails-api.md`
- `docs/project/tasks/RD-035.md` .. `RD-037.md`, STATUS, BACKLOG, HANDOFF, ROADMAP

## Decisions

- See ADR-0010 closed table (no Apply from tray; Quit always via dialog; OS is Launch at Login source of truth).

## Verification

| Command | Result |
| --- | --- |
| `go test ./...` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `npm run build` (cwd `frontend/`) | PASS |
| `wails3 generate bindings` | PASS |
| `wails3 task run` | NOT RUN (owner: close/open/quit dialog / Launch at Login) |
| RD-033 live smoke | NOT RUN (owner SOCKS5) |

## Not done

- Owner interactive tray + autostart smoke
- RD-033 interactive scenarios
- Agent did not commit

## Exact next step

Owner: `wails3 task run` → close window (app stays; Dock + menu bar) → tray Open / Quit Cancel / Quit Confirm → Settings Launch at Login toggle.

## Suggested commit message

```text
add macOS tray lifecycle with quit dialog and launch-at-login settings (rd-035..rd-037)
```

## Context to read

- `AGENTS.md`
- `docs/decisions/ADR-0010-tray-autostart.md`
- `docs/architecture/lifecycle.md`
- `docs/project/tasks/RD-037.md`
