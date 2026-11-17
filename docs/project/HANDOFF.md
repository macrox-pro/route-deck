# Handoff

Timestamp: 2026-11-17 19:55 UTC
Agent: agent (RD-057 CI macOS package)
Task ID: RD-057
Branch / commit: working tree (agent MUST NOT commit)

## Done

- Added `.github/workflows/build-macos.yml`: `macos-14`, `wails3@v3.0.0`, `wails3 task darwin:package`, zip + upload artifact
- Triggers: `workflow_dispatch`, push `main`, tags `v*`
- No Developer ID / notarization / cert secrets
- Documented in `docs/operations/build.md`, `docs/operations/release.md`, ADR-0005
- Left `.github/workflows/ci.yml` as the PR check gate (RD-007)

## Verification

| Command | Result |
| --- | --- |
| `test -f .github/workflows/build-macos.yml` | PASS |
| `golangci-lint run ./...` | PASS (0 issues; cache write warnings in sandbox) |
| Remote Actions `build-macos` | NOT RUN (owner push) |

## Exact next step

Owner: commit, push, then Actions → `build-macos` → download `RouteDeck-macos-<sha>` → unzip → Right-click Open (Gatekeeper).

## Suggested commit message

```text
add unsigned macos package workflow for ci smoke (rd-057)
```

## Context to read

- `docs/project/tasks/RD-057.md`
- `.github/workflows/build-macos.yml`
- `docs/operations/build.md`
- `docs/operations/release.md`

## Notes

- Artifact is arm64-only on `macos-14`; universal/DMG/notarize are out of scope.
- Ad-hoc codesign comes from `build/darwin/Taskfile.yml` (`codesign --sign -`).
