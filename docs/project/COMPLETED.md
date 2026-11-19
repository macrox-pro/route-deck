# Completed tasks

One-line archive of DONE `RD-*` work. Full task bodies lived under `docs/project/tasks/` and remain in git history.

Agents: do not load this file at session start. Open it only when the owner asks about past work.

When a task becomes DONE: append one row here, remove `docs/project/tasks/RD-XXX.md`, and drop the BACKLOG Done row (keep the Active index for non-DONE only).

| ID | Goal / outcome |
| --- | --- |
| RD-000 | Owner confirms `CONCEPTION.md` bounds for v1. |
| RD-001 | Record real Go, Wails, test, and build results for this scaffold. |
| RD-002 | Align module path, bundle metadata, and Dock icon with MACROX RouteDeck. |
| RD-003 | Owner accepts ADR-0001, ADR-0002, ADR-0004, and ADR-0005 (ADR-0003 is already the GORM decision to confirm in the same pass). Confirm ADR-0006 (vendor) in the same pass or note it already matches owner action. |
| RD-004 | Owner confirms aggregates and the three state dimensions. |
| RD-005 | Replace any TBD path list with facts from current ProxyBridge macOS docs. |
| RD-006 | When persistence work starts, verify `AutoMigrate` and `gorm.io/driver/sqlite` against the installed module docs. |
| RD-007 | Add CI workflow matching local checks. Do not require ESLint. |
| RD-008 | Add `internal/domain` value types and parsers with no aggregates and no protection derivation. |
| RD-009 | Add the three domain structs named in the domain model, with no password field and no UDPLimitation home. |
| RD-010 | Add Profile and ApplicationPolicy and reject unresolved proxy ids without rewriting actions to Direct. |
| RD-011 | Write the 14-step ProtectionState derivation into CONCEPTION, add two glossary sentences, and implement `DeriveProtection` with table tests. |
| RD-012 | Add only `ProxyStore.Get` and `SecretStore.Load` in `internal/application/port`, with no adapters and no third method. |
| RD-013 | Make `golangci-lint run ./...` catch unformatted Go and non-exhaustive enum switches, fix the RD-008 parsers that slipped through, and write the rules down. |
| RD-014 | Persist Application, ProcessIdentity, ProxyEndpoint (no password), and Profile in SQLite via GORM, with ports List/Save/Get and round-trip tests. |
| RD-015 | Add application services that save catalog aggregates and build an honest protection snapshot without ProxyBridge. |
| RD-016 | Remove GreetService and the time event; wire sqlite and SnapshotService so the app returns a snapshot DTO. |
| RD-017 | Ship a HashRouter shell with six screens; Dashboard loads the Go snapshot and shows protection with BadgeWithDot. |
| RD-018 | Wire structured `slog` logging to a rotating JSON file and stderr so every later adapter can leave an honest trail. |
| RD-019 | Extend domain policies and application ports so the slice can persist settings, delete with reference checks, inspect apps, and store secrets without inventing APIs later. |
| RD-020 | Persist policy scope, username, settings, and reject Delete of referenced proxies/apps with typed errors. |
| RD-021 | Implement SecretStore with `github.com/keybase/go-keychain` so proxy passwords leave SQLite. |
| RD-022 | Inspect a selected `.app` (and nested helpers) through a platform-neutral port so Catalog can add applications without embedding plist logic in services. |
| RD-023 | Discover ProxyBridge (custom path then documented paths), parse `--version`, and feed Snapshotter so the snapshot can leave permanent Unknown. |
| RD-024 | Probe proxy TCP/SOCKS5/UDP ASSOCIATE health and keep a registry refreshed on triggers and a 30s ticker. |
| RD-025 | Compile a profile into a wire-neutral rule plan and AppCoverage facts using the owner protocol-scope mapping. |
| RD-026 | Serialize a domain rule plan into deterministic macOS `.pbprofile` JSON, loading passwords only inside the serializer. |
| RD-027 | Manage one ProxyBridge CLI process with stdin profile, graceful stop, and a versioned line classifier with fixtures from real 4.1.3. |
| RD-028 | Orchestrate Apply/Stop with confirm, one rollback, Conflict handling, and honest ProtectionFacts — all tested with fake ports. |
| RD-029 | Expose Apply/Stop and an extended snapshot over Wails, emit product events, stop CLI on shutdown, and re-check on wake/network change. |
| RD-030 | Expose Proxies, Applications, Profiles, and Settings over Wails with typed errors, including debug logging and open-logs-folder. |
| RD-031 | Ship Untitled UI forms for proxies, application picker with helper checklist, and Settings (path, debug logging, open logs folder). |
| RD-032 | Edit multi-profile policies and Apply/Stop from Dashboard with honest banners and event-driven snapshot refresh. |
| RD-033 | Prove on macOS with ProxyBridge 4.1.3 that Apply reaches Active; Stop → Unprotected; logs without secrets (Shadowrocket 127.0.0.1:1082). |
| RD-034 | Browse installed macOS applications (with icons) behind ports, and add a chosen app to the catalog without auto-import. |
| RD-035 | Keep RouteDeck running when the window closes: menu-bar tray, Open restores the window, every Quit path shows a Question dialog. |
| RD-036 | Tray Status and tooltip track ProtectionState (and BridgeState when not Ready); Stop enables only when RouteDeck owns the CLI; Quit dialog states that protection will stop when CLI is owned, then Stop and quit. |
| RD-037 | Settings Launch at Login via Wails Autostart without auto-Apply. |
| RD-038 | Show recent RouteDeck JSON log records on the Activity screen without passwords or a connection journal. |
| RD-039 | Probe HTTP proxies with CONNECT and report AuthenticationFailed on 407. |
| RD-040 | Ship Setup guides screen with VPN-client detection and Port open facts (ADR-0013). |
| RD-041 | One-click Add this proxy from a guide (save + Recheck); form fallback without default port. |
| RD-042 | Suggestion banner on Dashboard/Proxies for detected open local ports. |
| RD-043 | Dashboard first-run checklist that auto-hides when complete. |
| RD-044 | Distinguish WARN and ERROR in RouteDeck logs so Activity severity matches. |
| RD-045 | Native dialogs for ProxyBridge GUI conflict and extension approval; remove Dashboard Conflict banner. |
| RD-046 | Honest ProtectionState via snapshot Apps; Applied vs Protection UX; proxy/profile icons; health labels; helpers H2. |
| RD-047 | Profiles form single-column layout; add-all processes by application (skip duplicates silently). |
| RD-048 | Profiles app-grouped editor (icon/Action/helpers + Apply to helpers); Review helpers preselect. |
| RD-049 | ADR-0014 and desktop UX docs/contracts. |
| RD-050 | Nested helper discovery (Resources, Current, IsUpdater) + ChatGPT fixture tests. |
| RD-051 | HelperGroupDTO, sync ConfirmHelpers, rule collapse, ErrConflictingMatch, catalog.changed. |
| RD-052 | EditorDialogs host + dialog shell routes; tray focuses open editor. |
| RD-053 | Profile editor window with HelperGroupDTO; profiles page list-only with delete confirm. |
| RD-054 | Proxy editor window; proxies page list-only; guide prefill; password never pre-filled. |
| RD-055 | Helpers dialog from HelperGroupDTO; auto-open after Add; Applications uses pendingCount. |
| RD-056 | Dock attention badge `!` when applied protection is not Active; vendor dock service. |
| RD-057 | CI macOS package workflow (`darwin:package` + zip artifact, no notarization). |
| RD-058 | Link-time app version/commit for Settings About and sidebar badge. |
| RD-059 | Persist Appearance in SQLite and sync native window chrome via recreate. |
| RD-060 | Opt-out Auto Apply on launch (default ON) via ADR-0015 gates + existing Apply path. |
| RD-061 | First-run checklist: partial proxy step (UDP-limited/Checking) + session dismiss with regression. |
