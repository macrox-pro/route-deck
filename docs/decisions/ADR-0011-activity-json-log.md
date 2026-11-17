# ADR-0011: Activity screen reads the JSON log

Status: Accepted
Date: 2026-11-17
Owners: MACROX RouteDeck

## Context

CONCEPTION and the roadmap require an Activity screen for RouteDeck events and health. The vertical slice left a placeholder. Product events notify the UI but are not persisted. The rotating JSON log at `~/Library/Logs/route-deck/route-deck.log` already records protection, health, bridge, and operation lines without passwords.

## Decision

- Activity v1 lists recent complete records from `route-deck.log` and rotations `.1` through `.3` (`logfile.DefaultPath`, `DefaultKeep`).
- The screen is health and ops. It is not a ProxyBridge connection journal. `connection.observed` stays absent.
- Wails product events stay notifications. Activity does not store a second copy of those events.
- A port and `logfile` adapter open an injected path read-only. The Wails inbound service calls an application service and does not import `logfile`.
- Pagination is **Load more** with an exclusive RFC3339 cursor `before` (empty = newest page). Page size is **100**. The response includes `hasMore`. Offset/page-number pagination is rejected for a growing rotating log.
- The DTO always exposes `time`, `level`, `msg`, and the stable RouteDeck keys in `docs/engineering/logging.md` (`component`, `proxy_id`, `profile_id`, `bridge_state`, `protection_state`, `health`, `line_class`, `operation_id`, `err`).
- Owner decision (A): every other non-empty **string** attribute is also returned in `attrs` as meta chips, except dropped secrets `password`, `secret`, `token`, and `profile`. Non-string JSON values stay omitted. Empty strings are omitted.
- A missing or unreadable active log file is a wrapped error. Broken or non-JSON lines are skipped. An empty parseable file is an empty list.
- There is no live file watcher. The UI loads the newest page on open, on Refresh, and again on the same product event names used for snapshot reload. Load more appends older pages.

## Alternatives

- In-memory ring of Wails events. Rejected: does not survive restart.
- SQLite table for activity. Rejected: duplicates the log for v1.
- Parse ProxyBridge verbose connection lines. Rejected: no stable connection API in v1.
- Fixed cap of 200 with no cursor. Superseded by owner decision: page size 100 and `before` cursor.
- Offset / page-number pagination. Rejected: unstable on an appending rotating log.

## Consequences

- Bridge lines that pass CLI text as a slog attribute named `msg` collide with the slog message key; the UI shows the decoded message as written. Renaming that attribute is out of this ADR.
- Redaction remains write-time key replace. A secret inside a message string is not scanned again on read.

## Follow-up

Optional live tail. Optional rename of the bridge `msg` attribute. Diagnostic bundle remains separate.

## Open decisions (owner)

None.
