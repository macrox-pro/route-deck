# ADR-0004: Product identity

Status: Accepted
Date: 2026-11-16
Owners: MACROX RouteDeck

## Context

The Wails scaffold used placeholder names (`changeme`, template bundle ids). The product name is fixed.

## Decision

| Item | Value |
| --- | --- |
| Product | MACROX RouteDeck |
| Binary / Task `APP_NAME` | `route-deck` |
| Task IDs | `RD-*` |
| Bundle ID | `pro.macrox.routedeck` |
| Module | `github.com/macrox-pro/route-deck` |
| Docs and UI | English |
| Icon master | `docs/branding/route-deck-icon.png` |

Code comments are English and only where the why is non-obvious.

## Alternatives

- Keep the scaffold module name. Rejected.
- Reverse-DNS `macrox.pro.routedeck` as the bundle id. Rejected in favor of `pro.macrox.routedeck`.

## Consequences

- `go.mod` still says `changeme` until owner task RD-002.
- Generated plists may lag until `build/config.yml` is updated and assets are regenerated.

## Follow-up

RD-002 performs the rename and icon install. Agents do not rename the module inside unrelated tasks.

## Open decisions (owner)

None.
