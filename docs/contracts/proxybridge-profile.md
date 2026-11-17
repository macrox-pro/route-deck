# ProxyBridge profile contract

Sources checked 2026-11-16:

- https://proxybridge.interceptsuite.com/docs/macos-profiles/
- https://proxybridge.interceptsuite.com/docs/macos-cli/
- https://proxybridge.interceptsuite.com/docs/macos-rules/
- https://proxybridge.interceptsuite.com/docs/macos-how-it-works/

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
| `processNames` | Patterns compared to bundle identifier and executable name. Either match is enough. Case-insensitive, whole-string. `*` is the only wildcard. A pattern with `/` does not match. Separators: `;` or comma |
| `targetHosts` | Host patterns. Separators: `;` only (a comma is not a separator). No CIDR. Weak if DNS is off or the app uses DNS-over-HTTPS |
| `targetPorts` | Port patterns. Separators: `;` or comma. Docs show lists such as `80; 443` or `*` |
| `protocol` | `TCP`, `UDP`, or `BOTH` (missing means `BOTH` on import) |
| `action` | `DIRECT`, `BLOCK`, or a proxy `id` |
| `enabled` | Optional. Missing means true |

## Match semantics

- The JSON field is still `processNames`. The CLI matches each pattern against the program's bundle identifier and its executable name.
- `*` or empty matches every program.
- No matching rule means the connection is not redirected. Absence of a rule is not `Block`. RouteDeck MUST NOT report that as protected.
- ProxyBridge's own app and extension flows are never redirected.
- Domain entries in `targetHosts` need Enable DNS. With DNS off, or when the app uses DNS-over-HTTPS / DNS-over-TLS, domain patterns do not match. IP patterns and `processNames` matching still apply.
- Among enabled rules: a specific program match beats a catch-all. A catch-all is `*` with no host limit and no port limit. Only the first catch-all applies. A `*` rule that has a host or port limit is not a catch-all; on TCP it matches immediately when the limit fits. Among specific rules (and limited `*` rules), earlier list order wins.

## Fail-open (CLI behavior)

The CLI warns and still runs when a rule action is not `DIRECT`, `BLOCK`, or a known proxy id. Traffic for that rule goes direct (`rule N sends traffic to a proxy that is not in this profile, it will go direct`).

RouteDeck MUST reject the profile before start if any proxy id does not resolve. MUST NOT rely on that warning.

No matching rule means the connection is not redirected. Absence of a rule is not `Block`.

## UDP

For a specific program, UDP matches on program and protocol only. Target hosts and ports are ignored for that UDP.

A rule whose program is `*` and which also has a host or port limit never matches UDP. Do not treat that rule as UDP coverage.

HTTP proxies do not carry UDP. The CLI sends that UDP direct and may log a warning. For a rule that needs UDP and the chosen proxy is HTTP, RouteDeck compiles UDP as `BLOCK` unless the user has stored an explicit Direct-UDP limitation. Do not treat the CLI's direct UDP as protection. Do not mark that case protected.

## TCP proxy failure

TCP sent to a proxy that is down, rejects, or is silent is closed. The CLI does not rewrite that traffic to Direct.

## Profile settings RouteDeck sets

| Field | RouteDeck v1 |
| --- | --- |
| `dnsEnabled` | true unless a task says otherwise. Domain rules are secondary to `processNames` match |
| `trafficLoggingEnabled` | false unless a task needs fixtures. v1 Activity does not require CLI connection lines |
| `closeToMenuBar`, `runAtStartup` | false. The CLI ignores them. RouteDeck tray is separate |

## RouteDeck protocol-scope mapping (ADR-0008)

Each application policy has a protocol scope (`TCP`, `UDP`, or `BOTH`) and an optional Direct-UDP limitation.

| Situation | Compiled rule(s) | Coverage |
| --- | --- | --- |
| Proxy / Direct / Block within scope | One rule with that action and protocol | Per CONCEPTION coverage facts |
| Protocol outside the selected scope | Separate rule with action `BLOCK` for that protocol | Covered (Block) |
| HTTP proxy with UDP in scope, limitation Absent | UDP rule action `BLOCK` | UDP not covered |
| HTTP proxy with UDP in scope, limitation DirectUDP | UDP rule action `DIRECT` | UDP not covered |

`targetHosts` and `targetPorts` stay `*` in the vertical slice. Domain rules are out of scope.

## Internal vs wire format

Domain types never match this JSON. The ProxyBridge adapter owns the DTO and serializer. Output MUST be deterministic for the same input (stable ordering).

## Apply

The CLI reads stdin at start. A new profile means stop and start. See `docs/architecture/lifecycle.md` for the one rollback attempt.
