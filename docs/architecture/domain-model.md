# Domain model

Types below are conceptual. Implement them only in a READY task. Names are normative.

## Aggregates and values

| Name | Role |
| --- | --- |
| Application | User-facing app entry (display name, bundle id, icon ref) |
| ProcessIdentity | Bundle id, executable name, path, and signing hints. Bundle id and executable name may be written into `processNames` |
| ProxyEndpoint | Global proxy: type `socks5` or `http`, host, port, optional username (SQLite), Keychain ref for the password, `sendDomain` |
| Profile | Named set of rules. References proxy IDs. At most one profile is applied. Profile names are unique |
| ApplicationPolicy | Policy for one process identity inside a profile: `Proxy`, `Direct`, or `Block`, plus `ProtocolScope` and `UDPLimitation`. A Proxy action requires a non-empty scope |
| RoutingAction | `PROXY` (proxy id), `DIRECT`, or `BLOCK` |
| ProtocolScope | `TCP`, `UDP`, or `BOTH` |
| ProtectionSnapshot | Point-in-time BridgeState, per-proxy ProxyHealth, aggregate ProtectionState, coverage gaps |
| UDPLimitation | Absent, or explicit user accept of Direct-UDP. Default compile is UDP `BLOCK` when the proxy cannot carry UDP |

## Matching

ProxyBridge rules use the JSON field `processNames`. The CLI compares each pattern with the app's bundle identifier and with its executable name. Either match is enough. Match is case-insensitive and whole-string; `*` is the only wildcard; a pattern containing `/` does not match.

Code signature and path stay UX warnings. They are not CLI match fields.

Helper processes are not added to a rule until the user confirms them. Catalog stores at most one `ProcessIdentity` per bundle id per application (ADR-0014). Confirm sync removes unchecked groups; referenced helpers return `ErrReferenced`.

When several stored identities share a bundle id and the same action, `BuildRulePlan` emits one rule whose `processNames` is only that bundle id. Different actions on a shared matchable name are resolved **last-wins** (last policy in the list) before Save and Apply (ADR-0019); they are not rewritten to Direct. `ValidateMatchConflicts` remains for asserts after resolve.

## Invariants

- A profile MUST NOT be applied if any rule references a missing proxy ID.
- Explicit `Direct` is a user policy. It is not protection.
- Loss of a proxy MUST NOT rewrite `Proxy` rules to `Direct`.
- `ProtectionState.Active` requires bridge ready, profile confirmed, refs valid, and every protected app covered for required protocols.
- Desired and applied may differ. UI shows both when they differ.
- Passwords are not fields on the domain aggregate that get persisted. Store a Keychain reference.
- Deleting a proxy, application, or helper identity that any profile references returns `ErrReferenced`.
- Duplicate application bundle ids return `ErrDuplicateBundleID`. Duplicate profile names return `ErrDuplicateProfileName`.
- Conflicting actions on identities that share a matchable name are resolved last-wins on Save/Apply (ADR-0019); they MUST NOT be rewritten to Direct.
- Protocol outside a Proxy policy's scope compiles to `BLOCK` (covered). HTTP+UDP uses `BLOCK` or `DIRECT` per `UDPLimitation`; UDP is not covered either way (ADR-0008).

## Apply failure

If a new CLI start does not confirm the profile:

1. Desired stays the new profile.
2. Applied becomes none and ProtectionState is `Unprotected`.
3. One restart of the previous valid applied profile is attempted.
4. If that fails, stay `Unprotected` and return both errors.

## Mapping sketch to ProxyBridge

| RouteDeck | ProxyBridge profile |
| --- | --- |
| ProxyEndpoint | `proxyConfigs[]` (`id`, `name`, `type`, `host`, `port`, `username`, password at serialize time, `sendDomain`) |
| Policy + ProcessIdentity | `proxyRules[]` (`processNames`, `protocol`, `action`, `enabled`) |
| Proxy action | action = proxy id |
| Direct | action = `DIRECT` |
| Block | action = `BLOCK` |

Serializer and ProxyBridge DTO live in the ProxyBridge adapter, not in domain.
