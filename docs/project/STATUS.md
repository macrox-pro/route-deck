# Project Status

Updated: 2026-11-17 16:50 UTC
Updated by: agent (RD-045 native dialogs)

## Current phase

Native GUI-conflict and extension-approval dialogs shipped locally; awaiting owner smoke.

## Current task

RD-045

## State

REVIEW

## Completed

- RD-045: Native dialogs for GUI conflict + extension approval (REVIEW)
- RD-044: Warn/Error severity for failure logs (REVIEW)
- RD-040..RD-043: Setup guides batch (REVIEW)
- RD-039: HTTP CONNECT health
- RD-007: CI workflow
- RD-038: Activity JSON log
- RD-035..RD-037 tray + Launch at Login
- Prior: RD-000..RD-034

## In progress

- None. Waiting on owner dialog smoke and commits.

## Blockers

- RD-033 Apply→Active needs owner UI + extension approval (SOCKS5 endpoint is ready: Shadowrocket 127.0.0.1:1082, health TCPAvailableUDPUnavailable)

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

## Next actions

1. Owner: `wails3 dev` → Apply with ProxyBridge GUI open → Warning dialog; quit GUI.
2. Owner: first-run Apply → Info dialog for extension approval; Allow in System Settings; Apply reaches Active.
3. Owner commits (agent MUST NOT commit).

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/adapter/outbound/proxybridge/ ./internal/application/service/ ./internal/adapter/inbound/wails/` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |

Environment: macOS darwin/arm64; module `github.com/macrox-pro/route-deck`; Go `1.26.0`; Wails v3.0.0
