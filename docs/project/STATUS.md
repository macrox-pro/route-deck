# Project Status

Updated: 2026-11-16 21:15 UTC
Updated by: agent (Phase 2 backlog write)

## Current phase

Phase 2 Domain ports — backlog text started. Code under `internal/` not started.

## Current task

RD-008

## State

READY

## Completed

- Docs-prep files created in this working tree. Not committed. Owner commits.
- RD-000: owner accepted `CONCEPTION.md` and ADR-0001 as written.
- RD-001: toolchain results recorded (see prior Verification).
- RD-002 identity work DONE (module, config, icons, `bin/RouteDeck.app`, lint PASS).
- RD-003: owner accepted ADR-0001..0006 as written via Phase 0 exit plan (2026-11-16). No ADR body edits.
- RD-004: match rule revised — `processNames` matches bundle id or executable name; dashboard aggregate unchanged.
- RD-005: ProxyBridge contracts frozen against live docs; installed CLI `ProxyBridge CLI 4.1.3`.
- Phase 2 backlog accepted: RD-008 READY; RD-009 through RD-012 TODO; task bodies live in `docs/project/tasks/RD-*.md`; BACKLOG is an index only.

## In progress

- BACKLOG split to per-task files (index + `docs/project/tasks/`). No `internal/` packages yet.

## Blockers

- None for RD-008. RD-009 through RD-012 wait on prior DONE. CONCEPTION is not edited until RD-011.

## Known limitations

- No domain, persistence, or ProxyBridge adapter code yet
- Wails template still includes `GreetService` and a demo time event
- ESLint is not part of the frontend scaffold and is not required
- `vendor/` committed per ADR-0006; refresh after dependency changes
- First `go test ./...` before a frontend build FAILs on missing `frontend/dist` embed
- Darwin icon is PNG→`.icns` only (no Icon Composer `Assets.car`)
- `build/` is excluded from golangci-lint (`^build/` in `.golangci.yml`). It is the Wails scaffold
- After `wails3 task common:update:build-assets`, re-apply `LSMinimumSystemVersion` `13.0.0` and drop `CFBundleIconName` when not shipping `Assets.car`
- Linker may warn about macOS 12.0 link minimum; plist floor is 13.0.0
- Homebrew and other non-documented ProxyBridge install paths stay “not documented; do not guess”
- Phase 2 exit waits until RD-011 is DONE (owner choice)

## Next actions

1. Open `docs/project/tasks/RD-008.md` and start the RD-008 loop at `internal/domain/errors.go`, then `internal/domain/id.go`.
2. After RD-008 DONE: promote RD-009 in its task file + BACKLOG index. Optional later: RD-006 / RD-007.

## Verification

| Command | Result |
| --- | --- |
| `golangci-lint run ./...` | PASS — 0 issues. Cache write warnings only (sandbox) |
| `go test ./...` | PASS — no test files in product packages; exit 0 |
| Closed-list review of RD-008..012 | PASS — only RD-008 READY; RD-011 TODO not BLOCKED; ports only Get/Load; Username/Password/List/Save/Delete only in MUST NOT; CONCEPTION untouched; no `internal/` |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go declared `1.26.0`; Wails v3.0.0; ProxyBridge CLI 4.1.3 at documented paths
