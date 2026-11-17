# Lifecycle

Wails process lifetime and ProxyBridge CLI lifetime are different. Stopping one does not imply the other.

## Wails

`ApplicationShouldTerminateAfterLastWindowClosed` is `false` (ADR-0010). Closing the window Cancels destroy and Hides it. RouteDeck stays in the menu bar. The Dock icon stays for the process lifetime (`ActivationPolicyRegular`).

Quit is a separate action (tray Quit, Cmd+Q, Dock Quit). Every Quit path shows a native Question dialog before exit is accepted. Cancel keeps the process resident. Confirm with a managed CLI stops that CLI gracefully, then quits. `OnShutdown` still stops the managed CLI as a safety net.

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

Launch at Login is controlled from Settings via Wails `app.Autostart`. OS registration is the source of truth. Launching RouteDeck at login does not prove the VPN client is up and MUST NOT auto-Apply. If the proxy is down, show degradation. MUST NOT flip policies to `Direct`.
