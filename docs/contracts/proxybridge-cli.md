# ProxyBridge CLI contract

Source checked: 2026-11-16 against https://proxybridge.interceptsuite.com/docs/macos-cli/

Re-read that page before changing the adapter (RD-005). MUST NOT add flags that are not on that page.

## What the CLI is

`proxybridge` is a launcher, not a separate daemon. The installer places it at `/usr/local/bin/proxybridge`. It runs:

`/Applications/ProxyBridge.app/Contents/MacOS/ProxyBridge`

with CLI options. No `sudo`. The GUI and the CLI share one system extension. Only one may control it.

## Options (documented)

| Option | Behavior |
| --- | --- |
| `--version` | Print version and exit |
| `--help`, `-h` | Help and exit |
| `--list-profiles` | List saved profiles and exit |
| `--profile <name>` | Run a profile saved in the ProxyBridge app |
| `--profile -` | Read a macOS `.pbprofile` JSON from stdin |
| `--verbose <0-3>` | 0 silent, 1 logs, 2 connections, 3 both |
| `--update` | Check for an update and exit. It does not download or install |

`--profile=Work` and `--verbose=2` are accepted. `--cli` is accepted and ignored.

MUST NOT pass a filesystem path to `--profile`. The app is sandboxed and refuses paths. Pipe the profile: `proxybridge --profile -`.

Stdin limit: 4 MiB. Stdin must be a pipe or file, not a terminal.

## How RouteDeck uses it

| Action | Command shape |
| --- | --- |
| Version | `proxybridge --version` |
| Run RouteDeck profile | `proxybridge --profile - --verbose <n>` with JSON on stdin |
| Update check | Only when the user asks. `--update` never installs |

RouteDeck does not rely on profiles saved inside the ProxyBridge GUI.

## Discovery

Documented locations:

| Path | Role |
| --- | --- |
| `/usr/local/bin/proxybridge` | Installer launcher |
| `/Applications/ProxyBridge.app/Contents/MacOS/ProxyBridge` | App executable |

Also honor a user-configured absolute path. Unknown version: warning, not silent success. Signature check: only if a READY task specifies a concrete API. Do not invent a codesign policy here.

Other well-known paths (Homebrew, and so on) are **not documented** on the CLI page. RD-005 may add a path only with a source URL. Until then, do not guess them.

## Process control

One managed process. If the GUI is running, the CLI refuses to start (`The ProxyBridge app is running...`). If the GUI starts later, the tunnel stops and the CLI exits 1 (`The tunnel was stopped by something else...`).

On that conflict: set BridgeState `Conflict`, ProtectionState not protected, tell the user to quit the GUI, do not restart in a loop.

Stop: SIGINT/SIGTERM/SIGHUP clear rules (docs: Ctrl+C, SIGTERM, SIGHUP). `SIGKILL` does not clear rules. After a kill -9, RouteDeck MUST NOT report protected.

Graceful stop is signal, wait, then timeout. MUST NOT use context cancel as the only stop (it can kill immediately).

## Readiness

`Extension: active` is not enough. The CLI docs say routing is on when it prints `Redirecting traffic. Press Ctrl+C to stop.` A profile load line looks like `Profile "<name>": N proxies, M rules applied`.

RouteDeck confirms apply only after the redirect-ready line (and no fatal error). Partial or unknown lines stay unclassified and MUST NOT flip state to `Active` by themselves.

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
