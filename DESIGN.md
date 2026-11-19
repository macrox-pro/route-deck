---
version: alpha
name: RouteDeck Desktop
description: macOS desktop shell for MACROX RouteDeck. Frosted materials, one burgundy action color, Untitled UI components. Target system; theme.css still ships the Untitled purple brand scale until a branding task applies these tokens.
---

# RouteDeck Desktop

Target visual system for the Wails shell. Read this file before generating UI. Behavior that is already decided lives in `docs/engineering/desktop-ux.md`, ADR-0007, and `docs/branding/README.md`. This file does not override those. It does not recolor `frontend/src/styles/theme.css` by itself.

Sources blended here:

- [macOS Vibrancy / Translucency](https://designmd.app/library/macos-vibrancy) for material tiers (sidebar, content, dialog), hairline borders, and blur-plus-saturate. The library file is a landing-page skin (`#007aff`, `#ff9f0a`, a fake wallpaper, 4rem headlines). Those tokens are not used.
- The 2025–2026 burgundy accent trend: one oxblood color for actions and selected chrome, on neutral or frosted surfaces, instead of SaaS purple or system blue.
- The shipped kit: Untitled UI in `frontend/src/components/`, semantic color classes, React Aria.

## Overview

RouteDeck is a macOS instrument for per-app proxy rules. The window should feel like a native utility: a quiet sidebar, readable lists, one obvious action. The app mark carries the blue–cyan–magenta deck. The interface accent is burgundy. Status colors stay separate from both.

Personality: precise, calm, dense enough for a catalog, never a marketing page.

### What the shell is today

The kit is in place and several screens already follow it. Settings is the reference page: a short header, a two-column section grid, grouped rows on `bg-secondary_subtle` with a hairline ring. Activity uses `EmptyState`. Protection and proxy health already map closed states onto success, warning, error, and gray badges.

The rest of the shell still reads as an unstyled catalog:

- Dashboard, Applications, and Help open with implementation notes (Applied versus Active, helper bypass). That copy belongs in Help, not in the page header.
- Profiles, Proxies, Applications, and the Dashboard proxy list are the same full-width divided ring. Empty rows are a single sentence. Only Activity uses the empty-state component.
- Profiles and Proxies confirm delete with `window.confirm`. Editor dialogs already use the native question dialog.
- Setup guides exist at `/help` and are absent from the sidebar.
- Active navigation and primary buttons use the Untitled purple `--color-brand-*` scale (`#7F56D9` family). The mark, sampled from `docs/branding/route-deck-icon.png`, is deep blue (`#031489`), cyan (`#05e9fe`), and magenta (`#982ef7`) on a white arrow (`#f2f5fd`).
- The sidebar is a flat `bg-secondary` panel. There is no material separation between chrome and content.
- Payment icons, the Untitled logo, and `mobile-header.tsx` are kit leftovers. The product shell does not use them.

### Material, and what is not material yet

Vibrancy in this system means three surfaces, not glass on every control.

| Tier | Role | Blur |
| --- | --- | --- |
| Sidebar | Long-lived chrome. Recedes. | Yes, when a backdrop exists behind it |
| Content | Lists, forms, settings groups. Must stay readable | No |
| Dialog | Native-frame editor window. Solid enough to type in | No |

`backdrop-filter` inside an opaque WKWebView blurs other web content. It does not sample the desktop. Real desktop frost is `Mac.Backdrop = MacBackdropTranslucent` on the Wails window. Vendor text in `webview_window_options.go` says a transparent webview needs the `private_mac_apis` build tag; otherwise the webview stays opaque above the backdrop. `MainWindowOptions` in `internal/adapter/inbound/wails/window_chrome.go` leaves backdrop at the zero value (`MacBackdropNormal`) and does not set `BackgroundColour`.

A later window task may turn on `MacBackdropTranslucent` for the main window and the editor windows. That task keeps the native frame (`Frameless` stays false), does not set a white `BackgroundColour`, and does not use `MacBackdropLiquidGlass` (macOS 15+; the product floor is macOS 13). Until that task ships, paint the sidebar with the solid fallback colors below. Do not fake a wallpaper gradient behind the app.

## Colors

Burgundy is the only chromatic UI accent: primary buttons, the selected sidebar item, focus rings, and brand-tinted text. The mark stays blue, cyan, and magenta and is not recolored. Cyan may tint a route glyph. It is not a button fill and not body text (it fails on white). Magenta stays on the mark.

Status colors are not burgundy. Active is success green. Degraded is warning amber. Unprotected and Failed are error red. Unknown, Inactive, and Initializing are gray. White text on burgundy `#6E2432` contrasts at about 10.7:1. Error red `#B42318` against that same burgundy is about 1.7:1, so the two fills must never sit on each other.

```yaml
colors:
  primary: "#6E2432"
  on-primary: "#FFFFFF"
  primary-hover: "#5C2430"
  accent-text: "#6E2432"
  accent-text-dark: "#E7B7BE"
  mark-field: "#031489"
  mark-cyan: "#05e9fe"
  mark-magenta: "#982ef7"
  mark-white: "#f2f5fd"
  error: "#B42318"
  on-error: "#FFFFFF"
  warning: "#B54708"
  success: "#067647"
  light:
    canvas: "#F4F1F2"
    sidebar: "#E7E2E4"
    surface: "#FFFcfb"
    surface-strong: "#FFFFFF"
    text: "#1D1D1F"
    text-muted: "#6E6568"
    hairline: "#1A121418"
    scrim: "#F4F1F2CC"
  dark:
    canvas: "#161214"
    sidebar: "#1C1618"
    surface: "#241C1E"
    surface-strong: "#2C2224"
    text: "#F4F1F2"
    text-muted: "#C4B4B8"
    hairline: "#FFFFFF1A"
    scrim: "#161214CC"
    primary: "#9A4452"
    on-primary: "#FFFFFF"
```

`colors.light.hairline` and `colors.dark.hairline` are 8-digit hex (`#RRGGBBAA`): black at 8% on light, white at 10% on dark.

When a translucent backdrop is actually enabled, the sidebar may use `rgba(255,255,255,0.62)` in light and `rgba(22,18,20,0.72)` in dark, with `backdrop-filter: blur(20px) saturate(1.6)` and the hairline. Content and inputs stay on `surface-strong`. If `prefers-reduced-transparency: reduce`, or the webview is still opaque, use the solid `sidebar` token and drop the filter.

Map these onto Untitled semantic classes. Do not add raw `purple-*`, `blue-*`, or `red-*` utilities in product screens.

| Token | Untitled role |
| --- | --- |
| `colors.primary` | `--color-brand-600` and solid brand buttons, after the branding task |
| `colors.light.canvas` | `bg-primary` page canvas in light mode |
| `colors.light.sidebar` | Sidebar background |
| `colors.light.surface-strong` | List and settings cards |
| `colors.error` / `warning` / `success` | Existing `text-error-primary`, `text-warning-primary`, `text-success-primary` |

## Typography

UI text uses the system face so the shell matches the macOS frame. Inter remains the fallback already declared in `frontend/src/styles/theme.css`. Technical values (host, port, bundle id, operation id, log time) use a monospace face. Page titles stay at the current `text-display-xs` size. Do not introduce a display size above 24px.

```yaml
typography:
  headline:
    fontFamily: "-apple-system, BlinkMacSystemFont, Inter, sans-serif"
    fontSize: 24px
    fontWeight: 600
    lineHeight: 1.25
    letterSpacing: -0.011em
  title:
    fontFamily: "-apple-system, BlinkMacSystemFont, Inter, sans-serif"
    fontSize: 16px
    fontWeight: 600
    lineHeight: 1.4
  body:
    fontFamily: "-apple-system, BlinkMacSystemFont, Inter, sans-serif"
    fontSize: 14px
    fontWeight: 400
    lineHeight: 1.45
  label:
    fontFamily: "-apple-system, BlinkMacSystemFont, Inter, sans-serif"
    fontSize: 13px
    fontWeight: 500
    lineHeight: 1.3
  caption:
    fontFamily: "-apple-system, BlinkMacSystemFont, Inter, sans-serif"
    fontSize: 12px
    fontWeight: 400
    lineHeight: 1.35
  mono:
    fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace"
    fontSize: 12px
    fontWeight: 400
    lineHeight: 1.4
```

## Layout

Desktop only. The shell is a fixed sidebar plus a scrolling main column. Widths already in `sidebar-simple.tsx`: expanded 232px, collapsed 68px. Main padding is 24px, 32px from the `md` breakpoint. Settings stays `max-w-3xl` and centered. Lists use the full main width.

```yaml
spacing:
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
  sidebar-expanded: 232px
  sidebar-collapsed: 68px
  page: 24px
  section: 24px
```

Page header, every route: one title, one sentence about the user outcome, one primary button. Secondary actions are `color="secondary"`. Setup guides is a sidebar item, same weight as Settings, not a spare button on the Dashboard.

Dashboard content order:

1. Protection status and bridge state on one row.
2. Setup checklist, when it is incomplete.
3. Warnings that need a decision (bridge missing, desired differs from applied).
4. Proxy health list.

Do not stack a glossary under the title.

Lists are rows: leading icon (32px), name, one muted meta line, trailing status, then row actions. Destructive actions sit last. More than three row actions collapse behind a menu when that component is already in the kit; do not invent a new menu.

Dialogs keep the sizes and footer rules in `docs/engineering/desktop-ux.md`. Body scrolls. Title and footer stay. macOS footer order is Cancel then Save.

## Elevation & Depth

Depth comes from material tiers and hairlines, not drop shadows.

- Sidebar sits on `sidebar` (or the translucent material). A 1px hairline separates it from the canvas. No shadow.
- Lists and settings groups sit on `surface-strong` with a 1px hairline and an inner top highlight `inset 0 1px 0 rgba(255,255,255,0.6)` in light mode. In dark mode the highlight is `rgba(255,255,255,0.08)`.
- Popovers and select menus may use the existing Untitled shadow. Product cards do not.
- Do not put `backdrop-filter` on buttons, badges, inputs, or list rows.

## Shapes

Controls stay close to the current Untitled radii. Materials are slightly rounder than inputs. Badges are pills. Do not use 20px-plus marketing corners.

```yaml
rounded:
  sm: 6px
  md: 8px
  lg: 12px
  full: 9999px
```

| Element | Radius |
| --- | --- |
| Inputs, row buttons, sidebar items | `md` (8px) |
| Lists, settings groups, guide cards | `lg` (12px) |
| Protection and health badges | `full` |
| App mark in the sidebar | 8px, as `RouteDeckMark` already does |

## Motion & Micro-interaction

Motion confirms a change. It does not introduce the screen. The sidebar width transition stays at 200ms. Hover and color changes stay at 100ms, which is the Untitled default. Honor `prefers-reduced-motion`. Do not stagger list items. Do not animate layout properties.

```yaml
motion:
  state-change:
    duration: 100ms
    delay: 0ms
    timingFunction: linear
  sidebar:
    duration: 200ms
    delay: 0ms
    timingFunction: linear
  surface-enter:
    duration: 150ms
    delay: 0ms
    timingFunction: ease-out
```

## Voice & Content Strategy

- Tone: precise and short. Say what the control does and what state the snapshot reported. Skip implementation lectures in page subtitles.
- Capitalization: sentence case for buttons, form labels, and empty states ("New profile", "Recheck health"). Navigation labels and page titles stay single title-case nouns (Dashboard, Proxies, Activity).
- Terminology: Protection states are the closed words Active, Degraded, Inactive, Unprotected, Failed, Unknown, Initializing. Applied is a profile fact, not a protection state. Proxy health display words are Healthy, Unreachable, Auth failed, Checking, and the TCP / UDP pair. Do not rename them in the UI.
- Buttons name the outcome. Save, Apply, Stop, Delete. Not OK, Submit, or Continue.
- Empty states name the missing object and the next action ("No proxies yet" plus "New proxy"). Loading states name the object ("Loading snapshot…"). Errors say what failed and the next step, in the existing alert pattern (`role="alert"`, error surface, error text).
- Numbers and addresses: times in Activity stay 24-hour local. Hosts, ports, and bundle ids use `typography.mono`.
- Language: English.

## Components

Use the files below. Add missing Untitled pieces with `npx untitledui@latest add <name> --yes` from `frontend/`. Do not run `init`. Do not restyle payment icons, marketing illustrations, or the Untitled logo into the product.

- **Buttons** → `@frontend/src/components/base/buttons/button.tsx` — one `primary` per view header; `secondary` for the rest; `secondary-destructive` for Delete. `md` in page headers and dialog footers. `sm` on list rows. Loading uses `isLoading` and `showTextWhileLoading`.
- **Inputs and selects** → `@frontend/src/components/base/input/input.tsx`, `@frontend/src/components/base/select/select.tsx` — label above the field, hint under it, error text under the hint. Focus ring uses burgundy, not system blue.
- **Toggles and checkboxes** → `@frontend/src/components/base/toggle/toggle.tsx`, `@frontend/src/components/base/checkbox/checkbox.tsx` — selected fill is burgundy. Settings rows keep the label on the left and the control on the right, as Settings already does.
- **Badges** → `@frontend/src/components/base/badges/badges.tsx` — protection via `@frontend/src/components/application/protection-badge/protection-badge.tsx`, health via `@frontend/src/components/application/proxy-health-badge/proxy-health-badge.tsx`. Burgundy is not a status color.
- **Lists** — product rows, not the data table, until a screen needs sorting. Container: hairline, `rounded.lg`, divided rows, 12px vertical padding.
- **Empty states** → `@frontend/src/components/application/empty-state/empty-state.tsx` — same composition as Activity. Featured icon theme `modern`, color `gray`.
- **Sidebar** → `@frontend/src/components/application/app-navigation/sidebar-navigation/sidebar-simple.tsx` — selected item: burgundy icon and label on `surface-strong`. Unselected icons stay quaternary. Collapsed rail keeps tooltips.
- **Dialogs** → `@frontend/src/components/application/dialog-shell/dialog-shell.tsx` — native frame, rules in `docs/engineering/desktop-ux.md`.
- **Alerts** — the existing `rounded-lg` ring pattern on Dashboard and Settings. One alert color per meaning. Do not stack the same warning style for unrelated facts; merge or move the explanation to Help.

```yaml
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.on-primary}"
    rounded: "{rounded.md}"
    padding: 10px
  button-primary-dark:
    backgroundColor: "{colors.dark.primary}"
    textColor: "{colors.dark.on-primary}"
    rounded: "{rounded.md}"
    padding: 10px
  button-secondary:
    backgroundColor: "{colors.light.surface-strong}"
    textColor: "{colors.light.text}"
    rounded: "{rounded.md}"
    padding: 10px
  sidebar:
    backgroundColor: "{colors.light.sidebar}"
    textColor: "{colors.light.text}"
  sidebar-active:
    backgroundColor: "{colors.light.surface-strong}"
    textColor: "{colors.primary}"
  card:
    backgroundColor: "{colors.light.surface-strong}"
    textColor: "{colors.light.text}"
    rounded: "{rounded.lg}"
    padding: 16px
  input:
    backgroundColor: "{colors.light.surface-strong}"
    textColor: "{colors.light.text}"
    rounded: "{rounded.md}"
    padding: 10px
```

## Do's and Don'ts

- Do keep one burgundy primary action per view.
- Do keep Active, Degraded, Unprotected, Failed, and Unknown on success, warning, error, and gray. Unknown must not look Active.
- Do use solid sidebar and surface colors until a window task enables `MacBackdropTranslucent`.
- Do honor `prefers-reduced-transparency` and `prefers-reduced-motion`.
- Do keep editor windows on the native frame, with Cancel then Save on macOS.
- Do put host, port, and bundle id in the mono style.
- Don't copy the designmd vibrancy landing tokens (`#007aff`, `#ff9f0a`, wallpaper gradient, 4rem type, glass on every button).
- Don't recolor `docs/branding/route-deck-icon.png` or paint the mark burgundy.
- Don't use burgundy as an error, and don't place an error fill on a burgundy fill.
- Don't set `Frameless`, a white `BackgroundColour`, or `MacBackdropLiquidGlass` for this look.
- Don't add a second font, a marketing hero, or a bento feature grid.
- Don't leave new empty lists as a bare sentence when `EmptyState` already exists.
- Don't explain Applied versus Active in a page subtitle. Link to Setup guides.
