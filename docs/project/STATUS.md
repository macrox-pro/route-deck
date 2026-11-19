# Project Status

Updated: 2026-11-19 20:01 UTC
Updated by: agent (RD-062 → DONE archive)

## Current phase

Idle after RD-062. Ready for next `RD-*` when owner opens one.

## Current task

None

## State

IDLE

## Completed

- DONE archive: [`COMPLETED.md`](COMPLETED.md) (RD-000..RD-062)
- Active backlog empty

## In progress

- None

## Blockers

- None

## Known limitations

- Shadowrocket Mac is a system VPN; Direct for other apps may still tunnel
- Hiddify guide: Detection unavailable (no sourced Mac bundle ID)
- Port open is TCP only, not ProxyHealth
- Suggestion dismiss is in-memory until restart
- Checklist dismiss is in-memory until restart (RD-061)
- Checklist UI copy is English on the frontend via stable message IDs (future i18n)
- Historical log lines stay INFO until new writes
- HTTP CONNECT probes example.com:443
- Keychain prompts on unsigned rebuild expected
- Another CLI / extension timeout still have no dedicated native dialog (out of RD-045 scope)
- Dialog copy is English-only; message IDs ready for future i18n
- Helper processes are never auto-confirmed (ADR-0008 / ADR-0014); Review helpers is explicit
- Proxy client icons only when host/port matches a documented Setup guide endpoint
- Quitting with a dirty editor discards edits (ADR-0014 known limitation)
- Old duplicate ProcessIdentity rows for one bundle id are collapsed at compile time, not migrated in SQLite
- CI `build-macos` DMG is ad-hoc codesign only (Gatekeeper will warn); arm64 (`macos-15`) only; no notarization
- Local `wails3 dev` always shows version `dev` (no ldflags injection in DEV builds)
- Appearance chrome recreate may reject the calling SetAppearance promise (persist-first; new window loads via Get)
- No automatic migration of prior localStorage theme into SQLite
- Auto Apply at login may fail honestly if the local proxy is not listening yet (no startup retry; ADR-0015)
- Complete with UDP-limited proxy may hide checklist while Protection stays Degraded (intentional; partial counts)
- RD-033 edge scenarios (GUI Conflict, proxy-down, CLI kill, wake) not re-run in the 2026-11-19 critical smoke pass
- Tag → GitHub Release smoke for RD-062 left to owner (`git push origin v*`)

## Next actions

1. Owner commits this archive (agent MUST NOT commit).
2. Owner: push a new `v*` tag and confirm Actions Release + DMG (if not already done).
3. Owner: open next `docs/project/tasks/RD-*.md` when scope is ready → BACKLOG Active + STATUS.

## Verification

| Command | Result |
| --- | --- |
| Archive RD-062 into COMPLETED.md | PASS |
| Remove `docs/project/tasks/RD-062.md` | PASS |
| `golangci-lint run ./...` (prior RD-062) | PASS (0 issues) |
| tag → GitHub Release smoke | NOT RUN (owner) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
