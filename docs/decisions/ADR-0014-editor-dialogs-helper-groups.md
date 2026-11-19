# ADR-0014: Editor dialogs, helper groups, sync confirm, Dock badge

Status: Accepted
Date: 2026-11-17
Owners: MACROX RouteDeck

## Context

Inline Profile and Proxy forms on the main window are hard to use. ChatGPT-style apps nest helpers under `Contents/Frameworks/.../Helpers` and `Contents/Resources`, share one bundle id across several executables, and leave Sparkle updaters that must not be auto-selected. ConfirmHelpers was additive (uncheck did not remove). Grouping by bundle id shipped in the UI (`afb231c`) but still lived in TypeScript. Owner closed product decisions for custom dialog windows, helper groups from Go, sync confirm, discovery scope, and a Dock attention badge (2026-11-17).

Supersedes the helper-discovery depth and ConfirmHelpers additive behaviour of [ADR-0008](ADR-0008-vertical-slice-decisions.md). Intake still uses the picker / browse + confirm; unconfirmed helpers never enter rules.

## Decision

### Editor dialogs

Profile create/edit, Proxy create/edit, and Review helpers each open in a separate Wails window (kinds: `profile`, `proxy`, `helpers`). One reused window per kind. Native frame (no Frameless). Parent main window is disabled with `SetEnabled` on every OS; `AttachModal` is also called (no-op on Linux). Dirty close asks `Dialog.Question` ("Discard changes?"). Bindings return immediately; dialogs save through existing services and emit `catalog.changed` so the main window refetches. Hash routes: `/#/dialog/profile/:id?`, `/#/dialog/proxy/:id?`, `/#/dialog/helpers/:appId`.

### Helper groups and matching

ProxyBridge matches bundle id **or** executable name; path never matches. One `ProcessIdentity` per bundle id per application. Go builds `HelperGroupDTO` (inspect executable names, confirmed, same-developer / updater hints). React renders groups and MUST NOT re-group. Confirm syncs: unchecked groups delete every identity with that bundle id; if any profile references any of them, return `ErrReferenced` and delete nothing. Repeat Confirm refreshes the stored path and executable name to the inspect representative. Checkbox defaults: only already-confirmed groups start checked; new groups start unchecked. UI offers Select all / Deselect all. **Superseded (ADR-0020):** one tri-state "Select all" checkbox (`isIndeterminate` when partially selected) replaces separate Select all / Deselect all buttons. Same-developer and Updater remain hints only (not auto-check). Utility windows (`guides`, `about`) are a separate kind set and are not editor dialogs.

### Compile

A unique bundle id still compiles to bundle id plus that identity's executable name. When several stored identities share a bundle id and the same action, `BuildRulePlan` emits one rule whose `processNames` is only that bundle id. Different actions on a shared matchable name return `ErrConflictingMatch`. Do not migrate old duplicate rows in SQLite; collapse at compile time. Never map either error to Direct. **Match-conflict hard-fail on Save/Apply superseded by [ADR-0019](ADR-0019-profile-delete-last-wins.md)** (last policy in list wins; persist resolved policies; never map to Direct).

### Discovery (macOS adapter only)

Walk every nested `.app` with an executable, including `Contents/Resources`. Dedupe by `EvalSymlinks`. Prefer `Versions/Current` in stored paths. Set `IsUpdater` for Sparkle-style helpers. No `.xpc` / `.appex` in this ADR. Domain and application layers MUST NOT mention `.app`, plist keys, or framework paths.

### Dock badge

`pkg/services/dock` badge `!` when a profile is applied and ProtectionState is not Active. No badge when nothing is applied. No notifications in this ADR. Unsupported platforms: Debug once and continue.

## Alternatives

- Frameless custom chrome. Rejected: accessibility and platform conventions.
- One ProcessIdentity per file. Rejected: CLI matching makes per-file actions on a shared bundle id ambiguous.
- Auto-confirm helpers. Rejected: ADR-0008 confirm rule stands.
- SQLite migration of duplicate rows. Rejected: compile-time collapse is enough.
- Notifications for protection loss. Deferred: needs signed app authorization on macOS.

## Consequences

- New Wails `EditorDialogsService`, event `catalog.changed`, Dock inbound adapter.
- Frontend dialog routes outside `AppShell`; Profiles and Proxies pages become list-only.
- Tasks RD-049..RD-056.

## Follow-up

Windows/Linux `ApplicationInspector` adapters. Optional notifications. Remove leftover duplicate identities only if an owner task asks.
