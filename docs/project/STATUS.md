# Project Status

Updated: 2026-11-17 20:05 UTC
Updated by: agent (RD-058 buildmeta version/commit → REVIEW)

## Current phase

CI packages unsigned `.app` with tag/commit identity; editor-dialogs epic still awaiting owner smoke/commits.

## Current task

RD-058

## State

REVIEW

## Completed

- Still in Active (REVIEW, not archived): RD-007, RD-033, RD-037..RD-058
- DONE archive: [`COMPLETED.md`](COMPLETED.md) (RD-000..RD-036 except REVIEW gaps)

## In progress

- None. Waiting on owner commits and remote `build-macos` evidence.

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

## Next actions

1. Owner commits RD-057 + RD-058 (and prior REVIEW work as desired). Agent MUST NOT commit.
2. Owner: push / tag `v*` and confirm Actions artifact + Settings About / sidebar label.
3. Owner: `wails3 dev` smoke for RD-049..RD-056 if not done.

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/buildmeta/` | PASS |
| `golangci-lint run ./...` | PASS (0 issues; cache write warnings in sandbox) |
| `cd frontend && npm run build` | PASS |
| Remote `build-macos` with tag | NOT RUN (not pushed) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
