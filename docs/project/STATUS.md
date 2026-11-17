# Project Status

Updated: 2026-11-17 16:05 UTC
Updated by: agent (Phase 7 tray + autostart)

## Current phase

Phase 7 — tray lifecycle and Launch at Login shipped (RD-035..RD-037). RD-033 live smoke remains REVIEW (owner SOCKS5).

## Current task

RD-037

## State

REVIEW

## Completed

- RD-035: tray resident, close→Hide, Open, Quit dialog shell
- RD-036: live Status/tooltip, Stop when CLI owned, Quit dialog copy
- RD-037: Launch at Login via Wails Autostart
- RD-034: browse installed apps + icon serialization hotfix
- RD-018..RD-032 vertical slice
- Prior: RD-000..RD-017

## In progress

- None. Waiting on owner tray/autostart smoke and commits.

## Blockers

- RD-033 needs owner SOCKS5 for Apply→Active (unchanged)

## Known limitations

- HTTP proxy auth verification out of vertical slice (ADR-0008)
- Keychain prompts on unsigned rebuild expected
- Activity viewer later
- Autostart mechanism may be LaunchAgent under `wails3 task run` vs SMAppService when bundled
- Hundreds of indexed apps; filter-only (no virtualization); icons load under a process-wide mutex
- Dock icon stays while window is hidden (no public Wails API to toggle Dock only then)

## Next actions

1. Owner: `wails3 task run` → close window (stays in menu bar) → Open from tray → Quit dialog Cancel/Confirm → Launch at Login toggle.
2. Owner runs RD-033 smoke with SOCKS5 when ready.
3. Owner commits (agent MUST NOT commit).
4. Next: Activity viewer or RD-007 CI per owner.

## Verification

| Command | Result |
| --- | --- |
| `go test ./...` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `npm run build` (frontend/) | PASS |
| `wails3 generate bindings` | PASS |
| `wails3 task run` | NOT RUN (owner tray/autostart smoke) |
| RD-033 live smoke | NOT RUN (owner SOCKS5) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
