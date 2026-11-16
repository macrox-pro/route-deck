# Handoff

Timestamp: 2026-11-16 22:15 UTC
Agent: agent (RD-010)
Task ID: RD-010
Branch / commit: uncommitted working tree on top of `6264c28` (agent MUST NOT commit)

## Done

- Promoted RD-010 to READY then IN_PROGRESS; recorded owner decisions in Scope
- Added `ApplicationPolicy`, `Profile`, `NewProfile`, `ValidateProxyRefs`
- Added sentinels `ErrUnresolvedProxy` and `ErrEmptyName` in `errors.go` only
- Unit tests: `TestNewProfile` and `TestValidateProxyRefs` (`known_id`, `missing_id`, `direct_is_not_a_proxy_ref`)
- Verification PASS for domain tests, lint, full `go test ./...`, import check, vendor status, gofmt

## Files changed

- `docs/project/tasks/RD-010.md`
- `docs/project/BACKLOG.md`
- `docs/project/STATUS.md`
- `docs/project/HANDOFF.md`
- `internal/domain/errors.go`
- `internal/domain/profile.go`
- `internal/domain/profile_test.go`

## Decisions

1. Promote RD-010 to READY. Owner, 2026-11-17.
2. Only `NewProfile(id ProfileID, name string, policies []ApplicationPolicy) (Profile, error)`. Error only when `name == ""` (`ErrEmptyName`, no trim). Empty ID, empty policies, and policy contents are not checked. `ValidateProxyRefs` is not called from the constructor. On error return the zero value. `ApplicationPolicy` is a plain struct with no constructor. Owner, 2026-11-17.
3. `ErrUnresolvedProxy` and `ErrEmptyName` are declared in `errors.go`. Do not change the package comment or any other text in that file. Owner, 2026-11-17.

## Verification

- Command: `go test ./internal/domain/...` — PASS — `ok github.com/macrox-pro/route-deck/internal/domain`
- Command: `golangci-lint run ./...` — PASS — 0 issues. Cache write warnings only (sandbox)
- Command: `go test ./...` — PASS
- Command: `go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/domain/...` — PASS — `errors fmt strings` only
- Command: `git status --short go.mod go.sum vendor/` — PASS — empty
- Command: `gofmt -l internal/domain/` — PASS — empty
- Command: `wails3 build` — NOT RUN — `main` does not import `internal/domain`
- macOS-only checks — NOT RUN — pure Go structs with no OS calls

## Not done

- RD-011, RD-012
- RD-006, RD-007
- Open RD-013 review items (unscheduled)
- Agent did not commit

## Risks / blockers

- `Profile` exported fields allow bypassing `NewProfile` and storing an empty name
- `NewProfile` does not validate empty `ID` or policy contents
- `ValidateProxyRefs` does not check `Name` or `ID`
- Package comment in `errors.go` still says “no aggregates”
- MUST NOT load every `docs/project/tasks/RD-*.md` into one session

## Exact next step

Open `docs/project/tasks/RD-011.md` or `docs/project/tasks/RD-012.md` and check DoR before promoting one to READY.

## Suggested commit message

Owner commits. Agent did not run git commit. Exactly one lowercase sentence; no fluff.

```text
add rd-010 profile validation so a missing proxy id returns an error and the action stays proxy
```

## Context to read

- `AGENTS.md`
- `docs/project/BACKLOG.md` (index only)
- `docs/project/tasks/RD-011.md` or `docs/project/tasks/RD-012.md`
- `docs/architecture/domain-model.md`
- `docs/contracts/proxybridge-profile.md`
- `docs/engineering/testing.md`

## Loop evidence (RD-010)

| Step | Evidence |
| --- | --- |
| Recover | STATUS + HANDOFF + RD-010 + domain-model + proxybridge-profile + ADR-0001 + golang/testing |
| Select | RD-010 DoR complete; RD-009 DONE; owner decisions closed; promoted to READY then IN_PROGRESS |
| Understand | Closed fields and constructor from RD-010 Scope + owner decisions; no ProxyStore; no DeriveProtection |
| Plan | Accepted plan `rd-010_profile_refs`; follow AGENTS.md, segmented todos, DoD, lint gate |
| Implement | `profile.go`, `profile_test.go`, two sentinels in `errors.go`; no Action rewrite |
| Verify | domain tests PASS; lint PASS; `go test ./...` PASS; imports `errors fmt strings`; vendor empty; gofmt empty |
| Review | No Direct rewrite; no ProxyStore; closed fields only; no one-off helper; errors.go is two sentinels; no agent commit |
| Document | STATUS, BACKLOG, HANDOFF, RD-010 share RD-010 DONE; next step is RD-011 or RD-012 DoR |
| Owner commit prep | message above; `git commit` not run |
