# ADR-0013: Setup guides, client detection, and one-click proxy add

Status: Accepted
Date: 2026-11-17
Owners: MACROX RouteDeck

## Context

Users need a local SOCKS5 or HTTP CONNECT endpoint from a third-party VPN client before RouteDeck can protect apps. CONCEPTION lists onboarding when ProxyBridge is missing. The owner closed product decisions for a Setup guides screen, detection of popular clients, one-click proxy add, a suggestion banner, and a first-run checklist (2026-11-17 plan approval).

## Decision

- Ship a **Setup guides** screen (`/help`) with English step cards. Go owns all guide content and detection; React renders only.
- Guide inventory (v1): ProxyBridge (always first), Shadowrocket, v2rayN, Hiddify, Happ, sing-box SFM, Surge Mac, Generic (always last). Clash Verge Rev and ClashX Meta / FlClash are out of scope.
- Detection: match installed apps by **bundle ID** via `InstalledApplicationLister` (ADR-0009). Probe documented default ports on `127.0.0.1` only (TCP dial). "Port open" is a TCP fact, never a health or protection word.
- Bundle ID sources (must be official; else Detection unavailable, no name matching):

| Guide | Bundle ID | Source |
| --- | --- | --- |
| ProxyBridge | `com.interceptsuite.ProxyBridge` | Verified on owner Mac |
| Shadowrocket | `com.liguangming.Shadowrocket` | Verified on owner Mac |
| v2rayN | `2dust.v2rayN` | Official `package-osx.sh` Info.plist |
| Hiddify | (none) | Detection unavailable until Mac app ID is sourced |
| Happ | `com.happ` | Homebrew cask preferences domain |
| sing-box SFM | `io.nekohasekai.sfamt.standalone` | Homebrew cask `sfm` uninstall quit |
| Surge Mac | `com.nssurge.surge-mac` | Homebrew cask preferences / helper paths |
| Generic | (none) | No detection |

- Default ports: Shadowrocket `1082` mixed; v2rayN `10808` mixed; Hiddify `12334` mixed; Happ `10808` SOCKS5; Surge `6153` SOCKS5 (also note `6152` HTTP); SFM and Generic have no default port (no probe).
- Mixed port is added as **SOCKS5**. Split-port clients use the SOCKS5 port.
- **One-click AddProxy** for guides with a default port: save via `Catalog.SaveProxy`, then `HealthMonitor.Recheck`. Closed port still saves; UI shows real health (`Unavailable`) and a "Connect client first" hint. SFM and Generic open the Proxies form with `127.0.0.1` + SOCKS5 pre-filled.
- **Suggestion banner** on Dashboard and Proxies when installed + port open + endpoint not in catalog. Dismiss is in-memory for the process lifetime (no SQLite).
- Every client guide includes a VPN/TUN warning and "Do not add this client itself to a RouteDeck Proxy rule" (warning only; no Go reject in this ADR).
- Prober accepts only catalog ports and dials only `127.0.0.1`. Installed-app list is cached in the use case (refresh on start and Help Refresh).
- First-run checklist: ProxyBridge found, proxy Available, app added, profile applied. Auto-hides when complete; reappears on regression; no stored flag.
- Docs links open via Wails `Browser.OpenURL` with URLs from the Go catalog by guide ID.

## Alternatives

- Prefill form only (no one-click save). Rejected: owner chose one-click add.
- Persist dismiss in SQLite. Rejected: session-only dismiss.
- Read client config files for the real port. Rejected: fragile; out of scope.
- Change VPN client settings. Rejected: CONCEPTION non-goal.
- Match by display name when bundle ID missing. Rejected: owner chose Detection unavailable.

## Consequences

- New Wails service `SetupGuidesService` and application use case.
- New ADR term "Port open" in the glossary.
- RD-040..RD-043 implement detection, one-click add, banner, and checklist.

## Follow-up

Owner interactive RD-033 Apply→Active. Optional later ADR for rejecting VPN clients in Proxy policies.

## Open decisions (owner)

None.
