# CONCEPTION — MACROX RouteDeck

## Problem

Users run a VPN client that exposes a local SOCKS5 or HTTP CONNECT proxy. They need chosen macOS apps to use that proxy, to go direct, or to be blocked, without RouteDeck implementing a VPN.

## What RouteDeck is

MACROX RouteDeck is a macOS desktop app. It assigns each selected application one policy and applies that policy by driving the official ProxyBridge CLI.

| Policy | Behavior |
| --- | --- |
| `Proxy` | Send matching connections through a chosen SOCKS5 or HTTP CONNECT proxy |
| `Direct` | Do not redirect matching connections (explicit user choice) |
| `Block` | Block matching connections |

RouteDeck does not implement a VPN protocol, does not pick a `utun` interface, and does not ship its own Network Extension.

## Identity

| Item | Value |
| --- | --- |
| Product | MACROX RouteDeck |
| Binary | `route-deck` |
| Bundle ID (target) | `pro.macrox.routedeck` |
| Module (target) | `github.com/macrox-pro/route-deck` |
| Task IDs | `RD-*` |
| UI language | English |
| Docs | English |
| Platform v1 | macOS 13+ desktop only |

Brand mark: `docs/branding/route-deck-icon.png` (layered glass deck, white split-route arrow, blue / cyan / magenta).

## Happy path

1. The user installs a VPN client that provides a local SOCKS5 (or HTTP CONNECT) endpoint.
2. The user installs and approves official ProxyBridge.
3. RouteDeck finds the ProxyBridge binary and checks its version.
4. The user adds a proxy and selects applications.
5. Go builds a ProxyBridge profile, validates it, and starts the CLI with the profile on stdin.
6. The user sees routing state, proxy health, and RouteDeck events.
7. Setting changes produce a new profile; RouteDeck confirms the result before showing success.
8. If the proxy is lost, the UI shows degradation. RouteDeck MUST NOT switch a protected app to `Direct` automatically.

## Non-goals (v1)

- Own VPN tunnel or Network Extension
- Managing arbitrary network interfaces
- Changing settings inside third-party VPN clients
- Cloud sync, accounts, or a RouteDeck server
- Windows, Linux, iOS, Android as product targets (Wails scaffold folders are not a promise)
- Rebuilding or forking ProxyBridge
- A guaranteed kill switch after ProxyBridge crashes

## Architecture

```text
Wails v3 / React
       |
Wails inbound services (DTO, events)
       |
Application services (use cases)
       |
Domain (entities, invariants)
       |
Ports
       |
Adapters: ProxyBridge CLI, Keychain, GORM+SQLite, macOS discovery, process
```

| Layer | MUST | MUST NOT |
| --- | --- | --- |
| Domain | Entities, invariants, typed errors | SQL, Wails, `os/exec`, ProxyBridge JSON |
| Application | Use cases, transactions, ports | Know `.pbprofile` field layout |
| Adapters | Implement ports | Own product policy |
| React | Render snapshot and send commands | Contain routing rules |

Go is the source of truth for user configuration. Wails events are notifications. After a window opens, the UI MUST request a snapshot.

## Data rules

| Topic | Decision |
| --- | --- |
| Proxies | Global catalog; profiles reference proxy IDs |
| Applied profile | At most one |
| Enforcement key | ProxyBridge `processNames` (each pattern matches bundle identifier or executable name) |
| Bundle ID / signature | Bundle id may be written into `processNames`. Signature and path are UX warnings only; they are not CLI match fields |
| Passwords | macOS Keychain. SQLite stores a reference only |
| Persistence | GORM + SQLite + `gorm.io/driver/sqlite` + `AutoMigrate` |
| Activity v1 | RouteDeck events and health/ops, not an assumed per-connection CLI log |

## States

Three dimensions. Do not collapse them into one enum.

**BridgeState:** `NotInstalled`, `Detected`, `Starting`, `Ready`, `Stopping`, `Stopped`, `Failed`, `Conflict`, `Unknown`.

**ProxyHealth:** `Unknown`, `Checking`, `Available`, `Unavailable`, `AuthenticationFailed`, `TCPAvailableUDPUnavailable`.

**ProtectionState:** `Inactive`, `Initializing`, `Active`, `Degraded`, `Unprotected`, `Failed`, `Unknown`.

`Active` is allowed only when the CLI is running, the profile is confirmed, proxy references resolve, every protected process is covered, and required protocols are actually provided. `Active` means routing is applied. It does not mean crash-proof isolation.

Dashboard aggregate: `Active` only if the bridge is ready AND every protected app is fully covered for its required protocols. Any gap is `Degraded`, `Unprotected`, or `Unknown` per the matrix. MUST NOT average mixed apps into a fake `Active`.

### Failure matrix

