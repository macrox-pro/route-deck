# Project Status

Updated: 2026-11-16 22:15 UTC
Updated by: agent (RD-010 DONE)

## Current phase

Phase 2 Domain ports — RD-010 Profile and ValidateProxyRefs DONE. Next: owner promotes RD-011 or RD-012 (only one READY at a time).

## Current task

RD-010

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
- RD-013: lint now enables the `gofmt` formatter and the `exhaustive` linter; RD-008 parsers convert once; `state.go` is `gofmt`-clean; rules in `docs/engineering/golang.md` and `.cursor/rules/golang.mdc`. Committed as `6264c28`.
- RD-009: `Application`, `ProcessIdentity` (`MatchableNames`), `ProxyEndpoint` (`NewProxyEndpoint`), `ErrInvalidPort`, `ErrEmptyHost`. Unit tests only for public methods/constructors. Lint and tests PASS.
- RD-010: `ApplicationPolicy`, `Profile`, `NewProfile`, `ValidateProxyRefs`, `ErrEmptyName`, `ErrUnresolvedProxy`. Lint and tests PASS.

## In progress

- None. RD-010 DONE. Owner promotes RD-011 or RD-012 after DoR check (only one READY at a time).

## Blockers

- None for promoting RD-011 or RD-012 after their DoR checks. CONCEPTION is not edited until RD-011.

## Known limitations

- No persistence or ProxyBridge adapter code yet (RD-011+)
- Wails template still includes `GreetService` and a demo time event
- ESLint is not part of the frontend scaffold and is not required
- `vendor/` committed per ADR-0006; refresh after dependency changes
- First `go test ./...` before a frontend build FAILs on missing `frontend/dist` embed
- Darwin icon is PNG→`.icns` only (no Icon Composer `Assets.car`)
- `build/` is excluded from golangci-lint linters and formatters (`^build/` in `.golangci.yml`). It is the Wails scaffold
- After `wails3 task common:update:build-assets`, re-apply `LSMinimumSystemVersion` `13.0.0` and drop `CFBundleIconName` when not shipping `Assets.car`
- Linker may warn about macOS 12.0 link minimum; plist floor is 13.0.0
- Homebrew and other non-documented ProxyBridge install paths stay “not documented; do not guess”
- Phase 2 exit waits until RD-011 is DONE (owner choice)
- Id parsers reject only `""`; whitespace-only ids are accepted until a later task revisits
- `domain-model.md` spells `RoutingAction` as `PROXY`/`DIRECT`/`BLOCK`; RD-008 uses opaque constructors with no string form
- Opaque `RoutingAction` would serialize as `{}`; future Wails inbound must map it to a DTO
- Open owner review items, not scheduled: zero `RoutingAction` has no validity check (a caller's `else` branch could treat it as Direct); five copy-paste id parsers in `id.go`; `ProtocolScopeBOTH` naming; `ErrEmptyID` wording for `KeychainRef`; `ProxyHealth` tests lack empty and case-variant tokens
- Lint shows one issue per line; a `gofmt` issue can hide behind another issue on the same line
- `ProxyEndpoint` has exported fields; a struct literal can bypass `NewProxyEndpoint` (invalid port, empty host, or `SendDomain: false`)
- `NewProxyEndpoint` does not validate `Kind` or empty `ID`; whitespace-only host is accepted (owner decision)
- `MatchableNames` does not deduplicate; equal bundle id and executable name yield two equal entries
- Unit tests target exported functions/methods only; plain structs with no behavior have no dedicated tests (owner, 2026-11-17; `docs/engineering/testing.md`)
- `Profile` has exported fields; a struct literal can bypass `NewProfile` and store an empty name
- `NewProfile` does not validate empty `ID` or policy contents; whitespace-only name is accepted (owner decision)
- `ValidateProxyRefs` does not check `Name` or `ID`; empty name is only a `NewProfile` error
- Package comment in `errors.go` still says “no aggregates” after Profile was added (owner decision not to edit it)
- Skeleton in `docs/engineering/testing.md` shows `ValidateProxyRefs` with one argument; real signature takes `known []ProxyID`

## Next actions

1. Open `docs/project/tasks/RD-011.md` or `docs/project/tasks/RD-012.md` and check DoR before promoting one to READY.
2. Owner commits RD-010 when ready.

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/domain/...` | PASS |
| `golangci-lint run ./...` | PASS — 0 issues. Cache write warnings only (sandbox) |
| `go test ./...` | PASS |
| `go list` imports for `./internal/domain/...` | PASS — `errors fmt strings` only |
| `git status --short go.mod go.sum vendor/` | PASS — empty |
| `gofmt -l internal/domain/` | PASS — empty |
| `wails3 build` | NOT RUN — `main` does not import `internal/domain` |
| macOS-only checks | NOT RUN — pure Go structs with no OS calls |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go declared `1.26.0`; Wails v3.0.0; golangci-lint 2.15.0; ProxyBridge CLI 4.1.3 at documented paths
