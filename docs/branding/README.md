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

## Install into the Wails build (owner task RD-002)

1. Replace `build/appicon.png` with this master (or a square export of it).
2. Regenerate platform icons with the Wails v3.0.0 task flow (`icons.icns`, Windows `.ico`, `Assets.car` / `.icon` as that version documents).
3. Keep this file as the source. Treat `build/` icon outputs as derived.
4. Confirm Dock and Finder show RouteDeck, not the Wails placeholder.

This docs pass does not regenerate `.icns` or `.ico`.
