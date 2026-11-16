# ADR-0003: GORM, SQLite, AutoMigrate

Status: Accepted
Date: 2026-11-16
Owners: MACROX RouteDeck

## Context

RouteDeck needs a local database for profiles, apps, proxies, rules, events, and settings. Passwords must not sit in that database.

## Decision

Use SQLite through GORM (`gorm.io/gorm`) and the official driver `gorm.io/driver/sqlite` (CGO, mattn/go-sqlite3). Schema changes use GORM `AutoMigrate` on an explicit model list at startup. GORM types stay in the persistence adapter. Domain and application code MUST NOT import GORM.

`gorm.io/driver/sqlite` is the driver because it is the official GORM SQLite driver and this product ships as a macOS Wails app where CGO is available. Pure-Go drivers are not used in v1.

## Alternatives

- goose or golang-migrate. Rejected for v1. The owner accepted AutoMigrate.
- `database/sql` without an ORM. Rejected. The owner required GORM.
- Pure-Go SQLite drivers. Rejected for v1 in favor of the official driver.

## Consequences

- Builds need a C toolchain (normal for this Wails macOS app).
- AutoMigrate adds missing tables and columns. It does not drop or reliably rename. Breaking changes need a READY task and a backup.
- No `migrations/` directory in v1.
- WAL, foreign keys, busy timeout, and synchronous mode are set in the adapter (see `docs/architecture/persistence.md`).

## Follow-up

RD-006 records installed module versions and a smoke AutoMigrate when persistence code starts.

## Open decisions (owner)

None.
