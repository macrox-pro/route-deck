# Project Status

Updated: 2026-11-16 20:16 UTC
Updated by: agent (vendor ADR-0006)

## Current phase

Phase 1 Toolchain

## Current task

RD-002

## State

TODO (awaiting owner start; not READY)

## Completed

- Docs-prep files created in this working tree. Not committed. Owner commits.
- RD-000: owner accepted `CONCEPTION.md` and ADR-0001 as written.
- RD-001: toolchain results recorded (see Verification and HANDOFF).

## In progress

- None

## Blockers

- `go.mod` module is still `changeme` (owner RD-002)
- Bundle ID in generated plist may not yet be `pro.macrox.routedeck` (owner RD-002)
- ProxyBridge install path list not frozen (RD-005)
- Phase 0 leftovers: RD-003 (ADR accept), RD-005 (ProxyBridge freeze) still TODO
- Phase 1 exit also needs golangci-lint green on the scaffold; currently FAIL (5 revive issues)

## Known limitations

- No domain, persistence, or ProxyBridge adapter code yet
- Wails template still includes `GreetService` and a demo time event
- ESLint is not part of the frontend scaffold and is not required
- Owner added `vendor/` (`go mod vendor`); ADR-0006 Accepted — commit and refresh after dependency changes
- First `go test ./...` before a frontend build FAILs on missing `frontend/dist` embed

## Next actions

1. Owner promotes RD-002 to READY and starts module/bundle/icon identity.
2. Later: RD-003, RD-005; then clear scaffold revive findings so Phase 1 lint exit can pass.

## Verification

| Command | Result |
| --- | --- |
| `go version` | PASS — `go version go1.27.4 darwin/arm64` |
| `node -v` / `npm -v` | PASS — `v22.23.5` / `10.9.9` |
| `golangci-lint version` | PASS — `2.15.0` |
| `wails3 doctor` | PASS — exit 0; system ready (doctor Go: go1.27.3) |
| `go test ./...` (before build) | FAIL — `pattern all:frontend/dist: no matching files found` |
| `wails3 task build` | PASS — exit 0; `bin/route-deck` built |
| `go test ./...` (after build) | PASS — no test files; exit 0 |
| `golangci-lint run ./...` (after build) | FAIL — 5 revive issues in scaffold |

Environment: macOS darwin/arm64; `go.mod` still `module changeme`; Go declared `1.26.0`; Wails v3.0.0
