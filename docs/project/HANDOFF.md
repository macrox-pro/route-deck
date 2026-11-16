# Handoff

Timestamp: 2026-11-16 21:56 UTC
Agent: agent (RD-009)
Task ID: RD-009
Branch / commit: uncommitted working tree on top of `6264c28` (agent MUST NOT commit)

## Done

- Promoted RD-009 to READY then IN_PROGRESS; recorded owner decisions in Scope
- Added `Application`, `ProcessIdentity` with `MatchableNames`, `ProxyEndpoint` with `NewProxyEndpoint`
- Added sentinels `ErrInvalidPort` and `ErrEmptyHost`
- Unit tests only for `MatchableNames` and `NewProxyEndpoint` (no `Application` test file; no reflect field inventories)
- Recorded public-method-only unit test rule in `docs/engineering/testing.md` and `.cursor/rules/golang.mdc`
- Verification PASS for domain tests, lint, full `go test ./...`, import check, vendor status, gofmt

## Files changed

- `docs/project/tasks/RD-009.md`
- `docs/project/BACKLOG.md`
- `docs/project/STATUS.md`
- `docs/project/HANDOFF.md`
- `docs/engineering/testing.md`
- `.cursor/rules/golang.mdc`
- `internal/domain/errors.go`
- `internal/domain/application.go`
- `internal/domain/process_identity.go`
- `internal/domain/process_identity_test.go`
- `internal/domain/proxy_endpoint.go`
- `internal/domain/proxy_endpoint_test.go`

## Decisions

1. Promote RD-009 to READY. Owner, 2026-11-17.
2. Only `NewProxyEndpoint(id ProxyID, name string, kind ProxyKind, host string, port int, ref KeychainRef) (ProxyEndpoint, error)`. No `Validate()`. Owner, 2026-11-17.
3. No extra checks: only `ErrInvalidPort` and `ErrEmptyHost`. Whitespace-only host is not trimmed. Owner, 2026-11-17.
4. Unit tests cover exported functions/methods only. Plain structs with no behavior are not unit-test targets. No tests written only for coverage or field inventories. Owner, 2026-11-17.

## Verification

- Command: `golangci-lint run ./...` (baseline) — PASS — 0 issues. Cache write warnings only (sandbox)
- Command: `go test ./internal/domain/...` — PASS — `ok github.com/macrox-pro/route-deck/internal/domain`
- Command: `golangci-lint run ./...` — PASS — 0 issues. Cache write warnings only (sandbox)
- Command: `go test ./...` — PASS
- Command: `go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/domain/...` — PASS — `errors fmt strings` only
- Command: `git status --short go.mod go.sum vendor/` — PASS — empty
- Command: `gofmt -l internal/domain/` — PASS — empty
- Command: `wails3 build` — NOT RUN — `main` does not import `internal/domain`
- macOS-only checks — NOT RUN — pure Go structs with no OS calls

## Not done

- RD-010, RD-011, RD-012
- RD-006, RD-007
- Open RD-013 review items (unscheduled)
- Agent did not commit

## Risks / blockers

- `ProxyEndpoint` exported fields allow bypassing `NewProxyEndpoint`
- `NewProxyEndpoint` does not validate `Kind` or empty `ID`
- `MatchableNames` does not deduplicate
- MUST NOT load every `docs/project/tasks/RD-*.md` into one session

## Exact next step

Open `docs/project/tasks/RD-010.md` and check DoR before promoting to READY.

## Suggested commit message

Owner commits. Agent did not run git commit. Exactly one lowercase sentence; no fluff.

```text
add rd-009 application, process identity, and proxy endpoint so profiles can reference typed catalog entries without passwords
```

## Context to read

- `AGENTS.md`
- `docs/project/BACKLOG.md` (index only)
- `docs/project/tasks/RD-010.md`
- `docs/architecture/domain-model.md`
- `docs/contracts/proxybridge-profile.md`
- `docs/engineering/testing.md`

## Loop evidence (RD-009)

| Step | Evidence |
| --- | --- |
| Recover | STATUS + HANDOFF + RD-009 + domain-model + proxybridge-profile + ADR-0001 + golang/testing after RD-013 |
| Select | RD-009 DoR complete; RD-008 and RD-013 DONE; owner promoted to READY |
| Understand | Closed fields and constructor from RD-009 Scope + owner decisions; no Profile; ports stay RD-012 |
| Plan | Accepted plan `rd-009_domain_structs`; follow AGENTS.md, segmented todos, DoD, lint gate |
| Implement | `application.go`, `process_identity.go`, `proxy_endpoint.go`, sentinels; tests only for public methods; testing.md updated |
| Verify | domain tests PASS; lint PASS; `go test ./...` PASS; imports `errors fmt strings`; vendor empty; gofmt empty |
| Review | No Direct fallback; no secrets; no GORM; no one-off helpers; no field-inventory tests; RD-013 review items untouched; no agent commit |
| Document | STATUS, BACKLOG, HANDOFF, RD-009 share RD-009 DONE; next step is `docs/project/tasks/RD-010.md` |
| Owner commit prep | message above; `git commit` not run |
