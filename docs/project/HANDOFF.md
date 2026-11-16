# Handoff

Timestamp: 2026-11-16 20:44 UTC
Agent: agent (RD-002 lint PASS)
Task ID: RD-002
Branch / commit: uncommitted working tree (agent MUST NOT commit)

## Done

- Module renamed to `github.com/macrox-pro/route-deck`; `go mod vendor` refreshed
- `build/config.yml` identity: MACROX TECHNOLOGIES / RouteDeck / `pro.macrox.routedeck` / 0.1.0
- Brand master copied to `build/appicon.png`; PNG-only `generate:icons`; template `Assets.car` removed
- Darwin plists: identifier `pro.macrox.routedeck`, min OS `13.0.0`, no `CFBundleIconName`
- Frontend bindings under `frontend/bindings/github.com/macrox-pro/route-deck`; `App.tsx` import updated
- `APP_BUNDLE_NAME` = `RouteDeck` → package `bin/RouteDeck.app` with binary `Contents/MacOS/route-deck`
- Windows MSIX templates aligned so leftover rejected-id grep is clean (v1 still macOS-only)
- Docs: STATUS, BACKLOG, HANDOFF, build.md, development.md, branding README, README identity
- Package comment on `greetservice.go` so revive `package-comments` passes for package `main`

## Files changed

- `go.mod`, `vendor/` (refresh)
- `build/config.yml`, `build/Taskfile.yml` (PNG-only icons)
- `Taskfile.yml`, `build/darwin/Taskfile.yml` (`APP_BUNDLE_NAME`)
- `build/appicon.png`, `build/darwin/icons.icns`, `build/windows/icon.ico`
- `build/darwin/Info.plist`, `build/darwin/Info.dev.plist` (+ asset-task side effects on ios/windows/linux)
- `build/windows/msix/template.xml`, `build/windows/msix/app_manifest.xml`
- `frontend/src/App.tsx`, `frontend/bindings/`
- `greetservice.go` (package comment)
- `docs/project/STATUS.md`, `BACKLOG.md`, `HANDOFF.md`
- `docs/operations/build.md`, `docs/operations/development.md`
- `docs/branding/README.md`, `README.md`
- Lint gate (this turn): `AGENTS.md`, `.cursor/rules/core.mdc`, `.cursor/rules/golang.mdc`, `docs/engineering/golang.md`, `docs/engineering/cursor-workflow.md`, `docs/prompts/plan-mode.md`, `feature.md`, `bugfix.md`, `refactor.md`, `test-gen.md`, `handoff.md`, `docs/templates/TASK.md`, `docs/templates/AGENT-REPORT.md`
- `.golangci.yml` excludes `^build/`; `docs/operations/development.md`; `docs/project/ROADMAP.md`
- ADR-0004 body unchanged (consequence “until RD-002” is now historical)

## Decisions

- Owner ordered `build/` excluded from golangci-lint. Path regex in `.golangci.yml` is `^build/`. That tree is the Wails scaffold, not product Go
- Owner rejected deferring lint. A recorded `golangci-lint run ./...` FAIL does not close a loop, and the next `RD-*` does not start while that command is FAIL, unless the READY task quotes an owner waiver for the command
- Company: `MACROX TECHNOLOGIES` (macrox.pro brand; English form)
- Version: `0.1.0`
- macOS bundle folder: `RouteDeck.app`; binary stays `route-deck` (ADR-0004)
- No Icon Composer `.icon` / `Assets.car` for v1

## Verification

- Command: `go mod vendor` — PASS
- Command: `wails3 task common:generate:icons` — PASS (PNG only)
- Command: `wails3 task common:update:build-assets` — PASS
- Command: `wails3 task common:generate:bindings` — PASS
- Command: `wails3 task darwin:package` — PASS — `bin/RouteDeck.app`
- Command: `go test ./...` — PASS — no test files
- Command: `golangci-lint config verify` — PASS
- Command: `gofmt -w greetservice.go` — PASS
- Command: `golangci-lint run ./...` — PASS — 0 issues. Cache write warnings only (sandbox)
- Command: PlistBuddy / resources / shasum / leftover grep — PASS (see STATUS)
- Owner Dock/Finder icon check — PASS — owner via `wails3 dev` (2026-11-16)

## Not done

- RD-003, RD-005 (still TODO; not started)
- Agent did not commit

## Risks / blockers

- Next `common:update:build-assets` resets `LSMinimumSystemVersion` to 12.0.0

## Exact next step

Owner reviews `docs/project/BACKLOG.md` RD-003 (ADR-0001..0006). Leave RD-003 TODO until that review.

## Suggested commit message

Owner commits. Agent did not run git commit. Exactly one lowercase sentence; no fluff.

```text
set rd-002 identity and satisfy golangci-lint with a package comment and build scaffold excluded
```

## Context to read

- `AGENTS.md`
- `docs/project/BACKLOG.md` (RD-002)
- `docs/decisions/ADR-0004-naming-identity.md`
- `docs/operations/build.md`

## Loop evidence (RD-002)

| Step | Evidence |
| --- | --- |
| Recover | STATUS + HANDOFF + BACKLOG RD-002 + plan |
| Select | RD-002 claimed READY → IN_PROGRESS |
| Understand | ADR-0004, ADR-0006, branding, Wails v3.0.0 asset/icon facts |
| Plan | Accepted RD-002 plan; owner chose RouteDeck.app bundle name |
| Implement | Module, config, icons, plists, bindings, APP_BUNDLE_NAME, package |
| Verify | `golangci-lint run ./...` PASS (0 issues); `go test ./...` PASS |
| Review | No Direct fallback; no secrets; no invented Wails flags; no commit |
| Document | STATUS, BACKLOG, HANDOFF share RD-002 DONE |
| Owner commit prep | Message above; `git commit` not run |
