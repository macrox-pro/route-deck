# Project Status

Updated: 2026-11-16 21:00 UTC
Updated by: agent (RD-005 Phase 0 exit)

## Current phase

Phase 0 Contracts — exit complete. Phase 1 Toolchain already DONE. Phase 2 Domain ports not started.

## Current task

RD-005

## State

DONE

## Completed

- Docs-prep files created in this working tree. Not committed. Owner commits.
- RD-000: owner accepted `CONCEPTION.md` and ADR-0001 as written.
- RD-001: toolchain results recorded (see prior Verification).
- RD-002 identity work DONE (module, config, icons, `bin/RouteDeck.app`, lint PASS).
- RD-003: owner accepted ADR-0001..0006 as written via Phase 0 exit plan (2026-11-16). No ADR body edits.
- RD-004: match rule revised — `processNames` matches bundle id or executable name; dashboard aggregate unchanged.
- RD-005: ProxyBridge contracts frozen against live docs; installed CLI `ProxyBridge CLI 4.1.3`.

## In progress

- None (RD-005 DONE; Phase 0 exit)

## Blockers

- No Phase 2 READY task yet. Owner must add a domain-ports task before any `internal/` packages.
- None for RD-005.

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

## Next actions

1. Owner adds a Phase 2 domain-ports READY task in `docs/project/BACKLOG.md` before any `internal/` packages.
2. Optional later: RD-006 (GORM smoke) when persistence starts; RD-007 (CI baseline) when the owner asks.

## Verification

| Command | Result |
| --- | --- |
| Four ProxyBridge doc URLs vs contracts | PASS — source date 2026-11-16 |
| `/usr/local/bin/proxybridge --version` | PASS — `ProxyBridge CLI 4.1.3` |
| `rg` for `not the runtime match key` / `Not Bundle ID` | PASS — gone from product docs (only historical if any) |
| `go test ./...` | PASS — no test files |
| `golangci-lint run ./...` | PASS — 0 issues. Cache write warnings only (sandbox) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go declared `1.26.0`; Wails v3.0.0; ProxyBridge CLI 4.1.3 at documented paths
