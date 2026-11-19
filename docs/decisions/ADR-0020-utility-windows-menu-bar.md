# ADR-0020: Utility windows, menu bar, and design-system UX

Status: Accepted
Date: 2026-11-20
Owners: MACROX RouteDeck

## Context

The frontend still used Untitled purple brand tokens, a `/help` Setup guides page inside the main window, `window.confirm` for some deletes, sparse empty states, and the Wails default application menu (Reload, Force Reload, Help → wails.io). Owner closed product decisions for a DESIGN.md-aligned desktop refactor (2026-11-20 plan): burgundy accent, shared page layout, Setup Guides and About as closable native windows, progress-step UI for checklist, native delete confirms, editor form grouping, and a real RouteDeck menu bar with accelerators. Editor modality (ADR-0014) stays.

## Decision

### Utility windows

- Kinds `guides` and `about` are **utility** windows (native frame, no Frameless). They do **not** disable the main window, have no dirty tracking, and are not part of Save and Quit (ADR-0018).
- One reused window per kind. Appearance recreate and Quit force-close them.
- Guides: title "Setup Guides", 880×640, resizable, min 720×520, hash `/#/dialog/guides` or `/#/dialog/guides/:step`.
- About: title "About RouteDeck", 420×460, not resizable, hash `/#/dialog/about`.
- Bound service: `OpenGuides(step)`, `OpenAbout()`, `Close(kind)`, `ShowMain(route)`. Closed `step` and `route` lists live in the wails-api contract.
- This **supersedes** ADR-0013 “Ship a Setup guides screen (`/help`)”. Guide content, detection, one-click add, checklist DTO, and ADR-0016 dismiss rules are unchanged. The UI surface moves to the `guides` window. The main-window `/help` route is removed.

### main.command

- Event `main.command` carries a closed `command` field: `navigate` (optional `route`), `apply`, `stop`, `recheck`, `toggle-sidebar`, `add-application`, inside the existing `EventEnvelope`.
- `app-shell` is the sole listener. `apply` / `stop` / `recheck` navigate to Dashboard with `location.state`; `add-application` navigates to Applications with `location.state`. Pages run existing handlers and reload snapshots. The event never carries protection state.
- `ShowMain` with any editor open focuses the editor and does **not** emit (ADR-0014 modality).

### Tray (ADR-0010 amend)

Tray menu closed list becomes:

1. Status (disabled) — unchanged
2. Open RouteDeck — unchanged visibility rule
3. Stop Protection — same enable rule as ADR-0010 (`CLIOwned`); Title Case label
4. Separator
5. Setup Guides…
6. About RouteDeck
7. Separator
8. Quit RouteDeck

This **supersedes** the ADR-0010 tray closed list that had only Open / Stop protection / Quit, and the exact label “Stop protection”.

### Menu bar

- Replace the Wails default application menu with a RouteDeck menu (RouteDeck, File, Edit, Protection, View, Window, Help). Accelerators and enable rules are listed in `docs/engineering/desktop-ux.md` and the wails-api contract.
- About opens the About utility window (not `Menu.ShowAbout`). Help → GitHub uses `Browser.OpenURL` (never navigates the app webview).
- Protection, View navigation items, and Add Application are **disabled** while any editor window is open. New Profile, New Proxy, About, Setup Guides, and Settings stay enabled. Stop follows `CLIOwned`.
- `Options.Name` stays `route-deck` (autostart fallback label). Menu titles use `SetLabel`.
- Not included: Reload, Force Reload, release DevTools, wails.io Learn More.

### Delete confirmation

- Proxies and applications confirm delete with native `Dialog.Question`, Cancel default, Cancel returns nil without deleting. Profiles already do (ADR-0019). Shared helper `confirmDelete` in the wails adapter.
- MUST NOT use `window.confirm` in the frontend.

### Helpers Select all (ADR-0014 amend)

- Helpers toolbar uses one **tri-state** “Select all” checkbox (`isIndeterminate` when partially selected) instead of separate Select all / Deselect all buttons. This **supersedes** the ADR-0014 sentence “UI offers Select all / Deselect all”.

### Guide open path (desktop-ux amend)

- Opening a proxy editor from a guide uses `EditorDialogs.OpenProxy("", guideId)` directly. This **supersedes** “Guide prefill via query stays” for the guides window. The `/proxies?guide=` query path is removed with `/help`.

### Visual language

- UI accent is burgundy (`--color-brand-*` scale; 600 `#6E2432`; dark solid brand-500 `#8A3A48`). App mark palette is unchanged.
- Applied badge is gray `modern`; green is reserved for Active. Unknown / Failed / Unprotected MUST NOT look like Active.
- Menu bar and tray use macOS Title Case. Buttons, labels, and empty states use sentence case. Sidebar labels stay title-case nouns.
- Menu-bar accelerators are in scope for desktop UX. Keyboard **list** navigation beyond Enter / Esc / Cmd-W remains out of scope.

## Alternatives

- Keep `/help` in the main window alongside a guides window. Rejected: owner chose one closable window and remove `/help`.
- Minimal menu (About + Settings only). Rejected: owner chose the full menu.
- Make editors non-modal. Rejected: owner kept ADR-0014 modality.
- Change `Options.Name` to `RouteDeck`. Rejected: autostart and role labels depend on the slug; use `SetLabel` instead.

## Consequences

- New types: `UtilityWindows`, `UtilityWindowsService` (file `utilities.go` — not `utility_windows.go`, which Go treats as Windows-only), `app_menu.go`, `confirm_delete.go`, `main.command` event, page-layout and progress-steps frontend components.
- Contracts and `DESIGN.md` updated. Frontend gate remains `npm run build` (no ESLint).
- Implemented across RD-070..RD-076.

## Follow-up

None required beyond RD-070..RD-076.

## Open decisions (owner)

None.
