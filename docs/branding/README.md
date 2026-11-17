# Brand

Canonical app icon: `route-deck-icon.png` in this folder.

Motif: layered glass “deck” tiles and a white Y-shaped split-route arrow on a deep blue field with magenta and cyan edges.

## Palette (cues, not a full design system)

| Role | Cue |
| --- | --- |
| Field | Deep blue |
| Accent | Cyan, magenta |
| Mark | White |

Dark-friendly UI is acceptable. Do not recolor the master into a generic purple gradient.

Untitled UI `frontend/src/styles/theme.css` may still use the CLI default purple `--color-brand-*` scale. Recoloring that scale to the RouteDeck palette is a separate branding task.

## Usage

| Surface | Rule |
| --- | --- |
| App icon / Dock / DMG | Derive from this master |
| README / CONCEPTION | Link this file |
| Tray (later) | May need a simplified mark; do not invent it here |

## Do not

- Redraw or replace the motif without owner approval
- Treat Wails template icons under `build/` as the brand master
- Commit a second “official” PNG elsewhere

## Install into the Wails build

1. Replace `build/appicon.png` with this master (square PNG; this file is already 1024×1024).
2. Regenerate `.icns` and `.ico` with `wails3 task common:generate:icons` (PNG only). Do not add an Icon Composer `.icon` or ship `Assets.car` for v1 (macOS 13 floor).
3. Keep this file as the source. Treat `build/` icon outputs as derived.
4. Confirm Dock and Finder show RouteDeck on `bin/RouteDeck.app`, not the Wails placeholder.
