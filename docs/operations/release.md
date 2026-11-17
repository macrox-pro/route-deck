# Release

v1 ships a macOS `.app` built with Wails v3.0.0. No Windows or Linux product build.

## Before a release

- `go test ./...` PASS
- `golangci-lint run ./...` PASS
- `npm run build` in `frontend/` PASS
- `wails3 task darwin:package` (or the task `wails3 task --list` shows) succeeds
- Supported macOS: 13+
- Supported ProxyBridge version range recorded from RD-005
- Release notes list fail-open limits: no kill switch, UDP vs HTTP, restart gap

## CI package (no notarization)

`.github/workflows/build-macos.yml` runs `wails3 task darwin:package` with `APP_VERSION` / `APP_COMMIT` derived from the git tag (`v*` → semver without `v`, else `dev`) and short SHA, then uploads `RouteDeck.app.zip`. That build uses Taskfile ad-hoc codesign only. It is for smoke and internal testers, not Gatekeeper-clean public distribution. Prefer release tags of the form `vMAJOR.MINOR.PATCH`.

## Signing

If a Developer ID cert is present, sign and notarize using Wails v3 docs for the installed version (local or a future CI job). If not, the release notes MUST say the build is unsigned / ad-hoc only and what Gatekeeper will do (Right-click → Open, Keychain re-prompts on rebuild).

## Install notes for users

RouteDeck does not embed ProxyBridge. The user installs ProxyBridge from its official source and approves the system extension.

## Rollback

Keep the previous `.app` and the Application Support database. See `recovery.md` before replacing a release that changed GORM models.
