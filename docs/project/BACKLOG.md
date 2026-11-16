# Backlog

Statuses: `TODO`, `READY`, `IN_PROGRESS`, `BLOCKED`, `REVIEW`, `DONE`, `DROPPED`.

Agents MUST NOT `git commit`. The owner commits after review.

Feature-drift rule: `STATUS.md` current task, this file's active task, and `HANDOFF.md` task ID MUST match.

---

# RD-000: Confirm conception boundaries

Status: DONE
Priority: P0
Owner: owner
Reviewer: owner

## Goal

Owner confirms `CONCEPTION.md` bounds for v1.

## Context

Docs-prep encoded product decisions. They are not owner-accepted until this task closes.

## Current behavior

`CONCEPTION.md` exists as a draft of locked decisions from the docs-prep plan.

## Desired behavior

Owner marks this task DONE or lists required edits.

## Scope

- Read `CONCEPTION.md`, `AGENTS.md`, ADR-0001.

## Out of scope

- Writing product code
- Renaming the Go module

## Acceptance criteria

- [x] Owner reply names RD-000 as accepted or lists concrete edits
- [x] No silent change to fail-open, UDP, or kill-switch wording

## Constraints

- MUST NOT invent new policies
- MUST NOT commit

## Verification

- command: owner review (no shell required)
- result: PASS — owner accepted CONCEPTION.md and ADR-0001 as written (2026-11-16). No edits requested. Fail-open, UDP, and kill-switch wording unchanged.

## Risks

- Shipping code before bounds are accepted

## Related ADR / contracts

- `docs/decisions/ADR-0001-architecture.md`
- `CONCEPTION.md`

## Handoff

Exact next step after accept: start RD-001 commands and record output in STATUS.

---

# RD-001: Verify scaffold toolchain

Status: DONE
Priority: P0
Owner: owner or agent
Reviewer: owner

## Goal

Record real Go, Wails, test, and build results for this scaffold.

## Context

`go.mod` declares `go 1.26.0` and `github.com/wailsapp/wails/v3 v3.0.0`. Module name is still `changeme`.

## Current behavior

Versions are documented but not re-verified in STATUS after docs-prep.

## Desired behavior

STATUS and HANDOFF contain command output for version and smoke checks.

## Scope

- Run the commands below on the owner machine
- Paste results into STATUS and HANDOFF

## Out of scope

- Fixing product bugs
- Module rename (RD-002)

## Acceptance criteria

- [x] `go version` output recorded
- [x] `wails3 doctor` output recorded (or NOT RUN with reason)
- [x] `go test ./...` result recorded
- [x] `wails3 build` or `wails3 task build` result recorded or NOT RUN with reason

## Constraints

- MUST NOT claim PASS without output
- MUST NOT commit

## Verification

- command: `go version` — PASS — `go version go1.27.4 darwin/arm64`
- command: `wails3 doctor` — PASS — exit 0; "Your system is ready for Wails development!" (doctor reports Go go1.27.3)
- command: `go test ./...` (before build) — FAIL — `main.go:17:12: pattern all:frontend/dist: no matching files found`
- command: `wails3 task build` — PASS — exit 0; produced `bin/route-deck` (ld macOS version warnings only)
- command: `go test ./...` (after build) — PASS — packages have no test files; exit 0
- command: `golangci-lint run ./...` (after build) — FAIL — 5 `revive` issues in scaffold (`greetservice.go`, `build/ios/...`)
- Environment: node v22.23.5; npm 10.9.9; golangci-lint v2.15.0

## Risks

- Doctor or build fails on missing Xcode CLI tools

## Related ADR / contracts

- `docs/decisions/ADR-0005-golang-toolchain.md`
- `docs/operations/development.md`

## Handoff

Exact next step: owner promotes RD-002 from TODO to READY and starts identity work.

---

# RD-002: Identity and app icon install

Status: DONE
Priority: P0
Owner: owner
Reviewer: owner

## Goal

Align module path, bundle metadata, and Dock icon with MACROX RouteDeck.

## Context

Scaffold uses module `changeme`. Target module is `github.com/macrox-pro/route-deck`. Target bundle ID is `pro.macrox.routedeck`. Brand master is `docs/branding/route-deck-icon.png`.

## Current behavior

Wails template identity and possibly `macrox.pro.routedeck` in a generated plist. Dock icon is the template icon.

## Desired behavior

`go.mod` module path is `github.com/macrox-pro/route-deck`. `build/config.yml` names MACROX TECHNOLOGIES / RouteDeck / `pro.macrox.routedeck`. Packaged app is `bin/RouteDeck.app` (binary `route-deck`) and shows the RouteDeck icon.

## Scope

- Rename module and fix imports
- Run `go mod vendor` after the rename and keep `vendor/` in sync (ADR-0006)
- Update `build/config.yml` info fields
- Replace `build/appicon.png` from the brand master and regenerate platform icons per Wails v3.0.0 docs
- Verify Finder/Dock icon on a local macOS package

