# Cursor workflow

Sources (principles only; no SaaS required):

- https://prodmoh.com/blog/10-best-practices-for-cursor-users
- https://www.vibecodingacademy.ai/blog/cursor-ai-best-practices-guide-2026

Upstream input is `CONCEPTION.md` + Accepted ADRs + contracts + the READY task file under `docs/project/tasks/`. That set replaces a vague prompt. `docs/project/BACKLOG.md` is an index; open one `RD-*.md` only.

## Prompt shape

Every task and packet starts with one line:

```text
Goal: ...
```

Then: Context, Current behavior, Desired behavior, Scope, Out of scope, Acceptance criteria, Constraints, Verification.

MUST `@` the task file, the contracts, and the packages you will edit. MUST NOT say “update auth” with no paths.

Screenshots are allowed for UI bugs. Acceptance criteria still go in the backlog.

If a criterion is missing, ask the owner and get an answer before planning or editing. Do not invent the criterion in Assumptions.

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
| Plan Mode | More than one package, architecture, ProxyBridge, protection — only after open questions are closed |
| Agent | Execute an accepted plan |
| Inline (Cmd+K) | One file, task already understood |

Plan flow: clarify open questions with the owner **first**, read docs and code, **then** write the plan, owner accepts, then build. MUST NOT draft a plan that invents product answers or papers over gaps. Save lasting plans under `.cursor/plans/` only. MUST NOT put plan files under `docs/`.

If a criterion, closed list, state mapping, or behavior choice is missing: stop and ask. Do not enter Plan Mode until the owner has answered.

Every plan MUST:

- Include an explicit **follow `AGENTS.md`** item (todo and/or DoD checkbox). Implied compliance is not enough.
- Segment todos into concrete steps (one command, one file group, or one verify). MUST NOT ship only two coarse todos such as “implement” and “document”.
- Include an explicit **Definition of Done (DoD)** section with checkable evidence (commands, acceptance boxes, docs, commit policy).
- List `golangci-lint run ./...` under Verification and as its own todo. If that command is already FAIL, the plan fixes it. A recorded FAIL does not close the plan. Skip the command only when the READY task quotes an owner waiver that names it.

See `docs/prompts/plan-mode.md`.

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
- golangci-lint run ./...
```
