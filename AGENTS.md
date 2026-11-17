# AGENTS.md — MACROX RouteDeck

Portable operating contract for any coding agent. Cursor-specific scoping lives in `.cursor/rules/`. Deep standards live in `docs/engineering/`.

Product: **MACROX RouteDeck**. Binary: `route-deck`. Tasks: `RD-*`. Bundle ID (target): `pro.macrox.routedeck`. Module (target): `github.com/macrox-pro/route-deck`. Platform v1: **macOS desktop only**.

## Read order (every session)

Context may be compressed or lost. Prior chat is untrusted. Re-read after long gaps.

1. `AGENTS.md` (this file)
2. `docs/project/STATUS.md`
3. `docs/project/HANDOFF.md`
4. `docs/project/BACKLOG.md` **index only** (status match) — then **one** task file `docs/project/tasks/RD-XXX.md` (same `RD-*` as STATUS)
5. Related ADR under `docs/decisions/` and contracts under `docs/contracts/`
6. `docs/engineering/golang.md`, `testing.md`, `cursor-workflow.md` when touching Go, tests, or agent workflow
7. Matching packet in `docs/prompts/`
8. The code and tests you will change

Then change code. MUST NOT start from chat memory alone. MUST NOT `@` every file under `docs/project/tasks/` in one session. DONE history is `docs/project/COMPLETED.md` (one-liners only; open on owner request).

## Source of truth

When sources conflict, use this order:

1. Real code behavior and command output from this session
2. Accepted ADRs
3. Contracts and architecture docs
4. `CONCEPTION.md`
5. Chat memory (lowest)

If code disagrees with an Accepted ADR: record the gap and propose a new ADR. MUST NOT silently “fix” either side.

## Loop

```text
Recover Context
  -> Select READY Task
  -> Understand Code + Contracts
  -> Plan (Plan Mode if multi-file / unclear)
  -> Implement (minimal scope)
  -> Verify (tests + golangci-lint + listed commands)
  -> Review (boundaries, fail-open, naming, no one-off helpers)
  -> Update Docs / HANDOFF / STATUS
  -> Prepare owner commit (diff + message)
  -> Next Loop
```

Agents MUST NOT `git commit`, `git push`, or `git commit --amend`. Prepare a diff summary and a suggested message only. The owner commits.

Suggested commit messages MUST be exactly one lowercase sentence: cite the `RD-*` when applicable, state the why, and omit fluff, second paragraphs, and bullet lists.

### Loop Verification Matrix

Run at the end of every Loop. A step without evidence fails the Loop.

| Loop step | MUST have evidence | Fail = |
| --- | --- | --- |
| Recover | STATUS + HANDOFF + task ID read this session; `@` paths used or listed | Starting from chat memory only |
| Select | Task status was READY; dependencies clear | Working a TODO/BLOCKED task without an explicit claim; starting another `RD-*` while `golangci-lint run ./...` is FAIL and the task has no owner waiver |
| Understand | Touched contracts/ADRs named; current vs desired stated | Coding without citing the owner of the behavior |
| Plan | Open questions closed with the owner before plan text; Plan Mode plan written or accepted; follow-AGENTS.md item; segmented todos; DoD present | Planning while product questions are open; inventing closed lists in the plan; multi-file edits with no plan; missing follow-AGENTS.md; two coarse todos; missing DoD |
| Implement | Diff matches Scope only; naming and extract rules held | Scope creep, one-off helper sprawl, invented APIs |
| Verify | Exact commands + PASS/FAIL/NOT RUN in HANDOFF | “should work”; unrun lint or tests; a plan whose Verification omits `golangci-lint run ./...`; a recorded lint FAIL treated as a finished loop |
| Review | Explicit check: no Direct fallback; GORM boundaries; secrets; owner-commit rule | Skipping review on protection or network changes |
| Document | STATUS, BACKLOG, HANDOFF share the same `RD-*`; next step is a command or path | Stale STATUS or the word “continue” as the next step |
| Owner commit prep | One lowercase sentence suggested; no agent commit | Agent ran `git commit`; multi-sentence, padded, or capitalized message |