## Out of scope

- Feature UI
- Tray glyph design

## Acceptance criteria

- [x] `go.mod` first line is `module github.com/macrox-pro/route-deck`
- [x] `go mod vendor` refreshed; `vendor/` present and not gitignored
- [x] `go test ./...` PASS after rename
- [x] `build/config.yml` productIdentifier is `pro.macrox.routedeck`
- [x] Local `.app` icon matches the brand master (owner visual check recorded) — PASS via `wails3 dev` Dock (2026-11-16)

## Constraints

- MUST NOT invent a second brand
- MUST keep `docs/branding/route-deck-icon.png` as the master
- MUST NOT drop `vendor/` (ADR-0006)
- MUST NOT commit unless the owner chooses to

## Verification

- command: `go mod vendor` — PASS
- command: `go test ./...` — PASS (after package; no test files)
- command: `wails3 task darwin:package` — PASS — `bin/RouteDeck.app`
- command: `golangci-lint run ./...` — PASS — 0 issues (`build/` excluded; package comment on `greetservice.go`)
- owner: Dock/Finder icon via `wails3 dev` — PASS

## Risks

- Regenerating assets overwrites hand-edited plist fields

## Related ADR / contracts

- `docs/decisions/ADR-0004-naming-identity.md`
- `docs/decisions/ADR-0006-go-vendor.md`
- `docs/branding/README.md`

## Handoff

Exact next step: RD-003 (DONE in Phase 0 exit).

---

# RD-003: Accept architecture ADRs

Status: DONE
Priority: P1
Owner: owner
Reviewer: owner

## Goal

Owner accepts ADR-0001, ADR-0002, ADR-0004, and ADR-0005 (ADR-0003 is already the GORM decision to confirm in the same pass). Confirm ADR-0006 (vendor) in the same pass or note it already matches owner action.

## Context

ADRs were written during docs-prep as Accepted drafts pending owner confirmation.

## Current behavior

Status fields say Accepted in-repo before owner sign-off.

## Desired behavior

Owner confirms or replaces status with comments. Agents MUST NOT change Accepted text silently.

## Scope

- Read the five ADRs
- Comment in this task or flip status if rejecting

## Out of scope

- New product features

## Acceptance criteria

- [x] Each ADR-0001..0005 has an owner note: accept or required edit — owner accepted ADR-0001..0006 as written via Phase 0 exit plan (2026-11-16). No ADR body edits.
- [x] ADR-0006 (vendor) accepted or required edit noted — accepted as written
- [x] No contradictory guarantee vs `CONCEPTION.md` — kill-switch and silent-Direct bans unchanged; match-rule revision is in CONCEPTION, not in ADR bodies

## Constraints

- MUST NOT edit Accepted ADR bodies without a superseding ADR

## Verification

- command: owner review — PASS (Phase 0 exit plan accepted 2026-11-16)

## Risks

- Agents treat draft Accepted as final before owner review

## Related ADR / contracts

- `docs/decisions/`

## Handoff

Exact next step: RD-004 review of glossary and states (DONE in same Phase 0 exit).

---

# RD-004: Review domain glossary and states

Status: DONE
Priority: P1
Owner: owner
Reviewer: owner

## Goal

Owner confirms aggregates and the three state dimensions.

## Context

`docs/architecture/domain-model.md` and `docs/project/GLOSSARY.md` define Application, ProcessIdentity, ProxyEndpoint, Profile, policy, and ProtectionSnapshot.

## Current behavior (at close)

`processNames` patterns match the bundle identifier or the executable name. Signature and path stay UX warnings. Dashboard aggregate unchanged: `Active` only when the bridge is ready and every protected app is fully covered.

## Desired behavior

Owner accepts that model or lists field-level edits.

## Scope

- Read domain-model and glossary
- Revise match wording to live ProxyBridge rules (Phase 0 exit)

## Out of scope

- Implementing types in Go

## Acceptance criteria

- [x] Owner accepts processNames as the enforcement key — accepted; patterns match bundle id or executable name (Phase 0 exit plan 2026-11-16)
- [x] Owner accepts Dashboard aggregate rule in CONCEPTION — unchanged and accepted

## Constraints

- MUST NOT add a fourth policy in this task

## Verification

- command: owner review — PASS (Phase 0 exit plan accepted 2026-11-16)

## Risks

- Later code treats a `*` rule that has a host or port limit as UDP coverage, or writes only one name and misses a helper process

## Related ADR / contracts

- `docs/architecture/domain-model.md`
- `docs/contracts/proxybridge-profile.md`

## Handoff

Exact next step: RD-005 against live ProxyBridge docs (DONE in same Phase 0 exit).

---

# RD-005: Freeze ProxyBridge contract against live docs

Status: DONE
Priority: P0
Owner: agent
Reviewer: owner

## Goal

Replace any TBD path list with facts from current ProxyBridge macOS docs.

