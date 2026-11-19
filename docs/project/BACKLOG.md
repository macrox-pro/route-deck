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

_None. Next work needs a new `RD-*` task (TODO → READY)._

## Done

Archived one-liners: [`COMPLETED.md`](COMPLETED.md) (RD-000..RD-062). Full task files removed from `tasks/`; recover from git if needed.
