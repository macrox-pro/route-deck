# ProxyBridge profile contract

Sources checked 2026-11-16:

- https://proxybridge.interceptsuite.com/docs/macos-profiles/
- https://proxybridge.interceptsuite.com/docs/macos-cli/
- Rules detail: https://proxybridge.interceptsuite.com/docs/macos-rules/ (re-read in RD-005 before coding rule fields)

## Format

macOS profile is JSON. Required lists: `proxyConfigs`, `proxyRules`. Windows/Linux profiles use different key casing and MUST be rejected.

Limits enforced by the CLI before start (from the CLI and profiles docs):

| Limit | Value |
| --- | --- |
| Proxies | 100 |
| Rules | 5000 |
| Stdin size | 4 MiB |
| Port | 1–65535 |
| Proxy type | `socks5` or `http` |
| Proxy id | Unique in the file |
| Protocol | `TCP`, `UDP`, or `BOTH` |
| SOCKS username/password | Each at most 255 bytes |

## Proxy object

| Field | Notes |
| --- | --- |
| `id` | Unique string. Rules point at it |
| `name` | Optional label |
| `type` | `socks5` or `http` |
| `host` | Required |
| `port` | Number |
| `username`, `password` | Optional. Clear text in the JSON the CLI reads |
| `sendDomain` | Optional. Docs: missing means true |

RouteDeck loads the password from Keychain only inside the serializer. Domain storage keeps a Keychain reference. Logs MUST NOT contain the profile body.

## Rule object

| Field | Notes |
| --- | --- |
| `name` | Optional |
| `processNames` | Match key. Not Bundle ID |
| `targetHosts` | Host patterns. Weak if DNS is off |
| `targetPorts` | Port patterns. Docs show lists such as `80; 443` or `*` |
| `protocol` | `TCP`, `UDP`, or `BOTH` (missing means `BOTH` on import) |
| `action` | `DIRECT`, `BLOCK`, or a proxy `id` |
| `enabled` | Optional. Missing means true |

## Fail-open (CLI behavior)

The CLI warns and still runs when a rule action is not `DIRECT`, `BLOCK`, or a known proxy id. Traffic for that rule goes direct (`rule N sends traffic to a proxy that is not in this profile, it will go direct`).

RouteDeck MUST reject the profile before start if any proxy id does not resolve. MUST NOT rely on that warning.

No matching rule means the connection is not redirected. Absence of a rule is not `Block`.

## UDP

HTTP proxies do not carry UDP. For a rule that needs UDP and the chosen proxy is HTTP, compile UDP as `BLOCK` unless the user has stored an explicit Direct-UDP limitation. Do not mark that case protected.

## Profile settings RouteDeck sets

| Field | RouteDeck v1 |
| --- | --- |
| `dnsEnabled` | true unless a task says otherwise. Domain rules are secondary to process match |
| `trafficLoggingEnabled` | false unless a task needs fixtures. v1 Activity does not require CLI connection lines |
| `closeToMenuBar`, `runAtStartup` | false. The CLI ignores them. RouteDeck tray is separate |

## Internal vs wire format

Domain types never match this JSON. The ProxyBridge adapter owns the DTO and serializer. Output MUST be deterministic for the same input (stable ordering).

## Apply

The CLI reads stdin at start. A new profile means stop and start. See `docs/architecture/lifecycle.md` for the one rollback attempt.
