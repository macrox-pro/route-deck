# Handoff

Timestamp: 2026-11-17 14:45 UTC
Agent: agent (RD-044 log severity)
Task ID: RD-044
Branch / commit: working tree (agent MUST NOT commit)

## Done

- Owner severity scale recorded in `docs/engineering/logging.md` (Info / Warn / Error)
- Retargeted failure call sites under `internal/` and `main.go` from Info to Warn or Error
- `service.ErrOperationInProgress` sentinel; Apply/Stop busy returns `errors.Is`
- Apply/Stop Wails logging: Warn for busy / not installed / conflict; Error otherwise
- Warn when Apply skipped (desired profile not set)
- Classified ProxyBridge CLI lines stay Info (ADR-0008)

## Verification

| Command | Result |
| --- | --- |
| `rg 'logger\(\)\.Info\("[^"]*fail' internal/ main.go` | PASS (no hits) |
| `go test ./internal/application/service/ ./internal/adapter/inbound/wails/ ./internal/adapter/outbound/keychain/ ./internal/adapter/outbound/proxybridge/ ./internal/adapter/inbound/tray/` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |

## Exact next step

Owner: restart or refresh `wails3 dev`, trigger a known failure (e.g. Apply without desired profile), confirm Activity shows WARN (not INFO). Then commit.

## Suggested commit message

```text
log operational failures as warn and infra failures as error (rd-044)
```

## Context to read

- `docs/engineering/logging.md`
- `docs/project/tasks/RD-044.md`
- `internal/adapter/inbound/wails/protection.go`
