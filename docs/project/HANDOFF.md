# Handoff

Timestamp: 2026-11-17 09:05 UTC
Agent: agent (RD-014..RD-017)
Task ID: RD-017
Branch / commit: uncommitted working tree on top of `6264c28` (agent MUST NOT commit)

## Done

- RD-014: SQLite catalog (`internal/adapter/outbound/persistence/sqlite`), GORM vendored, ports extended, `ErrUnsetAction` / `RoutingAction.IsSet`
- RD-015: `internal/application/service` Catalog + Snapshotter
- RD-016: removed Greet/time; `SnapshotService.Get`; main opens DB (fail-open to Unknown snapshot)
- RD-017: Untitled UI HashRouter shell; Dashboard badge from Go snapshot
- RD-006 closed by RD-014 smoke
- ADR-0007 Untitled UI; frontend AGENTS preface; recovery path documented

## Files changed (high level)

- `internal/domain`, `internal/application/port`, `internal/application/service`
- `internal/adapter/outbound/persistence/sqlite`
- `internal/adapter/inbound/wails`
- `main.go` (greetservice.go deleted)
- `frontend/` shell, styles move, React 19 / TS 6, bindings
- `vendor/`, `go.mod`, `go.sum`
- `docs/project/tasks/RD-014..017.md`, STATUS, BACKLOG, HANDOFF, ADR-0007

## Decisions

1. HashRouter for Wails asset 404 on deep paths.
2. DB open failure: window still opens; snapshot Unknown + error text.
3. ProtectionBadge uses BadgeWithDot; Unknown never success color.
4. Go method name collisions across ports solved with thin adapters on Store.

## Verification

- Command: `go test ./...` — PASS
- Command: `golangci-lint run ./...` — PASS — 0 issues
- Command: `npm run build` (frontend/) — PASS
- Command: `wails3 build` — PASS
- Launch: `bin/route-deck` created `~/Library/Application Support/route-deck/route-deck.db` with tables `applications`, `process_identities`, `proxies`, `profiles`, `application_policies`

## Not done

- ProxyBridge CLI adapter, Keychain adapter, health checks, forms, product events, tray
- RD-007 CI
- Agent did not commit

## Risks / blockers

- Large Untitled UI tree under `frontend/src/components` increases frontend bundle size
- Brand purple still default in theme.css

## Exact next step

Owner runs `wails3 dev`, reviews Dashboard Unknown badge, then commits. Next READY task is owner-chosen (ProxyBridge or Keychain).

## Suggested commit message

Owner commits. Agent did not run git commit. Exactly one lowercase sentence; no fluff.

```text
add rd-014..017 sqlite catalog and untitled ui shell so routedeck launches with an honest unknown snapshot
```

## Context to read

- `AGENTS.md`
- `frontend/AGENTS.md`
- `docs/project/STATUS.md`
- `docs/decisions/ADR-0007-untitled-ui.md`
- `docs/architecture/persistence.md`

## Loop evidence (RD-014..017)

| Step | Evidence |
| --- | --- |
| Recover | STATUS + HANDOFF + plan + AGENTS.md |
| Select | Owner accepted runnable catalog plan; task files RD-014..017 |
| Understand | ADR-0003, persistence.md, wails-api.md, DeriveProtection Unknown without CLI |
| Plan | Accepted plan; follow AGENTS.md; segmented todos; DoD; lint gate |
| Implement | sqlite + services + SnapshotService + HashRouter shell |
| Verify | go test PASS; lint 0; npm build PASS; wails3 build PASS; DB created |
| Review | no Direct fallback; no password column; GORM in adapter; no agent commit |
| Document | STATUS/BACKLOG/HANDOFF share RD-017 DONE |
| Owner commit prep | message above; `git commit` not run |
