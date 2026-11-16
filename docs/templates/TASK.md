# RD-XXX: Title

Status: TODO
Priority: P1
Owner:
Reviewer:

## Goal

One line.

## Context

## Current behavior

## Desired behavior

## Dependencies

## Scope

## Out of scope

## Acceptance criteria

- [ ] Machine-checkable result

## Constraints

- MUST NOT

## Verification

- command: golangci-lint run ./...
- result: PASS | FAIL | NOT RUN
- command:
- result: PASS | FAIL | NOT RUN

## Risks

## Related ADR / contracts

## Assumptions

None.

## DoR (all required before READY)

- [ ] One-line Goal
- [ ] Context / Current behavior / Desired behavior filled
- [ ] Scope and Out of scope explicit
- [ ] Dependencies closed or listed as non-blocking
- [ ] Acceptance criteria are checkboxes and machine-checkable
- [ ] Constraints include MUST NOT
- [ ] Verification lists exact commands
- [ ] Related ADR/contracts linked
- [ ] No blocking open ADR
- [ ] Fits one Loop (else split)
- [ ] Reviewer or owner named if product-sensitive

## DoD (all required before DONE)

### Code

- [ ] Matches Goal, Scope, architecture, Accepted ADRs
- [ ] Names reflect essence; no vague Util/Helper packages
- [ ] No one-off helper without reuse or a real test need
- [ ] Errors wrapped with `%w`; no silent Direct fallback
- [ ] Contexts and goroutines owned and stopped
- [ ] No secrets in code, logs, SQLite, fixtures
- [ ] GORM only in the persistence adapter; no `*gorm.DB` leak
- [ ] Public contracts unchanged unless this task allows it

### Verification

- [ ] Table-driven tests updated when logic changed
- [ ] `go test` command + PASS/FAIL recorded
- [ ] `golangci-lint run ./...` is PASS (command + output). NOT RUN only with an owner waiver in the READY task that names this command. A recorded FAIL blocks DONE
- [ ] Build recorded when relevant
- [ ] macOS-only checks recorded or NOT RUN with reason
- [ ] Unrun required checks, and required checks that FAILed, block DONE

### Docs

- [ ] STATUS, BACKLOG, HANDOFF share this RD-*
- [ ] Contracts/ADR updated if behavior changed
- [ ] Known limitations listed
- [ ] Exact next step is a command or path
- [ ] Diff + one lowercase-sentence suggested commit message; agent did not commit

## Loop evidence

| Step | Evidence |
| --- | --- |
| Recover | |
| Select | |
| Understand | |
| Plan | |
| Implement | |
| Verify | |
| Review | |
| Document | |
| Owner commit prep | |

## Handoff

Exact next step:
