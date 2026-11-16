# Project Status

Updated: 2026-11-16 20:44 UTC
Updated by: agent (RD-002 lint PASS)

## Current phase

Phase 1 Toolchain

## Current task

RD-002

## State

DONE

## Completed

- Docs-prep files created in this working tree. Not committed. Owner commits.
- RD-000: owner accepted `CONCEPTION.md` and ADR-0001 as written.
- RD-001: toolchain results recorded (see Verification and HANDOFF).
- RD-002 identity work implemented: module, config, icons, `bin/RouteDeck.app`. Owner Dock icon via `wails3 dev` PASS. `golangci-lint run ./...` PASS (`build/` excluded; package comment on `greetservice.go`).

## In progress

- None (RD-002 DONE)

## Blockers

- ProxyBridge install path list not frozen (RD-005)
- Phase 0 leftovers: RD-003 (ADR accept), RD-005 (ProxyBridge freeze) still TODO
- None for RD-002. Phase 1 exit criteria are met.

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

## Next actions

1. Owner reviews RD-003 in `docs/project/BACKLOG.md` (ADR-0001..0006). That task stays TODO until the owner accepts or lists edits.

## Verification

| Command | Result |
| --- | --- |
| `go mod vendor` | PASS |
| `wails3 task common:generate:icons` | PASS — PNG only; no `-iconcomposerinput` |
| `wails3 task common:update:build-assets` | PASS |
| `wails3 task common:generate:bindings` | PASS — `frontend/bindings/github.com/macrox-pro/route-deck` |
| `wails3 task darwin:package` | PASS — `bin/RouteDeck.app`; `Contents/MacOS/route-deck`; ld macOS 12.0 warnings only |
| `go test ./...` | PASS — no test files; module `github.com/macrox-pro/route-deck` |
| `golangci-lint run ./...` | PASS — 0 issues. `build/` excluded. Cache write warnings only (sandbox) |
| PlistBuddy CFBundleIdentifier | PASS — `pro.macrox.routedeck` |
| PlistBuddy CFBundleName | PASS — `RouteDeck` |
| PlistBuddy LSMinimumSystemVersion | PASS — `13.0.0` |
| Resources | PASS — `icons.icns` present; `Assets.car` absent; no `CFBundleIconName` |
| `shasum` brand vs `build/appicon.png` | PASS — identical |
| leftover template/rejected-id grep | PASS — clean (excl. vendor, docs, `project.nsi` comments) |
| Owner Dock icon (`wails3 dev`) | PASS — owner confirmed 2026-11-16 |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go declared `1.26.0`; Wails v3.0.0
