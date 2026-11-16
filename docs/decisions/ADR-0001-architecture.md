# ADR-0001: Layered architecture

Status: Accepted
Date: 2026-11-16
Owners: MACROX RouteDeck

## Context

RouteDeck must keep routing policy in Go, show state in React, and drive ProxyBridge without leaking SQL or CLI details into the domain.

## Decision

Use four layers: Wails inbound adapters, application services, domain, outbound adapters. Dependencies point inward. Composition starts in `main` and may move to `internal/bootstrap` later. React MUST NOT contain routing rules. GORM MUST NOT leave the persistence adapter.

## Alternatives

- Put policy in React and call the CLI from the UI. Rejected: two sources of truth and easy fail-open mistakes.
- A single `main` package for all logic. Rejected: ProxyBridge and SQLite would tangle with rules.

## Consequences

- More packages once features start, but no empty packages ahead of time.
- DTOs are mapped at the Wails boundary. Persistence models are mapped in the sqlite adapter.

## Follow-up

Implement packages only inside READY tasks. See `docs/architecture/target-layout.md`.

## Open decisions (owner)

None.
