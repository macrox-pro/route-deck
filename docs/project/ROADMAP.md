# Roadmap

Order follows dependencies. Calendar estimates are out of scope.

v1 delivery prefers one thin vertical slice after foundation docs: discover ProxyBridge CLI, compile one rule, confirm state, show a snapshot. Do not finish every screen before that slice.

## Epics

| Epic | Content | Depends on |
| --- | --- | --- |
| E00 | Requirements and ADRs | — |
| E01 | Go boundaries and composition root | E00 |
| E03 | Domain and application ports | E01 |
| E02 | SQLite via GORM adapter (after ports exist) | E03 |
| E04 | ProxyBridge CLI adapter | E00, E01 |
| E05 | Profile compiler | E03, E04 |
| E06 | Proxy health checks | E03 |
| E07 | macOS app discovery | E01, E03 |
| E08 | Protection coordinator | E03–E07 |
| E09 | Wails API and frontend foundation | E01, E03 |
| E10 | Main UI | E05–E09 |
| E11 | Tray, autostart, lifecycle | E08, E09 |
| E12 | Security and diagnostics | E02, E04, E06 |
| E13 | Release engineering | E10–E12 |
| E-slice | Vertical slice: one rule, confirmed state, snapshot | E04, E05, E09 (minimal) |

E02 MUST NOT start before domain ports exist. Persistence models map from domain types. GORM stays in the adapter.

E11 depends on the coordinator (E08), not only on the UI shell. Tray restore of a profile is coordinator behavior.

## Phases

| Phase | Exit |
| --- | --- |
| 0 Contracts | ADR-0001..0005 and ProxyBridge contract reviewed (RD-000, RD-003, RD-005) |
| 1 Toolchain | RD-001, RD-002, `golangci-lint run ./...` green (`build/` excluded) |
| 2 Domain ports | Entities and ports with table-driven tests; no GORM import |
| 3 Persistence | GORM adapter, AutoMigrate, Keychain port for secrets |
| 4 CLI adapter | One managed process, stdin profile, parser fixtures |
| 5 Slice | One app, one proxy, confirmed ProtectionState in UI snapshot |
| 6 Surfaces | Dashboard, applications, proxies, profiles, activity, settings |
| 7 Tray | Close window stays resident; quit dialog stops CLI; Launch at Login (ADR-0010) |
| 8 Release | Reproducible macOS build; signing limits stated if unsigned |

## Screens (after the slice)

Dashboard, Applications, Proxies, Profiles, Activity, Settings, onboarding.

Activity v1 shows RouteDeck events and health. It MUST NOT assume ProxyBridge prints a stable connection log.
