# Persistence

Accepted stack:

| Piece | Choice |
| --- | --- |
| Engine | SQLite |
| Access | [GORM](https://gorm.io/docs/) only (`gorm.io/gorm`) |
| Driver | `gorm.io/driver/sqlite` (CGO, mattn/go-sqlite3) |
| Schema | GORM [`AutoMigrate`](https://gorm.io/docs/migration.html) |

Pure-Go drivers, goose, and golang-migrate are out of scope unless a future ADR supersedes ADR-0003.

## Placement

GORM models and `*gorm.DB` stay inside `internal/adapter/outbound/persistence/sqlite` (create that tree only when a persistence task starts).

Three types, mapped explicitly:

| Kind | Example |
| --- | --- |
| Domain | `domain` profile entity |
| Persistence | GORM struct in the sqlite adapter |
| Wails DTO | inbound DTO |

MUST NOT pass GORM models to the frontend. MUST NOT import `gorm.io/gorm` from domain or application packages. MUST NOT return `*gorm.DB` from repositories.

## Repositories

Shape repositories around an aggregate or use case, not a generic CRUD dump.

- Pass `context.Context` into every database call
- Enforce uniqueness in the schema, not only in Go
- MUST NOT hide queries inside GORM hooks
- MUST NOT put business rules in hooks
- Application services own transactions
- MUST NOT hold a transaction open while waiting on ProxyBridge

## AutoMigrate

Run `AutoMigrate` at startup from the composition root, on the known model list only. Document that list next to the call when it exists.

AutoMigrate can add missing tables, columns, and indexes. It does not drop unused columns or tables and does not reliably rename fields. A breaking shape change needs a READY task and a backup step described in `docs/operations/recovery.md`.

There is no `migrations/` directory in v1.

## SQLite settings

Set and document these on the DSN or via pragmas (driver-supported form):

| Setting | Requirement |
| --- | --- |
| `journal_mode` | WAL |
| `foreign_keys` | ON |
| `busy_timeout` | Set explicitly (milliseconds) |
| `synchronous` | Choose and document (NORMAL with WAL is the default plan) |
| Path | Application Support for RouteDeck |
| Pool | One writer. Prefer a small pool. SQLite allows one writer per file. Prove the pool with a test when the adapter exists. |

## Secrets

SQLite MUST NOT contain proxy passwords. Store a Keychain item reference. The profile serializer loads the secret at apply time and MUST NOT log the profile body.

## Minimum stored facts

Profiles, applications, process identities, proxies (no passwords), rules, profile links, events, settings, desired vs applied identifiers. Exact columns are defined with the first persistence task, not before.