## Context

`docs/contracts/proxybridge-cli.md` links official docs and forbids invented flags. Install locations are not frozen.

## Current behavior (at close)

Contracts quote only documented flags and paths. Match semantics, UDP limits, catch-all order, and CLI process facts are frozen against live docs dated 2026-11-16. Installed CLI: `ProxyBridge CLI 4.1.3`.

## Desired behavior

Contract quotes only flags and limits that appear in the live docs for the version we support. Unknown items stay explicitly unknown.

## Scope

- Re-read https://proxybridge.interceptsuite.com/docs/macos-cli/ and profile/rules/how-it-works pages
- Update the two contract files
- Record doc date and CLI version if installed

## Out of scope

- Implementing the adapter
- Auto-updating ProxyBridge

## Acceptance criteria

- [x] Every CLI flag in the contract has a doc URL anchor or is marked unknown — all flags cite https://proxybridge.interceptsuite.com/docs/macos-cli/
- [x] Install path list is either filled from docs or explicitly “not documented; owner path required” — documented: `/usr/local/bin/proxybridge`, `/Applications/ProxyBridge.app/Contents/MacOS/ProxyBridge`; Homebrew and others not documented
- [x] No new flag added from memory

## Constraints

- MUST NOT invent flags
- MUST NOT treat stdout as JSON

## Verification

- command: open the four ProxyBridge doc URLs and diff the contract — PASS (2026-11-16)
- command: `/usr/local/bin/proxybridge --version` — PASS — `ProxyBridge CLI 4.1.3`

## Risks

- Docs drift from the installed binary

## Related ADR / contracts

- `docs/decisions/ADR-0002-proxybridge.md`
- `docs/contracts/proxybridge-cli.md`
- `docs/contracts/proxybridge-profile.md`

## Handoff

Exact next step: `docs/project/BACKLOG.md` — no Phase 2 domain-ports task yet; owner adds that task before any `internal/` packages.

---

# RD-006: Confirm GORM SQLite driver docs before persistence code

Status: TODO
Priority: P1
Owner: agent
Reviewer: owner

## Goal

When persistence work starts, verify `AutoMigrate` and `gorm.io/driver/sqlite` against the installed module docs.

## Context

ADR-0003 locks GORM, SQLite, AutoMigrate, and `gorm.io/driver/sqlite` (CGO). No alternate driver.

## Current behavior

GORM is not a dependency yet.

## Desired behavior

First persistence change records the installed `gorm.io/gorm` and `gorm.io/driver/sqlite` versions and a smoke open of an on-disk SQLite file with WAL and foreign keys.

## Scope

- Add dependencies only inside the persistence task that needs them
- Smoke-check AutoMigrate on a temp database

## Out of scope

- Pure-Go SQLite drivers
- goose or golang-migrate

## Acceptance criteria

- [ ] `go.mod` requires `gorm.io/gorm` and `gorm.io/driver/sqlite`
- [ ] Smoke test opens SQLite and runs AutoMigrate on a sample model
- [ ] `go test` for that package PASS recorded

## Constraints

- MUST NOT import GORM from domain or application packages
- MUST NOT store secrets in SQLite

## Verification

- command: `go test ./internal/adapter/outbound/persistence/sqlite/...` (path may differ; record the real path)

## Risks

- CGO toolchain missing on the build machine

## Related ADR / contracts

- `docs/decisions/ADR-0003-gorm-sqlite.md`
- `docs/architecture/persistence.md`

## Handoff

Exact next step: record module versions in HANDOFF.

---

# RD-007: CI baseline design

Status: TODO
Priority: P2
Owner: agent
Reviewer: owner

## Goal

Define CI commands that match local checks. Do not require ESLint.

## Context

`.golangci.yml` exists. Frontend check is `tsc` via Vite scripts in `frontend/package.json`.

## Current behavior

No CI workflow file is required by docs-prep.

## Desired behavior

A later change may add CI that runs `golangci-lint run ./...`, `go test ./...`, and `npm run build` in `frontend/`. Until that file exists, this task stays the spec.

## Scope

- Document the command list in this task (already here)
- Optional: add a CI workflow only when the owner marks this READY and asks for the file

## Out of scope

- ESLint
- macOS system tests in CI unless a runner exists

## Acceptance criteria

- [ ] Command list matches `docs/operations/development.md`
- [ ] No ESLint step
- [ ] If a workflow file is added, a dry description of jobs is in the PR/commit message the owner writes

## Constraints

- MUST NOT add ESLint
- MUST NOT commit (owner commits)

## Verification

- command: `golangci-lint run ./...`
- command: `go test ./...`
- command: `npm run build` (cwd `frontend/`)

## Risks

- CI macOS image lacks CGO or Wails

## Related ADR / contracts

- `docs/decisions/ADR-0005-golang-toolchain.md`
- `docs/operations/development.md`

## Handoff

Exact next step: owner decides whether to add a workflow file.
