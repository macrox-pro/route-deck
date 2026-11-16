# Handoff

Timestamp: 2026-11-16 20:34 UTC
Agent: agent (RD-002 owner icon PASS)
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

## Files changed

- `go.mod`, `vendor/` (refresh)
- `build/config.yml`, `build/Taskfile.yml` (PNG-only icons)
- `Taskfile.yml`, `build/darwin/Taskfile.yml` (`APP_BUNDLE_NAME`)
- `build/appicon.png`, `build/darwin/icons.icns`, `build/windows/icon.ico`
- `build/darwin/Info.plist`, `build/darwin/Info.dev.plist` (+ asset-task side effects on ios/windows/linux)
- `build/windows/msix/template.xml`, `build/windows/msix/app_manifest.xml`
- `frontend/src/App.tsx`, `frontend/bindings/`
- `docs/project/STATUS.md`, `BACKLOG.md`, `HANDOFF.md`
- `docs/operations/build.md`, `docs/operations/development.md`
- `docs/branding/README.md`, `README.md`
- ADR-0004 body unchanged (consequence “until RD-002” is now historical)

## Decisions

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
- Command: `go test ./...` — PASS
- Command: `golangci-lint run ./...` — FAIL — 5 revive (scaffold); blocks DONE
- Command: PlistBuddy / resources / shasum / leftover grep — PASS (see STATUS)
- Owner Dock/Finder icon check — PASS — owner via `wails3 dev` (2026-11-16)

## Not done

- Clearing scaffold golangci-lint revive FAILs (blocks RD-002 DONE / Phase 1 exit)
- RD-003, RD-005
- Agent did not commit

## Risks / blockers

- Lint FAIL keeps RD-002 at REVIEW
- Next `common:update:build-assets` resets `LSMinimumSystemVersion` to 12.0.0

## Exact next step

Clear scaffold revive findings (`golangci-lint run ./...`) or claim RD-003 while leaving RD-002 at REVIEW.

## Suggested commit message

Owner commits. Agent did not run git commit. Exactly one lowercase sentence; no fluff.

```text
set rd-002 module path, bundle id, and routedeck.app icon identity
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
| Verify | Commands above with PASS/FAIL/NOT RUN |
| Review | No Direct fallback; no secrets; no invented Wails flags; no commit |
| Document | STATUS, BACKLOG, HANDOFF share RD-002 REVIEW |
| Owner commit prep | Message above; `git commit` not run |
