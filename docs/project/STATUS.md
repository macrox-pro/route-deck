# Project Status

Updated: 2026-11-17 09:05 UTC
Updated by: agent (RD-014..RD-017 DONE)

## Current phase

Phase 3 Persistence + launchable shell — RD-014..RD-017 DONE. Next: owner review / commit.

## Current task

RD-017

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
- RD-006: folded into RD-014; GORM + sqlite driver vendored; AutoMigrate smoke PASS.
- RD-014: SQLite catalog adapter, five tables, ports List/Save, `ErrUnsetAction`, round-trip tests. Lint and tests PASS.
- RD-015: Catalog and Snapshotter services; Unknown snapshot without CLI; DB-failed path. Lint and tests PASS.
- RD-016: Greet/time removed; SnapshotService wired in main; bindings regenerated. Lint and tests PASS.
- RD-017: HashRouter shell, six pages, ProtectionBadge, Dashboard loads snapshot. `npm run build` and `wails3 build` PASS.

## In progress

- None. Owner reviews launchable shell and commits.

## Blockers

- None.

## Known limitations

- No ProxyBridge adapter, Keychain adapter, or health checks yet
- ESLint is not part of the frontend scaffold and is not required
- `vendor/` committed per ADR-0006; refresh after dependency changes
- Darwin icon is PNG→`.icns` only (no Icon Composer `Assets.car`)
- `build/` is excluded from golangci-lint linters and formatters (`^build/` in `.golangci.yml`). It is the Wails scaffold
- After `wails3 task common:update:build-assets`, re-apply `LSMinimumSystemVersion` `13.0.0` and drop `CFBundleIconName` when not shipping `Assets.car`
- Linker may warn about macOS 12.0 link minimum; plist floor is 13.0.0
- Homebrew and other non-documented ProxyBridge install paths stay “not documented; do not guess”
- Id parsers reject only `""`; whitespace-only ids are accepted until a later task revisits
- `domain-model.md` spells `RoutingAction` as `PROXY`/`DIRECT`/`BLOCK`; RD-008 uses opaque constructors with no string form
- Opaque `RoutingAction` would serialize as `{}`; future Wails inbound must map it to a DTO
- `ProxyEndpoint` has exported fields; a struct literal can bypass `NewProxyEndpoint`
- `Profile` has exported fields; a struct literal can bypass `NewProfile`
- Untitled UI brand scale in `theme.css` is still default purple (branding task later)
- Dashboard forms for proxies/profiles are empty shells
- Product events not emitted yet

## Next actions

1. Owner reviews `wails3 dev` / `bin/RouteDeck.app` and commits RD-014..017 (and vendor).
2. Next product stage: ProxyBridge CLI adapter or Keychain adapter per owner.

## Verification

| Command | Result |
| --- | --- |
| `go test ./...` | PASS |
| `golangci-lint run ./...` | PASS — 0 issues |
| `npm run build` (cwd `frontend/`) | PASS |
| `wails3 build` | PASS — `bin/route-deck` and app bundle |
| Launch smoke | PASS — `~/Library/Application Support/route-deck/route-deck.db` created with five tables |
| macOS-only checks | PASS — darwin/arm64 launch created Application Support DB |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go declared `1.26.0`; Wails v3.0.0; golangci-lint 2.15.0; gorm.io/gorm v1.31.2; gorm.io/driver/sqlite v1.6.0; ProxyBridge CLI 4.1.3 at documented paths
