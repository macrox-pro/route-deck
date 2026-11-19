# Release

v1 ships a macOS `.app` / `.dmg` built with Wails v3.0.0. No Windows or Linux product build.

## Before a release

- `go test ./...` PASS
- `golangci-lint run ./...` PASS
- `npm run build` in `frontend/` PASS
- `wails3 task darwin:package:dmg` (or `darwin:package`) succeeds
- Supported macOS: 13+
- Supported ProxyBridge version range recorded from RD-005
- Release notes list fail-open limits: no kill switch, UDP vs HTTP, restart gap

## CI package and GitHub Release (no notarization)

`.github/workflows/build-macos.yml` runs `wails3 task darwin:package:dmg` with `APP_VERSION` / `APP_COMMIT` derived from the git tag (`v*` → semver without `v`, else `dev`) and short SHA, then uploads `RouteDeck.dmg` as a workflow artifact.

On tag pushes matching `v*`, the same job creates a **published** GitHub Release for that tag and attaches **only** `RouteDeck.dmg` (no `.app` zip). Release notes are generated from commits plus a short ad-hoc / Gatekeeper blurb.

That build uses Taskfile ad-hoc codesign only. It is for smoke and internal testers, not Gatekeeper-clean public distribution. Prefer release tags of the form `vMAJOR.MINOR.PATCH`.

```bash
git tag v0.0.2
git push origin v0.0.2
```

## Signing

If a Developer ID cert is present, sign and notarize using Wails v3 docs for the installed version (local or a future CI job). If not, the release notes MUST say the build is unsigned / ad-hoc only and what Gatekeeper will do (Right-click → Open, Keychain re-prompts on rebuild).

## Install notes for users

RouteDeck does not embed ProxyBridge. The user installs ProxyBridge from its official source and approves the system extension.

## Rollback

Keep the previous `.app` and the Application Support database. See `recovery.md` before replacing a release that changed GORM models.
