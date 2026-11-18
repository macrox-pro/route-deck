# Desktop UX

English UI. Untitled UI components. Semantic color tokens. See ADR-0014 and ADR-0007.

## Editor dialogs

Kinds: `profile`, `proxy`, `helpers`. One reused native-frame window each. No Frameless.

| Kind | Size | Resize |
| --- | --- | --- |
| profile | 760×640 (min 600×480) | yes |
| proxy | 480×560 | no |
| helpers | 600×560 | yes |

Titles: "New Profile", "Edit Profile — {name}", "New Proxy", "Edit Proxy — {name}", "Review Helpers — {app}".

- Footer: macOS/Linux `[Cancel][Save]`; Windows `[Save][Cancel]`. Save is default. Enter saves, Esc cancels, Cmd/Ctrl+W closes.
- Body scrolls; title and footer stay put.
- Focus first field on open. Inline Go validation errors. Save shows loading.
- Save disabled when name blank, form clean, or save in flight. Create becomes dirty on first edit.
- Dirty close: `Dialog.Question` "Discard changes?".
- Destructive delete defaults to Cancel; Enter does not delete.
- After save: dialog closes, list refetches, saved row in view, focus returns to opener.
- Theme/dark mode via existing tokens. No white `BackgroundColour` override. Appearance (`system`|`light`|`dark`) is SQLite-backed; native frame uses Wails create-time `Mac.Appearance` (recreate main on change).
- Buttons: `sm` in lists, `md` in dialog footer.

## Profiles

Search filter, "Configured only" toggle, count "N of M" with an action. Helper groups from Go (`HelperGroupDTO`). "Apply app action to helpers" drafts only. Row action opens helpers dialog for that app.

## Helpers

One checkbox per bundle-id group. Executable names under the row. Hints: "Same developer", "Updater". Only already-confirmed groups start checked; new groups start unchecked. Toolbar: "Select all" / "Deselect all" with "N of M selected". Footer preview "Adds N · Removes M". `ErrReferenced` copy names profiles ("Used by Work, Home").

## Proxies

Password never pre-filled. Empty password on edit keeps Keychain secret. Guide prefill via query stays.

## Lists

Empty states name the next action. Applications catalog shows confirmed helper groups read-only.

## Out of scope here

Undo, toasts, saved window layouts, password reveal, notifications, keyboard list navigation beyond Enter/Esc/Cmd-W.
