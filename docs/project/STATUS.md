# Project Status

Updated: 2026-11-17 08:04 UTC
Updated by: agent (RD-012 DONE)

## Current phase

Phase 2 Domain ports — RD-011 and RD-012 DONE. Next: owner Phase 2 exit review. Do not open RD-006 or RD-007.

## Current task

RD-012

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
- RD-011: CONCEPTION 14-step derivation, Inactive/Initializing glossary rows, `DeriveProtection`, 21 table cases. Lint and tests PASS.
- RD-012: `ProxyStore.Get` and `SecretStore.Load` in `internal/application/port`, fake assignment test. Lint and tests PASS.

## In progress

- None. RD-012 DONE. Owner reviews Phase 2 exit.

## Blockers

- None. No Phase 2 exit checklist file exists. Do not open RD-006 or RD-007 without owner direction.

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
- `DatabaseDamaged`, `!VersionKnown`, referenced `Unknown`/`Checking`, `Stopping`, `Detected`, and applied-none are in the step 2/3/4/7 conditions and have no own row
- Exported `ProtectionFacts` can skip a constructor
- A proxy id missing from `Proxies` does not change the protection word

## Next actions

1. Owner reviews Phase 2 exit. No checklist file exists. Do not open RD-006 or RD-007.
2. Owner commits RD-010, RD-011, and RD-012 when ready.

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/application/port/...` | PASS |
| `go list` imports for `./internal/application/port/...` | PASS — `context` and `internal/domain` only |
| `golangci-lint run ./...` | PASS — 0 issues |
| `wails3 build` | NOT RUN — `main` does not import these packages |
| macOS-only checks | NOT RUN — no OS calls |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go declared `1.26.0`; Wails v3.0.0; golangci-lint 2.15.0; ProxyBridge CLI 4.1.3 at documented paths
