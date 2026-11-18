# Project Status

Updated: 2026-11-18 21:25 UTC
Updated by: agent (RD-059 Appearance + native frame sync → REVIEW)

## Current phase

Appearance preference in Settings with SQLite persistence and native window chrome recreate.

## Current task

RD-059

## State

REVIEW

## Completed

- Still in Active (REVIEW, not archived): RD-007, RD-033, RD-037..RD-059
- DONE archive: [`COMPLETED.md`](COMPLETED.md) (RD-000..RD-036 except REVIEW gaps)

## In progress

- None. Waiting on owner smoke/commits.

## Blockers

- RD-033 Apply→Active owner smoke still pending

## Known limitations

- Shadowrocket Mac is a system VPN; Direct for other apps may still tunnel
- Hiddify guide: Detection unavailable (no sourced Mac bundle ID)
- Port open is TCP only, not ProxyHealth
- Suggestion dismiss is in-memory until restart
- Historical log lines stay INFO until new writes
- HTTP CONNECT probes example.com:443
- Keychain prompts on unsigned rebuild expected
- Another CLI / extension timeout still have no dedicated native dialog (out of RD-045 scope)
- Dialog copy is English-only; message IDs ready for future i18n
- Helper processes are never auto-confirmed (ADR-0008 / ADR-0014); Review helpers is explicit
- Proxy client icons only when host/port matches a documented Setup guide endpoint
- Quitting with a dirty editor discards edits (ADR-0014 known limitation)
- Old duplicate ProcessIdentity rows for one bundle id are collapsed at compile time, not migrated in SQLite
- CI `build-macos` artifact is ad-hoc codesign only (Gatekeeper will warn); arm64 (`macos-14`) only
- Local `wails3 dev` always shows version `dev` (no ldflags injection in DEV builds)
- Appearance chrome recreate may reject the calling SetAppearance promise (persist-first; new window loads via Get)
- No automatic migration of prior localStorage theme into SQLite

## Next actions

1. Owner: `wails3 dev` smoke — Settings Appearance System/Light/Dark; title bar matches; dirty editor Cancel aborts.
2. Owner commits RD-059 (and prior REVIEW work as desired). Agent MUST NOT commit.

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/application/service/ ./internal/adapter/outbound/persistence/sqlite/ ./internal/adapter/inbound/wails/` | PASS |
| `golangci-lint run ./...` | PASS (0 issues; cache write warnings in sandbox) |
| `cd frontend && npm run build` | PASS |
| `wails3 generate bindings -ts` | PASS |
| macOS Appearance recreate smoke | NOT RUN (owner) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
