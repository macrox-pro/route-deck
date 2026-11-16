# Handoff

Timestamp: 2026-11-16 21:00 UTC
Agent: agent (RD-005 Phase 0 exit)
Task ID: RD-005
Branch / commit: uncommitted working tree (agent MUST NOT commit)

## Done

- Owner accepted ADR-0001..0006 as written (RD-003). No ADR body edits.
- Match rule revised (RD-004): `processNames` patterns match bundle identifier or executable name; signature and path stay UX warnings; dashboard aggregate unchanged.
- ProxyBridge contracts frozen (RD-005) against live docs dated 2026-11-16:
  - CLI flags with URL; documented paths only; second-CLI refuse; signal/5s timeout; extension approval wait; extension-rejected-rules warning is not confirmation
  - Profile: match semantics, separators, UDP `*`+host miss, HTTP-UDP CLI direct vs RouteDeck `BLOCK`, TCP close on proxy failure, catch-all vs limited-`*` order
- Installed CLI version recorded: `ProxyBridge CLI 4.1.3` via `/usr/local/bin/proxybridge --version`
- CONCEPTION, domain-model, glossary, threat-model updated for the match rule

## Files changed

- `CONCEPTION.md`
- `docs/architecture/domain-model.md`
- `docs/project/GLOSSARY.md`
- `docs/security/threat-model.md`
- `docs/contracts/proxybridge-cli.md`
- `docs/contracts/proxybridge-profile.md`
- `docs/project/STATUS.md`, `BACKLOG.md`, `HANDOFF.md`

## Decisions

- Approving the Phase 0 exit plan is the owner note for RD-003 and RD-004
- No new ADR: none of ADR-0001..0006 stated the old “bundle id is not a match key” sentence
- Homebrew and other paths stay not documented; do not guess
- Do not copy docs sample version `4.1.0`; use the installed probe result

## Verification

- Command: open four ProxyBridge URLs and diff contracts — PASS (2026-11-16)
- Command: `/usr/local/bin/proxybridge --version` — PASS — `ProxyBridge CLI 4.1.3`
- Command: `rg` for `not the runtime match key` / `Not Bundle ID` — PASS — cleared from product docs (BACKLOG RD-004 rewritten)
- Command: `go test ./...` — PASS — no test files
- Command: `golangci-lint run ./...` — PASS — 0 issues. Cache write warnings only (sandbox)

## Not done

- Phase 2 domain ports (no READY task yet)
- RD-006, RD-007 (not started)
- Agent did not commit

## Risks / blockers

- Later compiler must not treat CLI HTTP-UDP direct, a missing rule, or a `*` plus host/port limit as UDP coverage
- Docs may drift from the installed binary (`4.1.3` today)

## Exact next step

`docs/project/BACKLOG.md` — owner adds a Phase 2 domain-ports READY task before any `internal/` packages.

## Suggested commit message

Owner commits. Agent did not run git commit. Exactly one lowercase sentence; no fluff.

```text
align rd-003 through rd-005 so processnames matches bundle id or process name
```

## Context to read

- `AGENTS.md`
- `docs/project/BACKLOG.md` (RD-005)
- `docs/contracts/proxybridge-cli.md`
- `docs/contracts/proxybridge-profile.md`
- `CONCEPTION.md` (data rules and UDP)

## Loop evidence (RD-005 / Phase 0 exit)

| Step | Evidence |
| --- | --- |
| Recover | STATUS + HANDOFF + BACKLOG RD-003..005 + Phase 0 exit plan |
| Select | Phase 0 exit: RD-003/004 owner accept via plan; RD-005 agent freeze |
| Understand | Live macos-cli/profiles/rules/how-it-works; CONCEPTION dashboard + UDP BLOCK |
| Plan | Accepted Phase 0 exit plan |
| Implement | CONCEPTION, domain-model, glossary, threat-model, two contracts, STATUS/BACKLOG/HANDOFF |
| Verify | `/usr/local/bin/proxybridge --version` PASS (`4.1.3`); `golangci-lint run ./...` PASS; `go test ./...` PASS |
| Review | No Direct fallback; no secrets; no invented flags; no ADR body edits; no new packages; no commit |
| Document | STATUS, BACKLOG, HANDOFF share RD-005 DONE |
| Owner commit prep | Message above; `git commit` not run |
