# Logging

RouteDeck uses the Go standard library `log/slog`. Services and adapters receive a `*slog.Logger` from the composition root. The replaceable part is the `slog.Handler`. There is no custom logger port.

Owner decisions: ADR-0008. Secret bans: `docs/security/secrets.md`.

## Destinations

| Destination | Format | Role |
| --- | --- | --- |
| `~/Library/Logs/route-deck/route-deck.log` | JSON lines | Durable diagnostics |
| stderr | text | `wails3 dev` console |
| Wails `application.Options.Logger` | same logger | Framework messages |

Both file and stderr handlers are combined with `slog.NewMultiHandler` (Go 1.26+).

If the log directory cannot be created, the process continues with stderr only and logs that degradation once.

## Rotation

Package `internal/adapter/outbound/logfile` implements an `io.Writer`:

- Cap: 5 MiB per file
- Keep: 3 previous files (`route-deck.log.1` .. `.3`)
- Safe for concurrent writes; closed on shutdown

No third-party rotation dependency.

## Level

Default: Info. Settings stores a debug flag in SQLite. Changing it updates a shared `slog.LevelVar` without restart.

Severity for product events (owner 2026-11-17):

| Level | When |
| --- | --- |
| Info | Normal flow; every classified ProxyBridge CLI line (truncated) |
| Warn | Operational / expected: cancelled, desired profile not set, operation in progress, bridge not installed, conflict, probe and user-facing validation failures |
| Error | Unexpected infrastructure: Keychain, DB/settings/snapshot, apply/stop internal failure |

## Redaction

A `ReplaceAttr` on every handler replaces attribute values whose keys are `password`, `secret`, `token`, or `profile` with `[REDACTED]`. This is defense in depth. Callers MUST NOT pass secrets or the profile body to the logger.

## Stable attribute keys

| Key | Meaning |
| --- | --- |
| `component` | Package or subsystem (`sqlite`, `proxybridge`, `coordinator`, …) |
| `operation_id` | Long Apply/Stop operation |
| `proxy_id` | Catalog proxy id |
| `profile_id` | Profile id |
| `bridge_state` | BridgeState token |
| `protection_state` | ProtectionState token |
| `health` | ProxyHealth token |
| `line_class` | CLI stdout/stderr classification |
| `err` | Error value |

## What to log

### Info

- Startup: app identity, DB path, log path (no secrets)
- Bridge discovery result (path, version, compatibility warning text as attribute)
- Proxy health transitions (old → new)
- Apply/Stop start, confirmed, stop done with `operation_id`
- Confirm and rollback outcomes (success path)
- CLI exit code and classification
- Every classified CLI line (truncated)
- Protection state transitions
- Wake and network-change events

### Warn

- Typed operational errors returned to the UI (validation, referenced delete, already added, cancelled)
- Desired profile not set; operation in progress
- Bridge not installed; GUI/CLI conflict
- Bridge version probe failure while trying the next path
- Proxy health check / probe failure
- Tray Stop failure when busy
- Open docs / open logs folder failure
- Log path or log file unavailable (stderr-only fallback)
- Log file close failure on shutdown

### Error

- Keychain load/save/delete failure
- Database open/close or settings load failure
- Settings/DB read failure on Apply
- Snapshot build/get failure
- Apply/Stop internal failure (start/confirm/rollback/process), not the operational cases above
- Tray Quit Stop failure when the managed CLI cannot stop

### Debug

- Probe timings
- Snapshot build details
- Wails service call entry/exit

### Never

- Passwords, tokens, Keychain item bytes
- Full ProxyBridge profile JSON
- User traffic / connection payloads
- GORM SQL (GORM logger stays Silent)

### Domain

The domain package MUST NOT import `log/slog` or write logs. Coverage and rule-plan results are logged by the coordinator that calls them.

## UI

Settings exposes:

- Debug logging toggle
- Open logs folder (`Env.OpenFileManager` on the log directory)

The Activity screen (ADR-0011 / RD-038) lists complete JSON records from this log file and its rotations in pages of 100 with a `before` cursor. Non-secret string attributes beyond `time`/`level`/`msg` are shown as `attrs` chips (`password`, `secret`, `token`, `profile` stay omitted). It is not a live tail and not a ProxyBridge connection journal.
