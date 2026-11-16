# Project Status

Updated: 2026-11-16 19:47 UTC
Updated by: agent (RD-000 accept)

## Current phase

Phase 0 Contracts → Phase 1 Toolchain

## Current task

RD-001

## State

READY

## Completed

- Docs-prep files created in this working tree (see HANDOFF). Not committed. Owner commits.
- RD-000: owner accepted `CONCEPTION.md` and ADR-0001 as written. Fail-open, UDP, and kill-switch bounds unchanged.

## In progress

- None

## Blockers

- `go.mod` module is still `changeme` (owner RD-002)
- Bundle ID in generated plist may not yet be `pro.macrox.routedeck` (owner RD-002)
- ProxyBridge install path list not frozen (RD-005)

## Known limitations

- No domain, persistence, or ProxyBridge adapter code yet
- Wails template still includes `GreetService` and a demo time event
- ESLint is not part of the frontend scaffold and is not required
- golangci-lint v2.15.0 installed and `golangci-lint run ./...` FAIL on missing `frontend/dist` embed (scaffold). Recorded in HANDOFF. Config is in place.

## Next actions

1. Run RD-001 verification commands and record output in STATUS and HANDOFF.
2. After RD-001, owner starts RD-002 (`go.mod`, `build/config.yml`, app icon install).

## Verification

- Command: owner review of CONCEPTION.md + ADR-0001 (RD-000)
- Result: PASS — accepted 2026-11-16
- Environment: macOS, Go 1.26.0 declared in `go.mod`, Wails v3.0.0 required
