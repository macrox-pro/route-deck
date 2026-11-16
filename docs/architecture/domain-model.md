# Domain model

Types below are conceptual. Implement them only in a READY task. Names are normative.

## Aggregates and values

| Name | Role |
| --- | --- |
| Application | User-facing app entry (display name, bundle id, icon ref) |
| ProcessIdentity | Path, signing hints, and the process name string used for rules |
| ProxyEndpoint | Global proxy: type `socks5` or `http`, host, port, Keychain ref for the password, `sendDomain` |
| Profile | Named set of rules. References proxy IDs. At most one profile is applied |
| ApplicationPolicy | Policy for one process identity inside a profile: `Proxy`, `Direct`, or `Block` |
| RoutingAction | `PROXY` (proxy id), `DIRECT`, or `BLOCK` |
| ProtocolScope | `TCP`, `UDP`, or `BOTH` |
| ProtectionSnapshot | Point-in-time BridgeState, per-proxy ProxyHealth, aggregate ProtectionState, coverage gaps |
| UDPLimitation | Absent, or explicit user accept of Direct-UDP. Default compile is UDP `BLOCK` when the proxy cannot carry UDP |

## Matching

ProxyBridge rules use `processNames`. Bundle ID and code signature help the user recognize an app and warn on path or signature change. They are not the runtime match key.

Helper processes are not added to a rule until the user confirms them.

## Invariants

- A profile MUST NOT be applied if any rule references a missing proxy ID.
- Explicit `Direct` is a user policy. It is not protection.
- Loss of a proxy MUST NOT rewrite `Proxy` rules to `Direct`.
- `ProtectionState.Active` requires bridge ready, profile confirmed, refs valid, and every protected app covered for required protocols.
- Desired and applied may differ. UI shows both when they differ.
- Passwords are not fields on the domain aggregate that get persisted. Store a Keychain reference.

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