| Condition | Required outcome |
| --- | --- |
| ProxyBridge missing | Onboarding, not protected |
| Version unknown | Compatibility warning, not silent success |
| ProxyBridge GUI took the extension | `Conflict` and not protected |
| CLI fails to start | `Failed` |
| CLI exits unexpectedly | `Unprotected` |
| Proxy unreachable | `Unavailable`; no automatic `Direct` |
| Auth failure | `AuthenticationFailed` |
| SOCKS5 has no UDP | `TCPAvailableUDPUnavailable`; UDP not claimed protected |
| HTTP proxy used where UDP is required | Not protected for UDP |
| Proxy ID missing in profile | Reject apply |
| No rule for a process | That app is not protected |
| DNS observation incomplete | Domain rules `Degraded`; process match remains primary |
| Window closed | App keeps running (tray phase) |
| Quit | Warn and stop CLI under RouteDeck control |
| Database damaged | Diagnose / recover; do not pretend success |
| Wake or network change | Re-check health; do not assume still protected |

### ProtectionState derivation

First match wins.

1. Wake or network change, re-check not finished: `Unknown`.
2. Bridge `Unknown`, or CLI version unknown, or a referenced proxy is `Unknown` or `Checking`: `Unknown`.
3. Bridge `Starting` or `Stopping`: `Initializing`.
4. Confirm failed and the one rollback failed, or applied is none after a failed confirm: `Unprotected`, even if the bridge is `Failed`.
5. Bridge `Failed` on start, and step 4 did not apply: `Failed`. An unexpected CLI exit is step 8 or 10, not step 5. The matrix names that exit `Unprotected`.
6. Bridge `Conflict`: `Unprotected`.
7. Bridge `NotInstalled` or `Detected`: `Inactive`, with or without a desired profile. Owner.
8. Bridge `Stopped` and a profile is desired: `Unprotected`. Owner.
9. Bridge `Stopped` and no profile is desired: `Inactive`. Owner's split of step 8. No third test.
10. Bridge `Ready` and the applied profile is not confirmed: `Unprotected`.
11. Bridge `Ready`, profile confirmed, protected set empty: `Inactive`. Owner. Empty coverage is not `Active`.
12. Bridge `Ready`, profile confirmed, every protected app covered for its required protocols, every proxy used by a `Proxy` rule is `Available`: `Active`.
13. Bridge `Ready`, profile confirmed, at least one protected app covered and at least one gap: `Degraded`.
14. Bridge `Ready`, profile confirmed, protected set non-empty, nothing covered: `Unprotected`.

Coverage facts the caller may set:

- User `Block` for a protocol: that protocol's covered fact is true, including UDP. Proxy UDP limits do not clear it.
- UDP required and the proxy cannot carry it: covered fact is false. `Degraded` when another protocol for that app is covered (step 13).
- Direct-UDP limitation on required UDP: covered fact is false. `Direct` is not protection.
- Proxy `Unavailable` or `AuthenticationFailed` while steps 1–11 did not match: `Degraded`. Owner said this for `Unavailable`. Auth failure stays a proxy failure with no rewrite to `Direct`, which is the same step-13 gap.
- `TCPAvailableUDPUnavailable`: TCP may be covered. Required UDP is not.
- Domain rules present and DNS incomplete: gap, so not step 12. Process coverage alone does not clear it.
- Missing proxy id: not a protection word. `ValidateProxyRefs` returns `ErrUnresolvedProxy`. Snapshot word stays on the last confirmed profile, or step 10 when nothing is confirmed.
- Database damaged: input flag forces step 2's `Unknown`. The matrix does not name `Failed` here.
- Quit uses the `Stopped` steps. Window close changes no protection input.

## UDP and HTTP

- TCP may use SOCKS5 or HTTP CONNECT.
- UDP needs SOCKS5 UDP ASSOCIATE.
- HTTP CONNECT does not carry UDP. MUST NOT mark that path protected for UDP.
- Default compile: UDP for an HTTP proxy is `BLOCK`, unless the user records an explicit accept-Direct-UDP limitation.
- UDP rules match application and protocol. Do not promise address/port filtering for UDP the way TCP rules do.
- v1 primary match is `processNames` (bundle id or executable name), not domain. DNS-over-HTTPS can make domain rules incomplete.

## Apply and rollback

The CLI reads the profile at start. A new profile requires a restart. There is a window with no rules. MUST NOT promise a seamless switch.

If the new CLI does not confirm the profile:

1. Desired config stays the new one.
2. Applied state becomes none / `Unprotected`.
3. Attempt **one** restart of the previous valid applied profile.
4. If that restart fails, stay `Unprotected` and surface both errors.

## Security invariants

- An unresolved proxy ID MUST reject the profile. ProxyBridge can fail open to Direct on a missing proxy ID.
- A dead proxy MUST NOT be rewritten to `Direct`.
- `Unknown` MUST NOT be shown as protected.
- Do not promise zero leaks while ProxyBridge can drop rules on crash.
- Do not log passwords, tokens, full profiles that contain secrets, or user traffic payloads.

## Not guaranteed

- Kill switch after ProxyBridge crash or kill
- Seamless profile swap
- UDP safety on HTTP proxies
- Domain match when DNS is encrypted
- Protection of processes that have no rule
- Per-connection logs from CLI stdout
