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
| [`RD-006`](tasks/RD-006.md) | TODO | When persistence work starts, verify `AutoMigrate` and `gorm.io/driver/sqlite` against ... |
| [`RD-007`](tasks/RD-007.md) | TODO | Define CI commands that match local checks. Do not require ESLint. |
| [`RD-011`](tasks/RD-011.md) | TODO | Write the 14-step ProtectionState derivation into CONCEPTION, add two glossary sentence... |
| [`RD-012`](tasks/RD-012.md) | TODO | Add only `ProxyStore.Get` and `SecretStore.Load` in `internal/application/port`, with n... |

## Done

| ID | Status | Goal |
| --- | --- | --- |
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
