# Build

## Local

```bash
wails3 task build
```

`Taskfile.yml` dispatches to `build/<GOOS>/Taskfile.yml`. Product v1 packages macOS only:

```bash
wails3 task darwin:package
```

That produces `bin/RouteDeck.app` with executable `Contents/MacOS/route-deck` (`APP_BUNDLE_NAME` vs `APP_NAME` in `Taskfile.yml`). Confirm the task name with `wails3 task --list` if this fails. Do not invent a second build system.

## Version identity

Link-time package: `internal/buildmeta` (`Version`, `Commit`).

```bash
wails3 task darwin:package APP_VERSION=1.2.3 APP_COMMIT=$(git rev-parse --short=7 HEAD)
```

- Default without vars: `Version=dev`, empty `Commit` (local / `wails3 dev`)
- Git tags: `v1.2.3` → pass `APP_VERSION=1.2.3` (no leading `v`)
- UI sidebar / Settings label: `v1.2.3` or `dev` via `buildmeta.DisplayVersion`
- `Info.plist` `CFBundleShortVersionString` / `CFBundleVersion`: numeric form, or `0.0.0` when `dev`

## GitHub Actions (unsigned package)

Workflow: `.github/workflows/build-macos.yml`.

- Triggers: `workflow_dispatch`, push to `main`, tags `v*`
- Runner: `macos-14` (arm64)
- Command: `wails3 task darwin:package APP_VERSION=… APP_COMMIT=…` with `CGO_ENABLED=1` and `GOFLAGS=-mod=vendor`
- Tag `v*` → version without `v`; otherwise `dev`. Commit is `github.sha` short (7)
- Artifact: `RouteDeck.app.zip` (ad-hoc codesign from the Taskfile only)
- MUST NOT require Developer ID or notarization in this workflow
- PR check gate stays `.github/workflows/ci.yml` (tests, lint, frontend build only)

Dev loop:

```bash
wails3 dev
```

Config file: `build/config.yml`.

## Identity

| Field | Value |
| --- | --- |
| companyName | MACROX TECHNOLOGIES |
| productName | RouteDeck |
| productIdentifier | `pro.macrox.routedeck` |
| description | Per-app proxy policy for macOS |
| copyright | © 2026, MACROX TECHNOLOGIES |
| version | 0.1.0 |
| Binary / `APP_NAME` | `route-deck` |
| macOS bundle / `APP_BUNDLE_NAME` | `RouteDeck` → `bin/RouteDeck.app` |

Regenerating assets from `build/config.yml` overwrites generated plists. Follow the comment in that file: `wails3 task common:update:build-assets` only when the owner intends to refresh assets. After that refresh, re-apply `LSMinimumSystemVersion` `13.0.0` on the darwin plists (Wails v3.0.0 templates hardcode `12.0.0`) and remove `CFBundleIconName` when no `Assets.car` is shipped.

## Icon

Master: `docs/branding/route-deck-icon.png`.

1. Copy it over `build/appicon.png` (square PNG).
2. Regenerate `.icns` and `.ico` with `wails3 task common:generate:icons` (PNG only; no Icon Composer `.icon` / `Assets.car` — that path targets macOS 26).
3. Do not treat `build/` icon outputs as the source.

## Minimum OS

Product floor is macOS 13. Darwin `Info.plist` / `Info.dev.plist` must say `LSMinimumSystemVersion` `13.0.0`. The Go linker may still warn about link minimum 12.0; there is no `build/config.yml` key for that.
