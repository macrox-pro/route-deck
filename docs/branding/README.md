# Brand

Canonical app icon: `route-deck-icon.png` in this folder.

README hero banner: `readme-banner.jpg` (1024×512). Use only as the GitHub README header; do not treat it as the Dock / DMG icon master.

Motif: layered glass “deck” tiles and a white Y-shaped split-route arrow on a deep blue field with magenta and cyan edges.

## Palette (cues, not a full design system)

| Role | Cue |
| --- | --- |
| Field | Deep blue |
| Accent | Cyan, magenta |
| Mark | White |

Dark-friendly UI is acceptable. Do not recolor the master into a generic purple gradient.

Untitled UI `frontend/src/styles/theme.css` uses a burgundy `--color-brand-*` scale for UI chrome (ADR-0020 / DESIGN.md): brand-600 `#6E2432`, dark solid brand-500 `#8A3A48`. The app mark palette (deep blue / cyan / magenta) is unchanged. Do not treat burgundy as a protection status color.

## Usage

| Surface | Rule |
| --- | --- |
| App icon / Dock / DMG | Derive from `route-deck-icon.png` |
| README hero | `readme-banner.jpg` |
| README / CONCEPTION | Link this folder |
| Tray (later) | May need a simplified mark; do not invent it here |

## Do not

- Redraw or replace the motif without owner approval
- Treat Wails template icons under `build/` as the brand master
- Commit a second “official” app-icon PNG elsewhere (the README banner is a separate asset)

## Install into the Wails build

1. Replace `build/appicon.png` with this master (square PNG; this file is already 1024×1024).
2. Regenerate `.icns` and `.ico` with `wails3 task common:generate:icons` (PNG only). Do not add an Icon Composer `.icon` or ship `Assets.car` for v1 (macOS 13 floor).
3. Keep this file as the source. Treat `build/` icon outputs as derived.
4. Confirm Dock and Finder show RouteDeck on `bin/RouteDeck.app`, not the Wails placeholder.
