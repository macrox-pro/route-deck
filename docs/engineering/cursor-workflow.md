# Cursor workflow

Sources (principles only; no SaaS required):

- https://prodmoh.com/blog/10-best-practices-for-cursor-users
- https://www.vibecodingacademy.ai/blog/cursor-ai-best-practices-guide-2026

Upstream input is `CONCEPTION.md` + Accepted ADRs + contracts + the READY `RD-*` task. That set replaces a vague prompt.

## Prompt shape

Every task and packet starts with one line:

```text
Goal: ...
```

Then: Context, Current behavior, Desired behavior, Scope, Out of scope, Acceptance criteria, Constraints, Verification.

MUST `@` the task file, the contracts, and the packages you will edit. MUST NOT say “update auth” with no paths.

Screenshots are allowed for UI bugs. Acceptance criteria still go in the backlog.

If a criterion is missing, list assumptions and ask the owner before editing.

## Layers

| Layer | File |
| --- | --- |
| Any agent | `AGENTS.md` |
| Cursor scope | `.cursor/rules/*.mdc` |
| Long standards | `docs/engineering/` |
| Packets | `docs/prompts/` |

## Modes

| Mode | When |
| --- | --- |
| Plan Mode | More than one package, architecture, ProxyBridge, protection, unclear scope |
| Agent | Execute an accepted plan |
| Inline (Cmd+K) | One file, task already understood |

Plan flow: clarify, read docs and code, edit the plan, owner accepts, then build. Save lasting plans under `.cursor/plans/` or `docs/project/plans/`.

## Parallel and cloud

One worktree per `RD-*`. See `docs/operations/parallel.md`.

Cloud agents receive the READY task id, `@` contracts, and verification commands. Same DoD. MUST NOT point auto-run at secrets.

Auto-run: Auto-review, not Run Everything. Allowlist notes: `docs/operations/auto-run.md`.

## Drift

STATUS, BACKLOG, and HANDOFF share one `RD-*`. Suggested and owner commit messages MUST be exactly one lowercase sentence, cite that id when applicable, and omit fluff or second paragraphs. Unmet acceptance criteria means not DONE.

## Optional

Voice input and keyboard shortcuts are optional. They are not process gates.

Terminal AI tools may draft ops scripts. Product code still follows this repo’s DoD.

## Anti-example

```text
Make the proxy stuff better and commit it.
```

## Example

```text
Goal: Reject profiles that reference a missing proxy id.

Context: @docs/contracts/proxybridge-profile.md @docs/architecture/domain-model.md
Current behavior: no validator yet.
Desired behavior: ValidateProxyRefs returns ErrUnresolvedProxy and does not call the process port.
Acceptance criteria:
- go test for the package PASS
- missing id does not call Start
Constraints:
- MUST NOT rewrite the rule to DIRECT
Verification:
- go test ./internal/domain/...
```
