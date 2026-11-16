# Glossary

| Term | Meaning |
| --- | --- |
| RouteDeck | This product. Assigns per-app policy and drives ProxyBridge. |
| ProxyBridge | Third-party macOS tool. Official CLI is the only network enforcer. |
| Policy | `Proxy`, `Direct`, or `Block` for one application identity. |
| Proxy | A global SOCKS5 or HTTP CONNECT endpoint. Profiles reference its ID. |
| Profile | A named set of rules plus proxy references. At most one is applied. |
| Desired | Config the user asked to apply. |
| Applied | Config the running CLI has confirmed. |
| ProcessIdentity | How RouteDeck remembers an app (bundle id, executable name, path, signature hints). |
| processNames | ProxyBridge rule field. Each pattern matches the bundle identifier or the executable name. |
| BridgeState | Lifecycle of the ProxyBridge CLI integration. |
| ProxyHealth | Reachability and protocol support of one proxy. |
| ProtectionState | Whether routing for the protected set is confirmed. |
| Snapshot | Backend state the UI loads. Events do not replace it. |
| Keychain reference | Pointer stored in SQLite. The password is not in SQLite. |
| Fail-open | ProxyBridge may send traffic direct if a proxy ID is missing or no rule matches. RouteDeck must not allow the missing-ID case. |
| Degraded | Protection is partial or a dependency is weak. Not `Active`. |
| Unprotected | Rules are not in effect. |
| Direct-UDP limitation | User-recorded choice to allow UDP to go direct on a non-UDP proxy. Default is `BLOCK` for that UDP. |

## Not the same

| A | B |
| --- | --- |
| Product bundle id (`pro.macrox.routedeck`) | A value written into a rule's `processNames` field |
| Path or code signature | CLI match fields (`processNames`) |
| Wails event | Source of truth |
| Open TCP port | Healthy proxy |
| `Active` | Crash-proof kill switch |
| Explicit `Direct` policy | Silent fallback to Direct |
