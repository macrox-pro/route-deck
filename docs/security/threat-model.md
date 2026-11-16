# Threat model (v1)

RouteDeck is a local macOS app. It does not run a server and does not sync to a cloud account.

## Assets

| Asset | Where |
| --- | --- |
| Proxy passwords | Keychain |
| Proxy host, port, username | SQLite |
| Profile and rules | SQLite and the JSON passed to ProxyBridge stdin |
| Process names and bundle ids in rules | SQLite and the profile |

## Abuse cases

| Case | Response |
| --- | --- |
| Password in SQLite, logs, or export | Forbidden. Serializer may hold the password in memory only while writing stdin |
| Replaced ProxyBridge binary | Prefer the documented install paths. A user path is allowed. Record version. Do not claim a signature check until a task implements one with a real API |
| Shell injection | `exec.Command` with argument list only. No shell string |
| Log injection | Do not log the profile document |
| Wrong process match | Warn when path or signature changes. Enforcement is `processNames`, compared to bundle id or executable name. Signature and path are not match fields |
| Missing proxy id | Reject the profile. CLI would otherwise send that rule direct |
| CLI crash or SIGKILL | Rules may remain or vanish. State becomes not protected. No kill-switch claim |
| GUI takes the extension | `Conflict`, stop retry loop |
| Damaged SQLite | Recovery in `docs/operations/recovery.md`. Do not show `Active` |
| Export | Default export has no passwords |

## Out of scope

- Network attacks on the VPN client itself
- macOS system-extension exploits inside ProxyBridge
- Multi-user Fast User Switching (not specified)
