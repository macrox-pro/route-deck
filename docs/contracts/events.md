# Events contract

Wails events notify the UI. They are not the source of truth. On window open, the frontend calls the snapshot API.

Each product event payload MUST include:

| Field | Meaning |
| --- | --- |
| `schema` | Integer schema version, start at 1 |
| `timestamp` | RFC3339 UTC |
| `payload` | Event-specific object |

## Names

| Event | When |
| --- | --- |
| `protection.state.changed` | ProtectionState or coverage changed |
| `proxy.health.changed` | One proxy health changed |
| `proxybridge.state.changed` | BridgeState changed |
| `application.discovered` | Discovery found an app |
| `operation.progress` | Long operation tick (`operation_id`) |
| `operation.failed` | Long operation failed (`operation_id`, typed error) |
| `editors.save_for_quit` | Quit Save and Quit asked dirty profile/proxy editors to run Save (ADR-0018; `kinds`) |
| `main.command` | Menu bar, tray utility actions, or `UtilityWindowsService.ShowMain` asked the main window to run a command (ADR-0020). Payload `command` is closed: `navigate` (optional `route`), `apply`, `stop`, `recheck`, `toggle-sidebar`, `add-application`. The event is a trigger; `app-shell` is the sole listener; pages reload snapshots after handlers run. MUST NOT carry protection state |

`connection.observed` is not a v1 event. Activity reads the RouteDeck JSON log (ADR-0011), not a stored copy of these Wails events and not a parsed CLI connection stream. Product events may trigger an Activity reload; they are not proof of `Active`.

The scaffold event `time` was removed with RD-016. Product events are emitted from `main.go` / Protection and Proxies Wails services (RD-029+) with `EventEnvelope{schema,timestamp,payload}`.

## Rules

- Re-subscribe after reload.
- Ignore unknown `schema` versions (log and keep the last snapshot).
- MUST NOT treat an event as proof of `Active` without a snapshot that says so.
- Register payload types with Wails only using APIs that exist in v3.0.0. The scaffold uses `application.RegisterEvent[string]("time")`. Copy that pattern only after checking it still exists in the installed module.
