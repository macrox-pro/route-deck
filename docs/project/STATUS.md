# Project Status

Updated: 2026-11-19 20:43 UTC
Updated by: agent (RD-068 → DONE)

## Current phase

Idle after profile UX bugfix (RD-068).

## Current task

None

## State

IDLE

## Completed

- DONE archive: [`COMPLETED.md`](COMPLETED.md) (RD-000..RD-068)
- Active backlog empty

## In progress

- None

## Blockers

- None

## Known limitations

- Shadowrocket Mac is a system VPN; Direct for other apps may still tunnel
- Hiddify guide: Detection unavailable (no sourced Mac bundle ID)
- Port open is TCP only, not ProxyHealth
- Checklist UI copy is English on the frontend via stable message IDs (future i18n)
- Historical log lines stay INFO until new writes
- HTTP CONNECT probes example.com:443
- Keychain prompts on unsigned rebuild expected
- Dialog copy is English-only; message IDs ready for future i18n
- Helper processes are never auto-confirmed (ADR-0008 / ADR-0014); Review helpers is explicit
- Proxy client icons only when host/port matches a documented Setup guide endpoint
- Old duplicate ProcessIdentity rows for one bundle id are collapsed at compile time, not migrated in SQLite
- CI `build-macos` DMG is ad-hoc codesign only (Gatekeeper will warn); arm64 (`macos-15`) only; no notarization
- Local `wails3 dev` always shows version `dev` (no ldflags injection in DEV builds)
- Appearance chrome recreate may reject the calling SetAppearance promise (persist-first; new window loads via Get)
- No automatic migration of prior localStorage theme into SQLite
- Complete with UDP-limited proxy may hide checklist while Protection stays Degraded (intentional; partial counts)
- Tag → GitHub Release smoke for RD-062 left to owner (`git push origin v*`)
- Auto Apply does not retry when Apply succeeds with Degraded while the proxy is down (ADR-0017)
- Stop timeout / another CLI still have no dedicated native dialog (extension timeout covered by RD-066)
- Last-wins match conflicts may drop earlier policies that share a matchable name (ADR-0019; intentional)
- Deleting the Applied profile stops protection after the "Stop and Delete" confirm (ADR-0019)

## Next actions

1. Owner commits working tree including RD-063..RD-068 if not yet committed (agent MUST NOT commit).
2. Owner: optional `git tag` / Release smoke if not done for RD-062.
3. Owner: open next `docs/project/tasks/RD-*.md` when scope is ready.

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/domain ./internal/application/service ./internal/adapter/inbound/wails` | PASS |
| `cd frontend && npm run build` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
