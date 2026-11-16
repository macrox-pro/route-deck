# ADR-0002: Official ProxyBridge CLI

Status: Accepted
Date: 2026-11-16
Owners: MACROX RouteDeck

## Context

Per-app proxying on macOS needs a network extension. RouteDeck will not ship its own.

## Decision

The only packet enforcer is the official ProxyBridge CLI. RouteDeck discovers the binary, checks the version, writes a profile to the CLI stdin (`--profile -`), and parses human-readable stdout with versioned fixtures. RouteDeck does not fork ProxyBridge and does not auto-update it. A missing proxy ID rejects the profile. There is no kill-switch guarantee if the CLI dies.

## Alternatives

- Custom Network Extension. Rejected for v1.
- Passing a profile file path. Rejected: the contract is stdin.
- Treating CLI text as JSON. Rejected: it is not a stable API.

## Consequences

- Profile changes restart the CLI and leave a gap with no rules.
- GUI ProxyBridge can steal the extension. RouteDeck must surface `Conflict` and MUST NOT restart in a loop.
- Unknown CLI versions warn. They are not silent success.

## Follow-up

Freeze install paths in RD-005 from live docs. Implement the adapter only after that contract update if paths were unknown.

## Open decisions (owner)

None. Install path strings that are not in upstream docs stay "not documented" in the contract until RD-005, which does not reopen this decision.
