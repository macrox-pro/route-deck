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
| RD-034 | Browse installed macOS applications (with icons) behind ports, and add a chosen app to the catalog without auto-import. |
| RD-035 | Keep RouteDeck running when the window closes: menu-bar tray, Open restores the window, every Quit path shows a Question dialog. |
| RD-036 | Tray Status and tooltip track ProtectionState (and BridgeState when not Ready); Stop enables only when RouteDeck owns the CLI; Quit dialog states that protection will stop when CLI is owned, then Stop and quit. |
