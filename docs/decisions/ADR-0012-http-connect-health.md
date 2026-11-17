# ADR-0012: HTTP proxy health uses CONNECT

Status: Accepted
Date: 2026-11-17
Owners: MACROX RouteDeck

## Context

ADR-0008 probed HTTP proxies with TCP connect only and never reported `AuthenticationFailed`. That left HTTP credentials unchecked. Owner closed CONNECT-based health for RD-039.

## Decision

- HTTP health always sends one `CONNECT example.com:443` request after TCP dial. Do not return `Available` from TCP alone.
- With no username and an empty password, omit `Proxy-Authorization`.
- With a non-empty username or a non-empty password, send preemptive `Proxy-Authorization: Basic` on that single request. No challenge retry.
- `407` maps to `AuthenticationFailed` with a nil error (including when no credentials were sent). Checker errors become `Unknown` in HealthMonitor, so auth failure must not be an error return.
- HTTP status `200` maps to `Available`.
- TCP dial failure, short read, deadline after TCP, or any other status maps to `Unavailable` with a nil error.
- Read the status line only, then close. No TLS handshake and no response body use.
- The read deadline is the checker's dial timeout (default 3s).
- HTTP still does not count as UDP protection. ADR-0008 compile rules stay.
- Do not log CONNECT request bytes or passwords.

This supersedes the ADR-0008 sentence that HTTP health is TCP-only and never `AuthenticationFailed`.

## Alternatives

- TCP-only HTTP probe. Rejected by owner for RD-039.
- Challenge-response after `407`. Rejected: one request only.
- `TCPAvailableUDPUnavailable` for HTTP. Rejected: HTTP has no UDP ASSOCIATE path here.

## Consequences

- Recheck opens a real CONNECT through the proxy to `example.com:443` on the existing triggers (startup, save, before Apply, Recheck, 30s ticker). A proxy that refuses that host can show `Unavailable` while still accepting other targets.
- Keychain load failures still probe with a nil password (unchanged HealthMonitor behavior).

## Follow-up

None required for v1.

## Open decisions (owner)

None.
