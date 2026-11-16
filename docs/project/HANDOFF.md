# Handoff

Timestamp: 2026-11-16 21:42 UTC
Agent: agent (RD-008)
Task ID: RD-008
Branch / commit: uncommitted working tree (agent MUST NOT commit)

## Done

- Amended RD-008 Scope with owner decisions (`ErrUnknownProxyKind`, Parse* set, opaque `RoutingAction`)
- Added `internal/domain`: `errors.go`, `id.go`, `state.go`, `policy.go` and table tests
- Verification PASS for domain tests, import check, lint, full `go test ./...`, vendor status, bindings regenerate

## Files changed

- `docs/project/tasks/RD-008.md`
- `docs/project/BACKLOG.md`
- `docs/project/STATUS.md`
- `docs/project/HANDOFF.md`
- `internal/domain/errors.go`
- `internal/domain/id.go`
- `internal/domain/id_test.go`
- `internal/domain/state.go`
- `internal/domain/state_test.go`
- `internal/domain/policy.go`
- `internal/domain/policy_test.go`

## Decisions

1. Rejected `ProxyKind` returns `ErrUnknownProxyKind` (third sentinel). Owner, 2026-11-17.
2. Function set: five id parsers, three state parsers, `ParseProxyKind`; `ProtocolScope` / `UDPLimitation` constants only; type-prefixed state constants. Owner, 2026-11-17.
3. `RoutingAction` is opaque with unexported `proxyID` / `direct` / `block` (no internal kind type). Constructors `ProxyAction`, `DirectAction`, `BlockAction`. Zero value is not Direct. Owner, 2026-11-17.

## Verification

- Command: `go test ./internal/domain/...` — PASS — `ok github.com/macrox-pro/route-deck/internal/domain`
- Command: `go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/domain/...` — PASS — `errors fmt` only
- Command: `golangci-lint run ./...` — PASS — 0 issues. Cache write warnings only (sandbox)
- Command: `go test ./...` — PASS
- Command: `git status --short go.mod go.sum vendor/` — PASS — empty
- Command: `wails3 generate bindings -clean=true -time-type=Date -ts` then `git status --short frontend/bindings` — PASS — empty; processed 1 Service, 0 Models
- Command: `wails3 build` — NOT RUN — `main` does not import `internal/domain`
- macOS-only checks — NOT RUN — pure Go value types with no OS calls

## Not done

- RD-009, RD-010, RD-011, RD-012
- RD-006, RD-007
- Agent did not commit

## Risks / blockers

- Whitespace-only ids such as `" "` are accepted; RD-009 may revisit
- `domain-model.md` spells `RoutingAction` wire labels as `PROXY`/`DIRECT`/`BLOCK`; domain has no string form
- Opaque `RoutingAction` would serialize as `{}`; future Wails inbound must map to a DTO (`docs/contracts/wails-api.md`)
- MUST NOT load every `docs/project/tasks/RD-*.md` into one session

## Exact next step

Open `docs/project/tasks/RD-009.md` and check DoR before promoting to READY.

## Suggested commit message

Owner commits. Agent did not run git commit. Exactly one lowercase sentence; no fluff.

```text
add rd-008 domain ids, states, and routing actions so later tasks share typed values without a direct fallback
```

## Context to read

- `AGENTS.md`
- `docs/project/BACKLOG.md` (index only)
- `docs/project/tasks/RD-009.md`
- `docs/architecture/domain-model.md`
- `CONCEPTION.md` (States)

## Loop evidence (RD-008)

| Step | Evidence |
| --- | --- |
| Recover | STATUS + HANDOFF + RD-008 + CONCEPTION States + domain-model + boundaries + ADR-0001 + wails-api + golang/testing/cursor-workflow |
| Select | RD-008 was READY; DoR boxes all checked after owner-decision amendment |
| Understand | Closed lists in RD-008 Scope; no CONCEPTION edit; no aggregates; ports stay RD-012 |
| Plan | Accepted plan `rd-008_domain_value_types`; follow AGENTS.md, segmented todos, DoD, lint gate |
| Implement | `internal/domain` errors/id/state/policy + tests; STATUS/BACKLOG amended during loop |
| Verify | `go test ./internal/domain/...` PASS; import check PASS; `golangci-lint run ./...` PASS; `go test ./...` PASS; vendor empty PASS; bindings empty PASS |
| Review | No Direct fallback; Unknown not Active; no forbidden imports; no secrets; no extra Scope symbols; no agent commit |
| Document | STATUS, BACKLOG, HANDOFF, RD-008 share RD-008 DONE; next step is `docs/project/tasks/RD-009.md` |
| Owner commit prep | message above; `git commit` not run |
