# Project Status

Updated: 2026-11-16 20:34 UTC
Updated by: agent (RD-002 owner icon PASS)

## Current phase

Phase 1 Toolchain

## Current task

RD-002

## State

REVIEW

## Completed

- Docs-prep files created in this working tree. Not committed. Owner commits.
- RD-000: owner accepted `CONCEPTION.md` and ADR-0001 as written.
- RD-001: toolchain results recorded (see Verification and HANDOFF).
- RD-002 identity work implemented: module, config, icons, `bin/RouteDeck.app`. Owner Dock icon via `wails3 dev` PASS. Lint FAIL still blocks DONE.

## In progress

- None (RD-002 in REVIEW)

## Blockers

- ProxyBridge install path list not frozen (RD-005)
- Phase 0 leftovers: RD-003 (ADR accept), RD-005 (ProxyBridge freeze) still TODO
- Phase 1 / RD-002 DONE blocked by golangci-lint FAIL (5 scaffold revive issues)

## Known limitations

- No domain, persistence, or ProxyBridge adapter code yet
- Wails template still includes `GreetService` and a demo time event
- ESLint is not part of the frontend scaffold and is not required
- `vendor/` committed per ADR-0006; refresh after dependency changes
- First `go test ./...` before a frontend build FAILs on missing `frontend/dist` embed
- Darwin icon is PNG→`.icns` only (no Icon Composer `Assets.car`)
- After `wails3 task common:update:build-assets`, re-apply `LSMinimumSystemVersion` `13.0.0` and drop `CFBundleIconName` when not shipping `Assets.car`
- Linker may warn about macOS 12.0 link minimum; plist floor is 13.0.0

## Next actions

1. Clear scaffold revive findings so `golangci-lint run ./...` PASSes (or leave RD-002 REVIEW and claim RD-003).
2. Later: RD-003, RD-005.

## Verification

| Command | Result |
| --- | --- |
| `go mod vendor` | PASS |
| `wails3 task common:generate:icons` | PASS — PNG only; no `-iconcomposerinput` |
| `wails3 task common:update:build-assets` | PASS |
| `wails3 task common:generate:bindings` | PASS — `frontend/bindings/github.com/macrox-pro/route-deck` |
| `wails3 task darwin:package` | PASS — `bin/RouteDeck.app`; `Contents/MacOS/route-deck`; ld macOS 12.0 warnings only |
| `go test ./...` | PASS — no test files; module `github.com/macrox-pro/route-deck` |
| `golangci-lint run ./...` | FAIL — 5 revive (greetservice, android deps, ios app_options, ios deps superfluous-else) |
| PlistBuddy CFBundleIdentifier | PASS — `pro.macrox.routedeck` |
| PlistBuddy CFBundleName | PASS — `RouteDeck` |
| PlistBuddy LSMinimumSystemVersion | PASS — `13.0.0` |
| Resources | PASS — `icons.icns` present; `Assets.car` absent; no `CFBundleIconName` |
| `shasum` brand vs `build/appicon.png` | PASS — identical |
| leftover template/rejected-id grep | PASS — clean (excl. vendor, docs, `project.nsi` comments) |
| Owner Dock icon (`wails3 dev`) | PASS — owner confirmed 2026-11-16 |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go declared `1.26.0`; Wails v3.0.0
