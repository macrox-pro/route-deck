# ADR-0009: Installed application browse list

Status: Accepted
Date: 2026-11-17
Owners: MACROX RouteDeck

## Context

ADR-0008 chose a native `.app` picker for intake. The owner later asked for a browsable list of installed applications (with icons) while keeping Add explicit. Product v1 remains macOS-only, but discovery and icons must sit behind ports so another OS can add an adapter later.

The MacOSX SDK on the development machine has no public `LSCopyAllApplicationURLs`. Private Launch Services SPI is forbidden.

## Decision

- Show a browse list of installed applications. Rows enter the RouteDeck catalog only when the user presses **Add**. Helper confirmation and the native picker stay as in ADR-0008.
- Enumeration unions:
  1. The public application index (`kMDItemContentTypeTree == com.apple.application-bundle` via `NSMetadataQuery`).
  2. Top-level `.app` bundles in `/Applications`, `~/Applications`, and `/System/Applications`.
- Keep only directories ending in `.app`. Drop nested `.app` inside another `.app`, `.service`, and plugins. Nested helpers still come from inspect-on-add.
- Dedupe by bundle id. Preference: path under the user home directory, then `/Applications`, then `/System/Applications`, then any other registered location.
- Each listed row carries `IsSystem` from the platform adapter. On macOS, `IsSystem` is true when the path is outside `/Applications` and `~/Applications` (index hits under `/System`, `/Library`, and similar). The UI filters with this flag only; it MUST NOT hard-code OS paths. Other OS adapters define their own user-facing roots.
- A missing root is skipped. A failed metadata query still returns the three roots. Unreadable or empty-bundle-id plists are skipped at Debug.
- Icons use `-[NSWorkspace iconForFile:]` on the main thread, returned as 32×32 PNG on a separate call. Failed icons leave an empty UI slot.
- Ports: `InstalledApplicationLister` and `ApplicationIconReader`. Only a macOS adapter ships now; `!darwin` stubs return a typed unsupported-platform error.

## Alternatives

- Directory scan of the three roots only. Rejected: misses apps registered outside those folders.
- Private Launch Services `_LSCopyAllApplicationURLs` / `lsregister`. Rejected: private SPI.
- Auto-import every installed app into the catalog. Rejected: catalog stays intentional.
- Embed icon bytes in the list payload. Rejected: hundreds of apps would bloat the first paint.

## Consequences

- ADR-0008’s rejected “scan `/Applications`” alternative is superseded for browse listing; intake into the catalog remains Add + helper confirm + picker.
- macOS adapters may use AppKit / Foundation through cgo or Objective-C bridges verified against the installed SDK.
- UI needs a text filter; virtualization is out of this task.

## Follow-up

Windows/Linux adapters for the same ports. Optional Launch Services refresh notifications. RD-033 live ProxyBridge smoke remains separate.

## Open decisions (owner)

None.
