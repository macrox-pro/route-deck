# Handoff packet

```text
Goal: Record the Loop end state for RD-___

Context:
Current behavior:
Desired behavior:

Scope:
- docs/project/HANDOFF.md
- docs/project/STATUS.md
- docs/project/BACKLOG.md status line

Out of scope:
- New code

Acceptance criteria:
- [ ] Same task ID in STATUS, BACKLOG, and HANDOFF
- [ ] Verification commands show PASS, FAIL, or NOT RUN
- [ ] Exact next step is a command or path
- [ ] One lowercase-sentence suggested commit message included (no fluff)
- [ ] No git commit was run

Constraints:
- MUST NOT write "continue" as the next step
- MUST NOT mark the task DONE, and MUST NOT name the next `RD-*` as the next step, while `golangci-lint run ./...` is FAIL, unless the READY task quotes an owner waiver for that command

Verification:
- command: git status (read-only)
- command: git diff (read-only)

Related ADR / contracts:

Assumptions (empty if none):

Task ID:
```

Template: `docs/templates/HANDOFF.md`.
