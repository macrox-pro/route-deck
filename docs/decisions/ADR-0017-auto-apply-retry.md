# ADR-0017: One Auto Apply retry when proxy is Unavailable

Status: Accepted
Date: 2026-11-19
Owners: MACROX RouteDeck

## Context

ADR-0015 runs Auto Apply once at process start with no retry. Login often races the local VPN client: the proxy is not listening yet, Apply fails honestly, and the user must Apply from the Dashboard. Owner closed (2026-11-19): one limited retry plus a message that Auto Apply failed.

[`Protection.Apply`](../../internal/application/service/protection.go) rechecks health and still starts the CLI. There is no typed "proxy not listening" error. ADR-0002 forbids a restart loop when the GUI holds the extension. The Wails `ProtectionService.Apply` returns as soon as it schedules the coordinator; the real outcome is the synchronous coordinator `Apply`.

## Decision

- Startup path only (`MaybeAutoApplyAtStartup`). Manual Dashboard Apply stays one shot.
- Observe the **synchronous** `Protection.Apply` result. Do not retry the async Wails Apply return.
- Retry **once** only when that call returned an error **and** a proxy referenced by the desired profile is `ProxyHealthUnavailable`. Wait **5s** (engineering default), then call the same Apply once more. The second call starts only after the first has released `ErrOperationInProgress`.
- Do **not** retry when the error contains `BridgeLineGUIConflict`, `BridgeLineAnotherCLI`, or `BridgeLineExtensionTimeout`, even if a proxy is also `Unavailable`. Do not retry auth failure, "not installed", or a successful Apply that leaves Protection `Degraded`.
- If the first attempt is not retried, or the second attempt fails: existing failure events stay, plus a native Warning with stable ID `protection.auto_apply_failed`. English draft: Auto Apply did not succeed; Apply again from the Dashboard.
- MUST NOT rewrite policies to Direct. MUST NOT invent a new error type. MUST NOT loop until healthy.

This **supersedes** the ADR-0015 sentence that forbids a startup retry/backoff loop, for this one-shot Unavailable case only.

## Alternatives

- Retry until healthy. Rejected: owner asked for one limited retry.
- Retry every Apply error. Rejected: conflicts with ADR-0002 Conflict handling.

## Consequences

- RD-064 implements the gate, 5s wait, and Warning dialog.
- A dead port that still yields a successful CLI confirm (`Degraded`) is not retried and does not show this failure dialog.

## Follow-up

None required for v1.

## Open decisions (owner)

None.
