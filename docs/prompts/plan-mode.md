# Plan Mode packet

Use before multi-file edits, architecture changes, or ProxyBridge/protection work.

## Gate (MUST pass before writing this packet)

If Scope, Acceptance, Constraints, closed lists, or behavior are not fully decided in `CONCEPTION.md`, an Accepted ADR, a contract, or the READY task: **stop**. Ask the owner. Close every open question. Only then fill the template below.

MUST NOT invent product answers, field lists, error names, or state mappings inside the plan. Assumptions cover non-blocking engineering choices only after product decisions are closed.

```text
Goal:

Context:
Current behavior:
Desired behavior:

Scope:
Out of scope:

Acceptance criteria:
- [ ]

Constraints:
- MUST NOT

Verification:
- command: golangci-lint run ./...
- command:

Related ADR / contracts:

Assumptions:
- 

Files to change:
-

Risks:
-

Task ID:

Todos (segmented; one step each):
- [ ] Follow AGENTS.md for this loop (hard bans, source of truth, owner-commit, evidence)
- [ ] Run `golangci-lint run ./...` and keep this plan open until it PASSes (a recorded FAIL does not close the loop)
- [ ] 
- [ ] 

Definition of Done:
- [ ] Execution followed AGENTS.md (no silent Direct fallback, no invented APIs, no agent commit, verification recorded)
- [ ] Acceptance criteria above are met or explicitly waived by the owner
- [ ] `golangci-lint run ./...` is PASS, or the READY task quotes an owner waiver that names this command. A recorded FAIL is not Done
- [ ] Each other Verification command recorded PASS / FAIL / NOT RUN with output or reason. A required FAIL blocks Done
- [ ] STATUS, BACKLOG, and HANDOFF share the same RD-* and next step
- [ ] Suggested commit message is one lowercase sentence (agent did not commit)
```

## Todo rules

- MUST include a first or early todo: **Follow AGENTS.md**.
- MUST segment todos so progress is visible. Prefer many small steps over two buckets.
- Forbidden: only “do the work” and “update docs”, or only “run checks” and “record results”, when the task has more than those two real steps.
- Each todo names a concrete action (a command, a path set, or a verify gate).
- MUST include a todo that runs `golangci-lint run ./...`. If that command is already FAIL, clearing the findings is a todo in this plan, not a later task. Drop that todo only when the READY task quotes an owner waiver that names the command.

## DoD rules

- Every plan MUST include a **Definition of Done** section.
- DoD MUST include a checkbox that execution followed `AGENTS.md`.
- DoD is checkable evidence, not vibes. Unrun required checks block Done. A recorded `golangci-lint` FAIL blocks Done. Do not start the next `RD-*` while that command is FAIL.

Do not edit those files until the owner accepts this plan. Then run the feature or bugfix packet for implementation.

Skip this packet only for a one-line change the owner explicitly requested.
