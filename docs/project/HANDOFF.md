# Handoff

Timestamp: 2026-11-19 19:19 UTC
Agent: agent (RD-062 DMG + GitHub Release; macos-15 runner)
Task ID: RD-062
Branch / commit: working tree (agent MUST NOT commit)

## Done

- `build-macos.yml`: `darwin:package:dmg` → artifact `RouteDeck.dmg`
- On `v*` tag **push** only: published GitHub Release with DMG only (`softprops/action-gh-release@v2`)
- `contents: write`; `main` / `workflow_dispatch` skip the Release step (`if: push && refs/tags/v*`)
- Runner `macos-14` → `macos-15` (GitHub retired macos-14; annotation macos-14-arm64-us-default)
- Docs: `release.md`, `build.md`, ADR-0005 consequences

## Verification

| Command | Result |
| --- | --- |
| `golangci-lint run ./...` | PASS (0 issues) |
| workflow inspect | PASS |
| tag → Release smoke | NOT RUN (owner) |

## Exact next step

Owner commits, then `git tag v0.0.2 && git push origin v0.0.2` (or next free tag) and open the GitHub Release for the DMG.

## Suggested commit message

```text
publish github release with dmg on version tags (rd-062)
```

## Context to read

- `docs/project/tasks/RD-062.md`
- `.github/workflows/build-macos.yml`
- `docs/operations/release.md`

## Notes

- Owner decisions: DMG only; published immediately; no GoReleaser.
- Ad-hoc codesign / Gatekeeper blurb is in the Release body template.
