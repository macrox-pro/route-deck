# ADR-0007: Untitled UI for React

Status: Accepted
Date: 2026-11-17
Owners: MACROX RouteDeck

## Context

The Wails frontend needs a component system. The owner chose Untitled UI React (source copied into the project, not an npm component library).

## Decision

- UI components come from Untitled UI via `npx untitledui@latest add <name>` run from `frontend/`.
- MUST NOT run `npx untitledui@latest init` in this repository (it scaffolds a new Vite/Next app).
- PRO components, `example` pages, and `@untitledui-pro/icons` stay out of v1 until the owner authenticates.
- Agent frontend rules live in `frontend/AGENTS.md` (RouteDeck preface + Untitled embed).
- React and TypeScript versions follow the Untitled stack currently in `frontend/package.json`.

## Alternatives

- Install a traditional component library as a dependency. Rejected: owner wants source ownership.
- Keep the Wails greeter template. Rejected: product shell required.

## Consequences

- Components under `frontend/src/components/` are project source and may be customized.
- `npm run build` (tsc + Vite) remains the frontend gate. ESLint is not required.
- Brand `--color-brand-*` in `theme.css` may still be the Untitled default purple until a branding task recolors it.

## Follow-up

None for v1 shell. Forms and PRO examples need later READY tasks.

## Open decisions (owner)

None.
