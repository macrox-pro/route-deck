# Project Status

Updated: 2026-11-17 22:30 UTC
Updated by: agent (RD-049..RD-056 editor dialogs epic → REVIEW)

## Current phase

Editor dialogs, helper groups, and Dock badge epic implemented locally; awaiting owner smoke and commits.

## Current task

RD-056

## State

REVIEW

## Completed

- Still in Active (REVIEW, not archived): RD-007, RD-033, RD-037..RD-056
- DONE archive: [`COMPLETED.md`](COMPLETED.md) (RD-000..RD-036 except REVIEW gaps)

## In progress

- None. Waiting on owner UX smoke and commits.

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

## Next actions

1. Owner: `wails3 dev` smoke per HANDOFF (profile/proxy/helpers dialogs, sync confirm, Dock `!` badge).
2. Owner commits (suggested messages in HANDOFF). Agent MUST NOT commit.

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/...` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `cd frontend && npm run build` | PASS |
| `wails3 build` | PASS |
| macOS dialog/Dock smoke | NOT RUN (owner) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
