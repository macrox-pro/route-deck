# Backlog

Statuses: `TODO`, `READY`, `IN_PROGRESS`, `BLOCKED`, `REVIEW`, `DONE`, `DROPPED`.

Agents MUST NOT `git commit`. The owner commits after review.

Feature-drift rule: `STATUS.md` current task, this index, and `HANDOFF.md` task ID MUST match.

## How to read (context budget)

1. Use this file as an **index only**. Do not paste full task bodies here.
2. Open **one** task file: `docs/project/tasks/RD-XXX.md` matching `STATUS.md`.
3. MUST NOT `@` every file under `docs/project/tasks/` in one session.
4. DONE tasks live in [`COMPLETED.md`](COMPLETED.md) (one-line rows). Full bodies are removed from `tasks/` and remain in git history. Do not load COMPLETED unless the owner asks.

## Active

| ID | Status | Goal |
| --- | --- | --- |
| [`RD-059`](tasks/RD-059.md) | REVIEW | Persist Appearance in SQLite and sync native window chrome via recreate. |
| [`RD-058`](tasks/RD-058.md) | REVIEW | Link-time app version/commit for Settings About and sidebar badge. |
| [`RD-057`](tasks/RD-057.md) | REVIEW | CI macOS package workflow (`darwin:package` + zip artifact, no notarization). |
| [`RD-056`](tasks/RD-056.md) | REVIEW | Dock attention badge `!` when applied protection is not Active; vendor dock service. |
| [`RD-055`](tasks/RD-055.md) | REVIEW | Helpers dialog from HelperGroupDTO; auto-open after Add; Applications uses pendingCount. |
| [`RD-054`](tasks/RD-054.md) | REVIEW | Proxy editor window; proxies page list-only; guide prefill; password never pre-filled. |
| [`RD-053`](tasks/RD-053.md) | REVIEW | Profile editor window with HelperGroupDTO; profiles page list-only with delete confirm. |
| [`RD-052`](tasks/RD-052.md) | REVIEW | EditorDialogs host + dialog shell routes; tray focuses open editor. |
| [`RD-051`](tasks/RD-051.md) | REVIEW | HelperGroupDTO, sync ConfirmHelpers, rule collapse, ErrConflictingMatch, catalog.changed. |
| [`RD-050`](tasks/RD-050.md) | REVIEW | Nested helper discovery (Resources, Current, IsUpdater) + ChatGPT fixture tests. |
| [`RD-049`](tasks/RD-049.md) | REVIEW | ADR-0014 and desktop UX docs/contracts. |
| [`RD-048`](tasks/RD-048.md) | REVIEW | Profiles app-grouped editor (icon/Action/helpers + Apply to helpers); Review helpers preselect. |
| [`RD-047`](tasks/RD-047.md) | REVIEW | Profiles form single-column layout; add-all processes by application (skip duplicates silently). |
| [`RD-046`](tasks/RD-046.md) | REVIEW | Honest ProtectionState via snapshot Apps; Applied vs Protection UX; proxy/profile icons; health labels; helpers H2. |
| [`RD-045`](tasks/RD-045.md) | REVIEW | Native dialogs for ProxyBridge GUI conflict and extension approval; remove Dashboard Conflict banner. |
| [`RD-044`](tasks/RD-044.md) | REVIEW | Distinguish WARN and ERROR in RouteDeck logs so Activity severity matches. |
| [`RD-040`](tasks/RD-040.md) | REVIEW | Ship Setup guides screen with VPN-client detection and Port open facts (ADR-0013). |
| [`RD-041`](tasks/RD-041.md) | REVIEW | One-click Add this proxy from a guide (save + Recheck); form fallback without default port. |
| [`RD-042`](tasks/RD-042.md) | REVIEW | Suggestion banner on Dashboard/Proxies for detected open local ports. |
| [`RD-043`](tasks/RD-043.md) | REVIEW | Dashboard first-run checklist that auto-hides when complete. |
| [`RD-039`](tasks/RD-039.md) | REVIEW | Probe HTTP proxies with CONNECT and report AuthenticationFailed on 407. |
| [`RD-007`](tasks/RD-007.md) | REVIEW | Add CI workflow matching local checks. Do not require ESLint. |
| [`RD-038`](tasks/RD-038.md) | REVIEW | Show recent RouteDeck JSON log records on the Activity screen without passwords or a connection journal. |
| [`RD-037`](tasks/RD-037.md) | REVIEW | Settings Launch at Login via Wails Autostart without auto-Apply. |
| [`RD-033`](tasks/RD-033.md) | REVIEW | Prove on macOS with real ProxyBridge 4.1.3 that Apply reaches Active, failures stay honest, and the log file traces each scenario without secrets. SOCKS5: Shadowrocket 127.0.0.1:1082 (TCPAvailableUDPUnavailable). Apply→Active NOT RUN pending owner UI. |

## Done

Archived one-liners: [`COMPLETED.md`](COMPLETED.md) (RD-000..RD-036 except open REVIEW ids). Full task files removed from `tasks/`; recover from git if needed.
