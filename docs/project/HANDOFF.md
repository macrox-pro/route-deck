# Handoff

Timestamp: 2026-11-19 21:45 UTC
Agent: agent (RD-075 + RD-076)
Task ID: RD-076 (DONE)

## Done
- RD-075: proxy Connection/Authentication + Kind ButtonGroup; profile sticky toolbar, Action/Protocol ButtonGroup, proxy Select; helpers tri-state Select all; dialog shell/footer spacing
- RD-076: `app_menu.go` RouteDeck menu (SetApplicationMenu), enable rules while editor open + Stop via CLIOwned, table-driven `MenuItemEnabled` tests; wired in `main.go`
- UX follow-ups: Setup Guides in sidebar footer; AriaButton-unified sidebar; removed duplicate empty/header CTAs

## Verification
| Command | Result |
| --- | --- |
| `cd frontend && npm run build` | PASS |
| `go test ./internal/adapter/inbound/...` | PASS |
| `golangci-lint run ./...` | PASS |
| `wails3 task build` | PASS |
| Manual menu + accelerators (dirty editor Cmd+W / Quit) | NOT RUN (needs owner on macOS UI) |

## Known limitations
- Cmd+W handler in `dialog-footer.tsx` kept: CloseWindow role closes the focused window, but the footer path still runs dirty cancel semantics explicitly.
- Manual full menu pass not run this session; DONE for code+automated gates; owner should smoke-test menu items.

## Exact next step
Owner: `git status` / review diff, then commit. Optional: run `wails3 task build` and click every menu item once.

## Suggested commit message
```text
add routedeck menu bar and finish editor form ux for rd-075/rd-076
```
