# Handoff

Timestamp: 2026-11-17 15:25 UTC
Agent: agent (portable discovered ProxyBridge path + prior Applications polish)
Task ID: RD-033 (current)
Branch / commit: working tree (agent MUST NOT commit)

## Done

- Applications: Catalog-first UI, system filter via portable `isSystem`, icon crash fix
- Settings: `discoveredProxyBridgePath` placeholder via application `Settings` + `BridgeLocator`
- `BridgeLocator.PlatformPaths()`; macOS paths in `paths_darwin.go`, empty on `!darwin`
- Wails SettingsService thinned to DTO mapping over `service.Settings`

## Files changed (high level)

- `internal/application/port/{discover,bridge}.go`
- `internal/application/service/{settings,catalog,protection_test}.go`
- `internal/adapter/outbound/proxybridge/{locator,paths_*}.go`
- `internal/adapter/outbound/macos/lister*.go`
- `internal/adapter/inbound/wails/{settings,applications}.go`
- `main.go`, `frontend/src/pages/{settings,applications}.tsx`, bindings
- ADR-0008/0009, wails-api, STATUS/HANDOFF

## Decisions

- Discovered path = `Locate("")` over adapter `PlatformPaths` only; UI must not hard-code install paths.
- macOS `IsSystem` = outside `/Applications` and `~/Applications`.

## Verification

| Command | Result |
| --- | --- |
| `go test ./...` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `npm run build` (cwd `frontend/`) | PASS |
| `wails3 generate bindings` | PASS |
| `wails3 task run` | NOT RUN (owner confirm Settings placeholder + Applications) |
| RD-033 live smoke | NOT RUN (owner SOCKS5) |

## Not done

- RD-033 interactive scenarios
- Agent did not commit

## Exact next step

Owner: `wails3 task run` → Settings (placeholder shows discovered CLI path) and Applications (system checkbox).

## Suggested commit message

```text
expose portable bridge discovery and isSystem flags through ports for settings and browse ui
```

## Context to read

- `AGENTS.md`
- `docs/decisions/ADR-0008-vertical-slice-decisions.md`
- `docs/decisions/ADR-0009-installed-application-browse.md`
- `docs/project/tasks/RD-033.md`
