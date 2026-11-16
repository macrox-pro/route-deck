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

## Signing

If a Developer ID cert is present, sign and notarize using Wails v3 docs for the installed version. If not, the release notes MUST say the build is unsigned and what Gatekeeper will do.

## Install notes for users

RouteDeck does not embed ProxyBridge. The user installs ProxyBridge from its official source and approves the system extension.

## Rollback

Keep the previous `.app` and the Application Support database. See `recovery.md` before replacing a release that changed GORM models.
