# Desktop UX

English UI. Untitled UI components. Semantic color tokens. See ADR-0014, ADR-0007, ADR-0020, and `DESIGN.md`.

## Editor dialogs

Kinds: `profile`, `proxy`, `helpers`. One reused native-frame window each. No Frameless.

| Kind | Size | Resize |
| --- | --- | --- |
| profile | 760×640 (min 600×480) | yes |
| proxy | 480×560 | no |
| helpers | 600×560 | yes |

Titles: "New Profile", "Edit Profile — {name}", "New Proxy", "Edit Proxy — {name}", "Review Helpers — {app}".

- Footer: macOS/Linux `[Cancel][Save]`; Windows `[Save][Cancel]`. Save is default. Enter saves, Esc cancels, Cmd/Ctrl+W closes.
- Body scrolls; title and footer stay put. Body `px-6 py-5`, fields `gap-5`. Footer `px-6 py-4`.
- Focus first field on open. Inline Go validation errors. Save shows loading.
- Save disabled when name blank, form clean, or save in flight. Create becomes dirty on first edit.
- Dirty close: `Dialog.Question` "Discard changes?".
- Destructive delete defaults to Cancel; Enter does not delete. Proxies and applications confirm delete with the same native path (ADR-0020).
- After save: dialog closes, list refetches, saved row in view, focus returns to opener.
- Theme/dark mode via existing tokens. No white `BackgroundColour` override. Appearance (`system`|`light`|`dark`) is SQLite-backed; native frame uses Wails create-time `Mac.Appearance` (recreate main on change).
- Buttons: `sm` in lists, `md` in dialog footer.
- Button grouping: Action and Protocol are `ButtonGroup`s; Kind on the proxy editor is a `ButtonGroup`; proxy picker is `Select`. Helpers use one tri-state "Select all" checkbox (ADR-0020).

## Utility windows

Kinds: `guides`, `about`. Non-modal native-frame windows. They do not disable the main window, have no dirty tracking, and are not part of Save and Quit. See ADR-0020.

| Kind | Size | Resize |
| --- | --- | --- |
| guides | 880×640 (min 720×520) | yes |
| about | 420×460 | no |

Titles: "Setup Guides", "About RouteDeck". Hash routes: `/#/dialog/guides`, `/#/dialog/guides/:step`, `/#/dialog/about`.

Opening a proxy editor from a guide uses `EditorDialogs.OpenProxy("", guideId)` directly. The old `/help` route and `/proxies?guide=` query path are removed (ADR-0020 supersedes ADR-0013’s `/help` surface and the prior “Guide prefill via query stays” line).

## Profiles

Search filter, "Configured only" toggle, count "N of M" with an action. Helper groups from Go (`HelperGroupDTO`). "Apply app action to helpers" drafts only. Row action opens helpers dialog for that app. Visible row actions: Set desired, Edit. Duplicate and Delete live in the row menu.

## Helpers

One checkbox per bundle-id group. Executable names under the row. Hints: "Same developer", "Updater". Only already-confirmed groups start checked; new groups start unchecked. Toolbar: one tri-state "Select all" checkbox with "N of M selected". Footer preview "Adds N · Removes M". `ErrReferenced` copy names profiles ("Used by Work, Home").

## Proxies

Password never pre-filled. Empty password on edit keeps Keychain secret. Visible row action: Edit. Clear password and Delete live in the row menu. Header: Recheck health, New proxy.

## Lists

Empty states name the next action (approved copy in ADR-0020 / DESIGN.md). Applications catalog shows confirmed helper groups read-only. List pages use shared `PageHeader` / `ListPanel` / `ListRow` / `ScreenState`. Rows are at least `min-h-14`. Lists cap at `max-w-5xl`. Settings at `max-w-3xl`.

## Menu bar and tray

RouteDeck application menu replaces the Wails default (ADR-0020). Tray labels use Title Case; "Stop Protection" enable rule is unchanged (ADR-0010 / ADR-0020). Menu-bar accelerators are in scope. Keyboard **list** navigation beyond Enter / Esc / Cmd-W remains out of scope.

## Out of scope here

Undo, toasts, saved window layouts, password reveal, notifications, keyboard list navigation beyond Enter/Esc/Cmd-W (menu-bar accelerators are in scope per ADR-0020).
