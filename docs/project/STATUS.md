# Project Status

Updated: 2026-11-17 15:10 UTC
Updated by: agent (icon crash hotfix)

## Current phase

Phase 5+ — installed application browse list shipped; icon crash hotfix applied. RD-033 live smoke remains REVIEW.

## Current task

RD-033

## State

REVIEW

## Completed

- RD-034: browse installed apps + icon serialization hotfix (AppKit race)
- RD-018..RD-032 vertical slice
- Prior: RD-000..RD-017

## In progress

- None. Waiting on owner confirm Applications no longer crashes, then RD-033 smoke / commits.

## Blockers

- RD-033 needs owner SOCKS5 for Apply→Active

## Known limitations

- HTTP proxy auth verification out of vertical slice (ADR-0008)
- Keychain prompts on unsigned rebuild expected
- Activity / tray later
- Hundreds of indexed apps; filter-only (no virtualization); icons load under a process-wide mutex

## Next actions

1. Owner: `wails3 task run` → open Applications → confirm no crash.
2. Owner runs RD-033 smoke with SOCKS5 when ready.
3. Owner commits (agent MUST NOT commit).
4. Next phase: tray (E11) or RD-007 CI per owner.

## Verification

| Command | Result |
| --- | --- |
| `go test ./...` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `npm run build` (frontend/) | PASS |
| `wails3 build` | NOT RUN (use `wails3 task run` to confirm) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
