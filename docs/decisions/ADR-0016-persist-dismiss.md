# ADR-0016: Persist checklist and suggestion dismiss in SQLite

Status: Accepted
Date: 2026-11-19
Owners: MACROX RouteDeck

## Context

ADR-0013 kept checklist dismiss and suggestion-banner dismiss in memory for the process lifetime. Owner closed (2026-11-19): both dismisses must survive process restart via SQLite so first-run UX does not reset on every launch.

## Decision

- Persist **checklist dismiss** as one flag in SQLite, plus the four step booleans used for regression (bridge found, proxy available, application added, profile applied). Load them on process start before the first `Checklist` call.
- Clear checklist dismiss when any previously-true step regresses (same rule as ADR-0013 / RD-061). Improving a step after dismiss does not clear dismiss.
- Persist **suggestion dismiss** as rows of `(guideID, port)`. Load on process start. No undismiss API. Dismissed suggestions stay hidden across restart.
- Models and `*gorm.DB` stay in the sqlite persistence adapter. Application use cases call ports only.
- No secrets in these tables or logs.

This **supersedes** the ADR-0013 alternative that rejected SQLite persist for suggestion and checklist dismiss, and the ADR-0013 sentences that dismiss is in-memory for the process lifetime only.

## Alternatives

- Session-only dismiss. Rejected by owner for this stage.
- Persist checklist flag without fingerprint. Rejected: regression clear would break across restart.

## Consequences

- New `DismissStore` port (or equivalent) and sqlite models / AutoMigrate entries.
- `SetupGuides` loads persisted state at startup and writes on dismiss / regression clear.
- RD-063 implements the change. Contracts (`wails-api.md`) stop saying dismiss is session-only.

## Follow-up

None required for v1.

## Open decisions (owner)

None.
