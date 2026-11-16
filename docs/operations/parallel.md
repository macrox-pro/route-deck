# Parallel work

Use one git worktree per `RD-*` so two agents do not share a working tree.

```bash
git worktree add ../route-deck-rd-010 -b rd-010
```

Open that directory in its own Cursor window.

## Rules

- One owner per file. Two agents MUST NOT edit the same file.
- The owner of the task resolves conflicts.
- Cloud or background agents get the READY task id, `@` contracts, and the verification commands. Same DoD as local work.
- MUST NOT give cloud agents Keychain items, proxy passwords, or production credentials.
- MUST NOT use Cursor “Run Everything” auto-run on a machine that holds those secrets.

## Feature drift

STATUS, BACKLOG, and HANDOFF in that worktree name the same `RD-*`.
