# Project Status

Updated: 2026-11-17 18:55 UTC
Updated by: agent (RD-048 profiles app editor)

## Current phase

Profiles app-grouped editor shipped locally; awaiting owner smoke.

## Current task

RD-048

## State

REVIEW

## Completed

- RD-048: Profiles app-grouped editor + Review helpers preselect (REVIEW)
- RD-047: Profiles form layout + add-all (REVIEW; editor UX superseded by RD-048)
- RD-046: Protection UX clarity (REVIEW)
- RD-045: Native dialogs for GUI conflict + extension approval (REVIEW; may be uncommitted)
- RD-044: Warn/Error severity for failure logs (REVIEW)
- RD-040..RD-043: Setup guides batch (REVIEW)
- RD-039: HTTP CONNECT health
- RD-007: CI workflow
- RD-038: Activity JSON log
- RD-035..RD-037 tray + Launch at Login
- Prior: RD-000..RD-034

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
- Helper processes are never auto-confirmed (ADR-0008); Review helpers is explicit
- Proxy client icons only when host/port matches a documented Setup guide endpoint
- ConfirmHelpers remains additive (uncheck does not remove catalogued helpers)

## Next actions

1. Owner: `wails3 dev` smoke per HANDOFF (Profiles app list + Apply to helpers + Review helpers preselect).
2. Owner commits. Agent MUST NOT commit.

## Verification

| Command | Result |
| --- | --- |
| `cd frontend && npm run build` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
