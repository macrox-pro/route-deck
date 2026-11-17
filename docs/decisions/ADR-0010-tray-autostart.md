# ADR-0010: Tray lifecycle and Launch at Login

Status: Accepted
Date: 2026-11-17
Owners: MACROX RouteDeck

## Context

Phase 5 shipped Apply/Stop through the window UI. Closing the last window still terminated the process (`ApplicationShouldTerminateAfterLastWindowClosed: true`). CONCEPTION requires the window close to keep RouteDeck running in the menu bar, and Quit to warn and stop the managed ProxyBridge CLI. ADR-0008 deferred tray, quit warning, and autostart.

## Decision

- `ApplicationShouldTerminateAfterLastWindowClosed` is `false`. Window close Cancels destroy and Hides the window. ProtectionState is unchanged.
- Dock stays for the process lifetime (`ActivationPolicyRegular`). Hiding the Dock icon only while the window is hidden is not a public Wails API. Dock click Shows the window (Wails `ApplicationShouldHandleReopen` plus our Show helper).
- Menu bar uses a monochrome template icon (`SystemTray.SetTemplateIcon`).
- Tray menu (closed list):
  1. Status (disabled) — `ProtectionState`; if `BridgeState != Ready`, append the bridge word
  2. Open RouteDeck — Show + Focus; **hidden while the main window is visible**
  3. Stop protection — always visible; enabled only when RouteDeck owns the CLI (`BridgeState` in `Ready`, `Starting`, `Stopping`)
  4. Separator
  5. Quit RouteDeck
- Apply from tray is forbidden. Apply stays on the Dashboard.
- Tray icon click (left or right) opens the tray menu. Window restore is **Open RouteDeck** (when the window is hidden) or Dock click.
- Quit from tray, Cmd+Q, and Dock Quit share one path. A native Question dialog **always** appears. Cancel stays resident. Confirm with an owned CLI runs graceful `Protection.Stop`, then quit. Confirm without an owned CLI quits without Stop. `Conflict` is not owned; Stop stays disabled; do not signal a foreign process.
- If the main window is hidden when Quit starts, Show it before the dialog when needed so the alert can present.
- `OnShutdown` still stops the managed CLI as a safety net.
- This ADR **supersedes** the ADR-0008 sentence “Quit stops the CLI gracefully on `OnShutdown` with no dialog.”
- Launch at Login is a Settings toggle. OS registration via Wails `app.Autostart` is the source of truth (no SQLite mirror). Login does **not** auto-Apply (ADR-0008).

## Alternatives

- `ActivationPolicyAccessory` to hide the Dock forever. Rejected: removes Dock click restore.
- Persist Launch at Login in SQLite and sync to the OS. Rejected: drifts from OS state.
- Apply from tray. Rejected for v1: Dashboard owns Apply feedback.

## Consequences

- Tray controller lives in `internal/adapter/inbound/tray` (composition root adapter, not a Wails service).
- Settings DTO gains `launchAtLogin` from `Autostart.IsEnabled`.
- Activity viewer is ADR-0011 / RD-038. CI remains RD-007.

## Follow-up

RD-007 CI. Optional multi-state tray icons.

## Open decisions (owner)

None.