## Definition of Ready (DoR)

A task MUST NOT start until all are true:

- [ ] One-line `Goal` present
- [ ] `Context` / `Current behavior` / `Desired behavior` filled
- [ ] `Scope` and `Out of scope` explicit
- [ ] Dependencies closed or listed as non-blocking
- [ ] Acceptance criteria are checkboxes and machine-checkable
- [ ] `Constraints` include MUST NOT (files, APIs, behaviors)
- [ ] `Verification` lists exact commands
- [ ] Related ADR/contracts linked
- [ ] No blocking open ADR for this change
- [ ] Task fits one Loop (else decompose first)
- [ ] Reviewer or owner named when the change is product-sensitive

If any box fails, status stays `TODO` or `BLOCKED`, not `READY`.

## Definition of Done (DoD)

### Code

- [ ] Matches Goal, Scope, architecture boundaries, and Accepted ADRs
- [ ] Names reflect the thing they are; no vague `Util` / `Helper` / `Manager` packages
- [ ] No one-off helper extracted without a second call site or a real testability need
- [ ] Errors wrapped with `%w`; not swallowed; no silent Direct fallback
- [ ] Contexts and goroutines have an owner and a stop path
- [ ] No secrets in code, logs, SQLite, or fixtures
- [ ] GORM only inside the persistence adapter; no `*gorm.DB` returned outward
- [ ] Public contracts unchanged unless the task allows it

### Verification (record in HANDOFF)

- [ ] Table-driven unit tests added or updated when logic changed
- [ ] `go test` for affected packages: command + PASS/FAIL
- [ ] `golangci-lint run ./...` with repo `.golangci.yml` is PASS (command + output). `NOT RUN` only when the READY task contains an owner waiver that names this command. A recorded FAIL blocks DONE
- [ ] Affected build recorded when relevant (`wails3 build` or `wails3 task`)
- [ ] macOS-only checks: result on macOS, or `NOT RUN` with a reason
- [ ] Any unrun required check, or a required check that FAILed, blocks DONE (status stays `REVIEW` at most)

### Docs and handoff

- [ ] STATUS, BACKLOG, and HANDOFF updated for the same `RD-*`
- [ ] Contracts or ADR updated if behavior changed
- [ ] Known limitations listed
- [ ] Exact next step is a command or file path
- [ ] Diff and one lowercase-sentence suggested commit message prepared; agent did not commit

## Hard bans

- Silent Direct fallback when a proxy is missing, unresolved, or dead
- Showing `Unknown` as `Active` / protected
- Claiming tests, lint, or build passed without command output from this turn
- Inventing Wails, GORM, or ProxyBridge APIs or CLI flags
- Planning or writing backlog closed lists while product or contract questions are still open
- Inventing answers in a plan to cover missing owner decisions
- Expanding scope beyond the READY task
- Changing an Accepted ADR without a new ADR
- Secrets in logs, SQLite, or domain aggregates
- Creating empty packages “for later”
- Non-English docs or code comments
- Using any ORM other than GORM
- Implementing without Goal and acceptance criteria
- `git commit`, `git push`, amend, or force-push
- One-off helper sprawl and vague names
- Promising a kill switch or leak-free failure of ProxyBridge (v1 does not guarantee it)
- Adding ESLint as a required gate (not in v1; do not add it)
- Introducing goose or golang-migrate (schema is GORM `AutoMigrate`)
- Gitignoring or deleting committed `vendor/` (ADR-0006)
- Deferring a `golangci-lint run ./...` FAIL to a later task, or starting the next `RD-*` while that command is FAIL, unless the READY task quotes an owner waiver for that command

## Stack rules (short)

