# ADR-0008: Vertical slice product decisions

Status: Accepted
Date: 2026-11-17
Owners: MACROX RouteDeck

## Context

Phase 3 left a launchable shell with an honest Unknown snapshot. Phase 4 and Phase 5 need closed product answers before coding: secrets, app intake, helpers, health depth, rule compilation, profiles, apply lifecycle, logging, and failure handling. The owner closed those questions on 2026-11-17.

## Decision

### Scope

Ship a full vertical slice including UI: add a proxy and an application, edit multiple profiles, Apply and Stop through the official ProxyBridge CLI, and show a confirmed ProtectionState. Activity, tray, autostart, export, CI, and branding recolor stay out of this batch.

### Secrets

`github.com/keybase/go-keychain` (cgo, Security.framework) implements the `SecretStore` port. Keychain service is `pro.macrox.routedeck.proxy`; account is the proxy id. SQLite stores the reference only. Proxy username lives in SQLite. An empty password field on edit means keep the stored password; clearing the password is a separate explicit action.

### Application intake and helpers

The Wails open-file dialog selects a `.app`. `howett.net/plist` reads `Info.plist` behind a platform-neutral `ApplicationInspector` port; only a macOS adapter ships in this batch. Nested helper `.app` bundles (for example under `Contents/Frameworks` and `Contents/Helpers`) are listed with checkboxes. Confirmed helpers become separate `ProcessIdentity` rows under the same `Application`. Unconfirmed helpers never enter rules. A duplicate application bundle id is rejected.

### Health

Every proxy gets a TCP connect. SOCKS5 also runs the greeting, RFC1929 auth (`AuthenticationFailed` on failure), and a UDP ASSOCIATE request (`Available` on success, `TCPAvailableUDPUnavailable` on failure). HTTP proxies use TCP connect only and are never reported `AuthenticationFailed` in this batch. Triggers: startup, proxy save, before Apply, Recheck, and every 30 seconds. Wake (`SystemDidWake`) and network change (`NetworkChanged`) set `RecheckPending` until a fresh check finishes.

### Rule compilation

Each Proxy policy has a protocol scope (`TCP`, `UDP`, or `BOTH`) and a Direct-UDP limitation flag. A protocol outside the scope compiles to `BLOCK` (covered). For an HTTP proxy with UDP in scope, UDP compiles to `BLOCK` by default or `DIRECT` when the Direct-UDP limitation is set; UDP is not covered either way. `targetHosts` and `targetPorts` stay `*`. Domain rules are out of scope.

### Profiles and catalog

Multiple profiles: create, copy, rename, delete, set desired. Delete is refused for the desired or applied profile. Profile names are unique. Deleting a proxy or application that any profile references is rejected with a typed error.

### Apply lifecycle

No auto-apply at startup. Editing while running does not re-apply; the snapshot marks desired-differs-from-applied until Apply. Quit stops the CLI gracefully on `OnShutdown` with no dialog. An unexpected CLI exit (not GUI conflict) does not auto-restart. A foreign or orphan CLI maps to `Conflict` with no signal to processes RouteDeck did not start. A rejected-rules warning, even before redirect-ready, counts as unconfirmed apply and triggers the single rollback.

### ProxyBridge path and version

A custom absolute path in Settings is stored in SQLite and checked before the documented paths. Documented install locations are owned by the `BridgeLocator` adapter (`PlatformPaths`); macOS ships `/usr/local/bin/proxybridge` and the app-bundle executable, other OS adapters supply their own list. Settings exposes `discoveredProxyBridgePath` from `Locate("")` as the path input placeholder. A version is known when `--version` parses. Tested version is `4.1.3`; any other parsed version runs with a compatibility warning.

### Logging

Standard `log/slog` injected into services and adapters. Destinations: JSON lines in `~/Library/Logs/route-deck/route-deck.log` and text to stderr, combined with `slog.NewMultiHandler`. Own size-based rotation: 5 MB, keep 3. Info default; Debug toggle in Settings applied live via `slog.LevelVar`. Every classified CLI line is logged at Info (truncated). Settings offers "Open logs folder". No custom logger port. MUST NOT log passwords, tokens, Keychain data, profile JSON, or traffic payloads. Details: `docs/engineering/logging.md`.

## Alternatives

- Keychain via `/usr/bin/security` or `99designs/keyring`. Rejected for go-keychain behind the existing port.
- Scan `/Applications` or manual-only intake. Rejected for the native picker plus helper confirmation in this ADR. Superseded for browse listing by [ADR-0009](ADR-0009-installed-application-browse.md) (catalog still requires Add).
- Single Default profile. Rejected for multi-profile operations in the slice.
- Auto-apply or auto-reapply. Rejected for manual Apply with an honest Unprotected or desired-differs banner.
- Zap/zerolog or a custom Logger port. Rejected for stdlib `slog` with replaceable handlers.

## Consequences

- Platform-specific work stays in adapters (`keychain`, `macos`, `proxybridge`, `proxyhealth`, `logfile`).
- Second cgo dependency (go-keychain) alongside sqlite.
- Unsigned rebuilds may prompt for Keychain access; recorded as a known limitation.
- HTTP proxy auth is not verified in this batch.

## Follow-up

Tray quit warning, Activity log viewer, HTTP auth verification, and CI (RD-007) remain later tasks.

## Open decisions (owner)

None.
