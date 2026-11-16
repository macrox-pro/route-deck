# Build

## Local

```bash
wails3 task build
```

`Taskfile.yml` dispatches to `build/<GOOS>/Taskfile.yml`. Product v1 packages macOS only:

```bash
wails3 task darwin:package
```

Confirm the task name with `wails3 task --list` if this fails. Do not invent a second build system.

Dev loop:

```bash
wails3 dev
```

Config file: `build/config.yml`.

## Identity (RD-002)

`build/config.yml` still has template company and product strings until the owner updates:

| Field | Target |
| --- | --- |
| companyName | MACROX |
| productName | RouteDeck |
| productIdentifier | `pro.macrox.routedeck` |
| description | Per-app proxy policy for macOS |

Regenerating assets from `build/config.yml` overwrites generated plists. Follow the comment in that file: `wails3 task common:update:build-assets` only when the owner intends to refresh assets.

## Icon

Master: `docs/branding/route-deck-icon.png`.

1. Copy it over `build/appicon.png` (square PNG).
2. Regenerate `.icns`, `.ico`, and `Assets.car` with the Wails v3.0.0 asset task above.
3. Do not treat `build/` outputs as the source.

This docs pass does not run that regeneration.

## Minimum OS

Product floor is macOS 13. The generated plist may say a lower `LSMinimumSystemVersion` until RD-002. The product requirement wins over the template default.
