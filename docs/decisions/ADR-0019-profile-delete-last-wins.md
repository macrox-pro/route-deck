# ADR-0019: Profile delete Desired and last-wins match conflicts

Status: Accepted
Date: 2026-11-19
Owners: MACROX RouteDeck

## Context

Profile Delete refused Desired and Applied with an opaque referenced error (ADR-0008), so users could not remove the profile they were using without a separate Desired change. Profile Save and `BuildRulePlan` returned `ErrConflictingMatch` when policies assigned different fingerprints to identities that share a matchable name (ADR-0014). The editor keys policies by UI representatives while the list counted every stored process id, so badges could show Block while the editor showed None. Owner closed (2026-11-19): clear Desired then delete; Applied only after Stop; last policy in the saved list wins with overwrite — the user never resolves conflicts.

## Decision

### Delete

- Deleting the **Desired** profile (when it is not Applied) clears Desired in the same Delete path, then removes the profile.
- Deleting the **Applied** profile runs **Stop then delete** after a native Question ("Stop and Delete"). Cancel is the default. Confirm uses Wails `Dialog.Question`, not `window.confirm`.
- MUST NOT delete Applied without Stop. MUST NOT silently rewrite policies to Direct.

This **supersedes** ADR-0008 “Delete is refused for the desired or applied profile” for the Desired case only. Applied refusal remains.

### Match conflicts (last-wins)

- Before persisting a profile and before compiling a rule plan, resolve match conflicts by **last-wins**: for each matchable name, the fingerprint of the **last** policy in `profile.Policies` that touches that name wins.
- Persist the resolved policy list so stored data matches compile. Save and Apply MUST NOT surface `ErrConflictingMatch` to the user for this case.
- Never map unresolved or conflicting actions to Direct.

This **supersedes** ADR-0014 “Different actions on a shared matchable name return `ErrConflictingMatch`” for Save and Apply compile. Collapse of duplicate bundle ids at compile time is unchanged.

### Editor WYSIWYG

- On load, remap policies onto UI representatives (main = lowest-id process for the app bundle; helper = first `processIds` entry after scanning the group).
- On Save, persist only editor-owned keys; drop sibling and orphan process policies so list badges match what the user saw.

## Alternatives

- Keep refusing Desired delete. Rejected: owner chose clear-then-delete UX.
- Keep hard-fail conflicts for the user to fix. Rejected: user cannot resolve match conflicts in the UI.
- Prefer main over helper on conflict. Rejected: owner chose last in list.

## Consequences

- RD-068 implements Delete, last-wins resolve, and editor normalize/save.
- ADR-0008 and ADR-0014 gain footnotes pointing here.
- Domain-model and wails-api Delete / conflict lines update.

## Follow-up

None required for v1.

## Open decisions (owner)

None.
