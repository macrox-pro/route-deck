# Project Status

Updated: 2026-11-16 21:42 UTC
Updated by: agent (RD-008 DONE)

## Current phase

Phase 2 Domain ports — RD-008 value types DONE. Next: check DoR on RD-009.

## Current task

RD-008

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
- Phase 2 backlog accepted: task bodies in `docs/project/tasks/RD-*.md`; BACKLOG is an index only.
- RD-008: `internal/domain` ids, states, policy value types, opaque `RoutingAction`, three sentinels, table tests. Lint and tests PASS.

## In progress

- None. RD-008 DONE. Owner promotes RD-009 after DoR check.

## Blockers

- None for promoting RD-009 after its DoR check. RD-010 through RD-012 wait on prior DONE. CONCEPTION is not edited until RD-011.

## Known limitations

- No aggregates, persistence, or ProxyBridge adapter code yet (RD-009+)
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
- Id parsers reject only `""`; whitespace-only ids are accepted until a later task revisits
- `domain-model.md` spells `RoutingAction` as `PROXY`/`DIRECT`/`BLOCK`; RD-008 uses opaque constructors with no string form
- Opaque `RoutingAction` would serialize as `{}`; future Wails inbound must map it to a DTO

## Next actions

1. Open `docs/project/tasks/RD-009.md` and check DoR before promoting to READY.
2. Owner commits RD-008 when ready.

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/domain/...` | PASS |
| `go list` imports for `./internal/domain/...` | PASS — `errors fmt` only |
| `golangci-lint run ./...` | PASS — 0 issues. Cache write warnings only (sandbox) |
| `go test ./...` | PASS |
| `git status --short go.mod go.sum vendor/` | PASS — empty |
| `wails3 generate bindings ...` then bindings status | PASS — empty; 1 Service, 0 Models |
| `wails3 build` | NOT RUN — `main` does not import `internal/domain` |
| macOS-only checks | NOT RUN — pure Go value types |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go declared `1.26.0`; Wails v3.0.0; ProxyBridge CLI 4.1.3 at documented paths
