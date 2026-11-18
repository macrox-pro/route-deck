# Lifecycle

Wails process lifetime and ProxyBridge CLI lifetime are different. Stopping one does not imply the other.

## Wails

`ApplicationShouldTerminateAfterLastWindowClosed` is `false` (ADR-0010). Closing the **main** window Cancels destroy and Hides it. RouteDeck stays in the menu bar. The Dock icon stays for the process lifetime (`ActivationPolicyRegular`).

Editor dialog windows (`profile`, `proxy`, `helpers`) are separate Wails windows (ADR-0014). Closing an editor destroys that window and re-enables the main window. Dirty editors ask before close. Tray "Open RouteDeck" focuses an open editor when present, otherwise shows the main window.

Dock badge `!` appears when a profile is applied and ProtectionState is not Active (ADR-0014). Unsupported platforms log once at Debug.

Quit is a separate action (tray Quit, Cmd+Q, Dock Quit). Every Quit path shows a native Question dialog before exit is accepted. Cancel keeps the process resident. Confirm with a managed CLI stops that CLI gracefully, then quits. `OnShutdown` still stops the managed CLI as a safety net. Quitting with a dirty editor discards edits (known limitation).

## ProxyBridge CLI

Only one component (the process adapter) starts and stops the managed CLI.

Stop sequence:

1. Signal a graceful stop
2. Wait with a timeout
3. Classify timeout or crash separately from a clean exit

MUST NOT use `exec.CommandContext` cancellation as the only stop path. Context cancel may kill the process immediately and look like a crash.

The CLI reads the profile at start. Changing rules means stop then start. There is a gap with no rules.

## Desired vs applied

| Moment | Desired | Applied |
| --- | --- | --- |
| User saves | New | Previous until confirm |
| CLI confirms | New | New |
| Confirm fails | New | None, then one rollback attempt |
| Rollback fails | New | None, `Unprotected`, both errors shown |

## Sleep, wake, network change

Re-check proxy health and bridge state. MUST NOT keep showing `Active` across those events without a re-check.

## Autostart

Launch at Login is controlled from Settings via Wails `app.Autostart`. OS registration is the source of truth. Launch at Login alone does not Apply.

Opt-out **Auto Apply on launch** is a separate SQLite setting (ADR-0015, default on). When enabled and the desired profile has Proxy or Block rules, startup calls the same Apply path as the Dashboard. If the proxy is down, Apply fails honestly and the UI shows Unprotected or Degraded after a successful apply with gaps. MUST NOT flip policies to `Direct`.
