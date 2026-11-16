# Lifecycle

Wails process lifetime and ProxyBridge CLI lifetime are different. Stopping one does not imply the other.

## Wails

Current scaffold sets `ApplicationShouldTerminateAfterLastWindowClosed: true` in `main.go`. That matches a demo window, not the product.

When the tray phase starts, that flag MUST become `false`: closing the window keeps RouteDeck running. Quit is a separate action that stops the CLI and warns the user.

Do not change `main.go` until the tray task is READY.

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

## Autostart (later)

Launching RouteDeck at login does not prove the VPN client is up. If the proxy is down, show degradation. MUST NOT flip policies to `Direct`.
