# ADR-0018: Quit warns and can save dirty profile/proxy editors

Status: Accepted
Date: 2026-11-19
Owners: MACROX RouteDeck

## Context

ADR-0014 and lifecycle document that quitting with a dirty editor discards edits. Owner closed (2026-11-19): warn, and offer Save and Quit when a save already exists.

Profile and proxy editors save from React via `ProfilesService.Save` / `ProxiesService.Save`. The helpers dialog confirms identities and has no Save. Multiple editor windows may be open together.

## Decision

- When at least one **profile or proxy** editor is dirty: Quit Question body says unsaved edits will be lost. Buttons: **Cancel**, **Quit** (discard, then existing stop-CLI-and-quit), **Save and Quit**.
- Save and Quit emits one event; each dirty profile/proxy dialog runs the same function as its Save button. Quit proceeds only after every such save succeeds. The first save failure leaves the app running and surfaces the existing save error.
- A dirty **helpers** window gets the warning and **Cancel / Quit** only (no Save and Quit).
- If no dirty editor is open, buttons stay Cancel / Quit.
- Appearance chrome discard is unchanged.

This **supersedes** the ADR-0014 / lifecycle known limitation that quitting with a dirty editor always discards edits without a save path.

## Alternatives

- Block Quit until Save/Discard in the editor. Rejected: owner chose warn + optional Save and Quit.

## Consequences

- Tray `RequestQuit` and editor dialogs cooperate via a new event.
- RD-065 implements the change. Lifecycle docs update.

## Follow-up

None required for v1.

## Open decisions (owner)

None.
