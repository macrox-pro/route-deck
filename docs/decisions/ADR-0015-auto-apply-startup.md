# ADR-0015: Auto Apply at startup

Status: Accepted
Date: 2026-11-18
Owners: MACROX RouteDeck

## Context

ADR-0008 and ADR-0010 forbade Apply at process start so login launch would not pretend the VPN client was ready. Owner closed (2026-11-18): allow an opt-out **Auto Apply on launch** setting, **default ON**, independent of Launch at Login, so a desired profile with Proxy/Block rules can Apply through the existing path when the app starts.

## Decision

- Settings exposes **Auto Apply on launch** (`autoApply`). Persisted in SQLite on the settings row. Default **true** (including AutoMigrate on existing databases).
- Independent of Launch at Login (OS Autostart). Manual launch and login launch both may Auto Apply when the flag is on.
- On every process start, after services and event emit are wired, RouteDeck evaluates gates once. If all pass, it calls the existing `ProtectionService.Apply` (same stack as Dashboard: Health Recheck, Bridge locate, compile, start, confirm, rollback-on-fail).
- Closed gate list (verbatim):
  1. `settings.AutoApply == true`
  2. Desired profile set (`desiredSet` and non-empty ID)
  3. Desired profile has at least one application policy whose action is Proxy or Block (Direct-only or empty → quiet Info skip)
  4. Then the existing Apply path
- Skip or Apply failure leaves Unprotected / existing failure events. MUST NOT rewrite policies to Direct. MUST NOT invent a second Apply stack. MUST NOT Apply from tray.
- No startup retry/backoff loop in this ADR (superseded for one Unavailable retry by ADR-0017). Wake and network change keep health recheck only.
- This ADR **supersedes** the ADR-0008 sentence “No auto-apply at startup” and the ADR-0010 / lifecycle lines that Launch at Login MUST NOT auto-Apply. Launch at Login still does not imply Apply by itself; Auto Apply is a separate setting.

## Alternatives

- Default OFF. Rejected by owner: default ON.
- Auto Apply only when started from login items. Rejected: independent of Launch at Login.
- Retry until proxies are healthy. Rejected for this ADR: honest failure; user can Apply manually.
- Persist Launch at Login in SQLite. Still rejected (ADR-0010).

## Consequences

- RD-060 implements persistence, Settings toggle, gate helper, and composition-root one-shot.
- Login race (proxy not listening yet) may fail Apply honestly; known limitation.
- First-run extension approval Info may appear during Auto Apply (existing Apply path).

## Follow-up

One limited Unavailable retry at startup: ADR-0017 / RD-064.

## Open decisions (owner)

None.
