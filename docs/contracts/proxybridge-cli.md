# ProxyBridge CLI contract

Source checked: 2026-11-16 against https://proxybridge.interceptsuite.com/docs/macos-cli/

Re-read that page before changing the adapter. MUST NOT add flags that are not on that page.

## What the CLI is

`proxybridge` is a launcher, not a separate daemon. The installer places it at `/usr/local/bin/proxybridge`. It runs:

`/Applications/ProxyBridge.app/Contents/MacOS/ProxyBridge`

with CLI options. No `sudo`. The GUI and the CLI share one system extension. Only one may control it.

## Options (documented)

Each option below is on https://proxybridge.interceptsuite.com/docs/macos-cli/.

| Option | Behavior |
| --- | --- |
| `--version` | Print version and exit |
| `--help`, `-h` | Help and exit |
| `--list-profiles` | List saved profiles and exit |
| `--profile <name>` | Run a profile saved in the ProxyBridge app |
| `--profile -` | Read a macOS `.pbprofile` JSON from stdin |
| `--verbose <0-3>` | 0 silent (default), 1 logs, 2 connections, 3 both |
| `--update` | Check for an update and exit. It does not download or install |

`--profile=Work` and `--verbose=2` are accepted. `--cli` is accepted and ignored.

`--profile` is required to run unless the command is `--list-profiles`, `--update`, `--version`, or `--help`/`-h`.

MUST NOT pass a filesystem path to `--profile`. The app is sandboxed and refuses paths. Pipe the profile: `proxybridge --profile -`.

A profile name that contains `/` or ends in `.pbprofile` is treated as a path and refused. Rename such a profile in the ProxyBridge app.

Stdin limit: 4 MiB. Stdin must be a pipe or file, not a terminal.

## How RouteDeck uses it

| Action | Command shape |
| --- | --- |
| Version | `proxybridge --version` |
| Run RouteDeck profile | `proxybridge --profile - --verbose <n>` with JSON on stdin |
| Update check | Only when the user asks. `--update` never installs |

RouteDeck does not rely on profiles saved inside the ProxyBridge GUI.

## Discovery

Documented locations (https://proxybridge.interceptsuite.com/docs/macos-cli/):

| Path | Role |
| --- | --- |
| `/usr/local/bin/proxybridge` | Installer launcher |
| `/Applications/ProxyBridge.app/Contents/MacOS/ProxyBridge` | App executable |

Also honor a user-configured absolute path. Unknown version: warning, not silent success. Signature check: only if a READY task specifies a concrete API. Do not invent a codesign policy here.

Other well-known paths (Homebrew, and so on) are **not documented** on the CLI page. Do not guess them. An owner-configured absolute path is the only addition without a docs URL.

Do not copy sample version strings from the docs (for example `4.1.0`) as the supported version. Record the installed binary's `--version` output when probing.

## Process control

One managed process. If the GUI is running, the CLI refuses to start (`The ProxyBridge app is running...`). If another CLI is already running, a second one is refused (`Another ProxyBridge CLI is already running.`). If the GUI starts later, the tunnel stops and the CLI exits 1 (`The tunnel was stopped by something else...`).

On that conflict: set BridgeState `Conflict`, ProtectionState not protected, show a native Warning dialog telling the user to quit the GUI, do not restart in a loop.

Stop: SIGINT/SIGTERM/SIGHUP clear rules (docs: Ctrl+C, SIGTERM, SIGHUP). If the extension does not answer, the CLI exits after 5 seconds. `SIGKILL` does not clear rules. After a kill -9, RouteDeck MUST NOT report protected.

First-run extension approval prints `Extension: waiting for your approval...` (classified `extension_approval`, not a failure). RouteDeck shows a native Info dialog once per Apply start and keeps waiting. Approval can wait up to 5 minutes before `Timed out waiting for the extension...`.

Graceful stop is SIGTERM, wait up to ~20s (CLI may need ~5s for extension cleanup), then SIGKILL with up to 3 kill+reap attempts (~5s each, 1s between). Start owns the single `cmd.Wait`; Stop only observes that exit channel. Exit after our stop signal counts as success even when `Wait` returns a non-zero status. `os.ErrProcessDone` on signal is success once reaped. MUST NOT surface a stop timeout to the UI when the process was killed and reaped. MUST NOT use context cancel as the only stop (it can kill immediately); kill retries ignore caller cancel. After SIGKILL, RouteDeck MUST NOT report protected.

## Readiness

`Extension: active` is not enough. The CLI docs say routing is on when it prints `Redirecting traffic. Press Ctrl+C to stop.` A profile load line looks like `Profile "<name>": N proxies, M rules applied`. Classifier also maps `waiting for your approval` to `extension_approval` (informational only).

RouteDeck confirms apply only after the redirect-ready line (and no fatal error). Partial or unknown lines stay unclassified and MUST NOT flip state to `Active` by themselves.

A warning such as `warning: N rule(s) were rejected by the extension (invalid fields)` stays unclassified. It MUST NOT mean the profile is confirmed or that ProtectionState is `Active`.

Verbose 2 and 3 print connections. v1 Activity does not depend on parsing those lines. If a later task parses them, it needs fixtures. Unknown lines stay `Unclassified`.

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Help, version, list, update check, or clean stop |
| 1 | Bad profile, startup failure, or tunnel lost |
| 2 | Usage error |

Interpret the code together with stderr text and the operation. Code 0 on `--version` is not “protection active”.

## Output

Human text, not JSON. Parser fixtures per CLI version. MUST NOT infer Direct vs Block from a partial string match.

## Unit tests

MUST NOT launch the real binary. Feed fixture strings to the parser and a fake process to the process port.
