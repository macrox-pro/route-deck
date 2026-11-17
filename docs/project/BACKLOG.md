# Backlog

Statuses: `TODO`, `READY`, `IN_PROGRESS`, `BLOCKED`, `REVIEW`, `DONE`, `DROPPED`.

Agents MUST NOT `git commit`. The owner commits after review.

Feature-drift rule: `STATUS.md` current task, this index, and `HANDOFF.md` task ID MUST match.

## How to read (context budget)

1. Use this file as an **index only**. Do not paste full task bodies here.
2. Open **one** task file: `docs/project/tasks/RD-XXX.md` matching `STATUS.md`.
3. MUST NOT `@` every file under `docs/project/tasks/` in one session.
4. DONE tasks stay as files for history; do not load them unless the owner asks.

## Active

| ID | Status | Goal |
| --- | --- | --- |
| [`RD-037`](tasks/RD-037.md) | REVIEW | Settings Launch at Login via Wails Autostart without auto-Apply. |
| [`RD-033`](tasks/RD-033.md) | REVIEW | Prove on macOS with real ProxyBridge 4.1.3 that Apply reaches Active, failures stay honest, and the log file traces each scenario without secrets. |
| [`RD-007`](tasks/RD-007.md) | TODO | Define CI commands that match local checks. Do not require ESLint. |

## Done

| ID | Status | Goal |
| --- | --- | --- |
| [`RD-036`](tasks/RD-036.md) | DONE | Tray Status/tooltip track protection; Stop when CLI owned; Quit dialog copy then Stop. |
| [`RD-035`](tasks/RD-035.md) | DONE | Keep RouteDeck running when the window closes: menu-bar tray, Open restores the window, every Quit path shows a Question dialog. |
| [`RD-034`](tasks/RD-034.md) | DONE | Browse installed macOS applications (with icons) behind ports, and add a chosen app to the catalog without auto-import. |
| [`RD-032`](tasks/RD-032.md) | DONE | Edit multi-profile policies and Apply/Stop from Dashboard with honest banners and event-driven snapshot refresh. |
| [`RD-031`](tasks/RD-031.md) | DONE | Ship Untitled UI forms for proxies, application picker with helper checklist, and Settings (path, debug logging, open logs folder). |
| [`RD-030`](tasks/RD-030.md) | DONE | Expose Proxies, Applications, Profiles, and Settings over Wails with typed errors, including debug logging and open-logs-folder. |
| [`RD-029`](tasks/RD-029.md) | DONE | Expose Apply/Stop and an extended snapshot over Wails, emit product events, stop CLI on shutdown, and re-check on wake/network change. |
| [`RD-028`](tasks/RD-028.md) | DONE | Orchestrate Apply/Stop with confirm, one rollback, Conflict handling, and honest ProtectionFacts — all tested with fake ports. |
| [`RD-027`](tasks/RD-027.md) | DONE | Manage one ProxyBridge CLI process with stdin profile, graceful stop, and a versioned line classifier with fixtures from real 4.1.3. |
| [`RD-026`](tasks/RD-026.md) | DONE | Serialize a domain rule plan into deterministic macOS `.pbprofile` JSON, loading passwords only inside the serializer. |
| [`RD-025`](tasks/RD-025.md) | DONE | Compile a profile into a wire-neutral rule plan and AppCoverage facts using the owner protocol-scope mapping. |
| [`RD-024`](tasks/RD-024.md) | DONE | Probe proxy TCP/SOCKS5/UDP ASSOCIATE health and keep a registry refreshed on triggers and a 30s ticker. |
| [`RD-023`](tasks/RD-023.md) | DONE | Discover ProxyBridge (custom path then documented paths), parse `--version`, and feed Snapshotter so the snapshot can leave permanent Unknown. |
| [`RD-022`](tasks/RD-022.md) | DONE | Inspect a selected `.app` (and nested helpers) through a platform-neutral port so Catalog can add applications without embedding plist logic in services. |
| [`RD-021`](tasks/RD-021.md) | DONE | Implement SecretStore with `github.com/keybase/go-keychain` so proxy passwords leave SQLite. |
| [`RD-020`](tasks/RD-020.md) | DONE | Persist policy scope, username, settings, and reject Delete of referenced proxies/apps with typed errors. |
| [`RD-019`](tasks/RD-019.md) | DONE | Extend domain policies and application ports so the slice can persist settings, delete with reference checks, inspect apps, and store secrets without inventing APIs later. |
| [`RD-018`](tasks/RD-018.md) | DONE | Wire structured `slog` logging to a rotating JSON file and stderr so every later adapter can leave an honest trail. |
| [`RD-017`](tasks/RD-017.md) | DONE | Ship a HashRouter shell with six screens; Dashboard loads the Go snapshot and shows prot... |
| [`RD-016`](tasks/RD-016.md) | DONE | Remove GreetService and the time event; wire sqlite and SnapshotService so the app retur... |
| [`RD-015`](tasks/RD-015.md) | DONE | Add application services that save catalog aggregates and build an honest protection sna... |
| [`RD-014`](tasks/RD-014.md) | DONE | Persist Application, ProcessIdentity, ProxyEndpoint (no password), and Profile in SQLite... |
| [`RD-006`](tasks/RD-006.md) | DONE | When persistence work starts, verify `AutoMigrate` and `gorm.io/driver/sqlite` against ... |
| [`RD-012`](tasks/RD-012.md) | DONE | Add only `ProxyStore.Get` and `SecretStore.Load` in `internal/application/port`, with n... |
| [`RD-011`](tasks/RD-011.md) | DONE | Write the 14-step ProtectionState derivation into CONCEPTION, add two glossary sentence... |
| [`RD-010`](tasks/RD-010.md) | DONE | Add Profile and ApplicationPolicy and reject unresolved proxy ids without rewriting act... |
| [`RD-009`](tasks/RD-009.md) | DONE | Add the three domain structs named in the domain model, with no password field and no U... |
| [`RD-013`](tasks/RD-013.md) | DONE | Make `golangci-lint run ./...` catch unformatted Go and non-exhaustive enum switches, fix t... |
| [`RD-008`](tasks/RD-008.md) | DONE | Add `internal/domain` value types and parsers with no aggregates and no protection deri... |
| [`RD-000`](tasks/RD-000.md) | DONE | Owner confirms `CONCEPTION.md` bounds for v1. |
| [`RD-001`](tasks/RD-001.md) | DONE | Record real Go, Wails, test, and build results for this scaffold. |
| [`RD-002`](tasks/RD-002.md) | DONE | Align module path, bundle metadata, and Dock icon with MACROX RouteDeck. |
| [`RD-003`](tasks/RD-003.md) | DONE | Owner accepts ADR-0001, ADR-0002, ADR-0004, and ADR-0005 (ADR-0003 is already the GORM ... |
| [`RD-004`](tasks/RD-004.md) | DONE | Owner confirms aggregates and the three state dimensions. |
| [`RD-005`](tasks/RD-005.md) | DONE | Replace any TBD path list with facts from current ProxyBridge macOS docs. |