| Topic | Rule |
| --- | --- |
| Go | 1.26+ (Wails v3.0.0). Details: `docs/engineering/golang.md` |
| Vendor | Commit `vendor/` (`go mod vendor`). Refresh after dependency changes. See ADR-0006 |
| Tests | Table-driven by default. Details: `docs/engineering/testing.md` |
| Lint | `golangci-lint run ./...` using `.golangci.yml`. FAIL blocks DONE |
| DB | SQLite via GORM only (`gorm.io/gorm`) + `gorm.io/driver/sqlite` (CGO) |
| Schema | GORM `AutoMigrate` at controlled startup. No separate migration tool in v1 |
| GORM boundary | Models and `*gorm.DB` stay in the persistence adapter |
| Secrets | macOS Keychain. SQLite stores references only |
| UI | React has no routing business rules. After an action, request a snapshot |
| ProxyBridge | Official CLI only. Profile via stdin `--profile -`. Human text is not a stable API |
| Platform | macOS 13+ desktop. Scaffold folders for other OSes are not product scope |

Verify library APIs against **`vendor/`** (then installed module docs) before use. MUST NOT invent them.

## Modes

| Mode | When |
| --- | --- |
| Plan Mode | Multi-file work, architecture, ProxyBridge or protection changes, more than one package — only after open questions are closed |
| Agent | Execute an accepted plan or a READY `RD-*` across files |
| Inline edit | Single-file fix after the task is understood |

Skip Plan Mode only for a one-line change the owner explicitly requested.

### Plan gate (before any plan text)

If anything needed for Scope, Acceptance, Constraints, closed lists, or behavior is missing from `CONCEPTION.md`, an Accepted ADR, a contract, or the READY task: **stop**. Ask the owner. Close every open question. **Then** write the plan.

MUST NOT:

- Start Plan Mode while product or contract questions are still open
- Invent answers, field lists, error names, or mappings to “fill gaps” in the plan
- Treat Assumptions as a substitute for owner decisions on product behavior
- Write backlog or plan text that papers over unresolved ambiguity

Assumptions in a plan are only for non-blocking, reversible engineering choices the owner can reject. Product words, closed lists, and state mappings need owner answers first.

### Plan requirements

Every Plan Mode plan MUST include:

1. **Follow `AGENTS.md`** — an explicit plan item or DoD checkbox that execution follows `AGENTS.md` (source of truth, hard bans, owner-commit, verification evidence). MUST NOT be implied only.
2. **Segmented todos** — one todo per concrete step (command, file group, or verify). MUST NOT collapse the loop into two coarse items such as “run checks” and “update docs”. Prefer a checklist the owner can track step by step.
3. **Definition of Done (DoD)** — an explicit section stating what evidence closes the loop (acceptance boxes, commands + PASS/FAIL/NOT RUN, docs updated, no commit unless the owner asked).
4. **Lint gate** — Verification lists `golangci-lint run ./...` and one todo runs that command. If the command is already FAIL, clearing those findings is in this plan. A recorded FAIL does not close the plan. Omit the command only when the READY task contains an owner waiver that names it.

Packet: `docs/prompts/plan-mode.md`.

## Ambiguity

Order is fixed: **clarify → owner answers → plan → implement**.

If a decision is not in `CONCEPTION.md`, an Accepted ADR, a contract, or the READY task: **stop and ask the owner before planning**. Do not draft a plan, backlog task body, or closed list until the answer exists. Coding under open ambiguity is also forbidden.

Assumptions may appear in a plan only after product decisions are closed, and only for non-blocking engineering choices.

## Parallel work

- Prefer one git worktree per `RD-*`. See `docs/operations/parallel.md`.
- Two agents MUST NOT edit the same file.
- Cloud agents get the same READY task, contracts, verification commands, and DoD.
- MUST NOT use “Run Everything” auto-run against machines that hold secrets.

## Report

Use `docs/templates/AGENT-REPORT.md`. Forbidden phrases without evidence: “done”, “should work”, “tests pass”, “protection is active”.

## Correct vs forbidden

Correct: unresolved proxy ID returns a typed error and does not start the ProxyBridge process.

Forbidden: missing proxy ID is rewritten to action Direct so the profile still applies.
