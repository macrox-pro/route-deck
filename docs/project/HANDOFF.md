# Handoff

Timestamp: 2026-11-17 22:30 UTC
Agent: agent (RD-049..RD-056 editor dialogs epic)
Task ID: RD-056
Branch / commit: working tree (agent MUST NOT commit)

## Done

- **RD-049:** ADR-0014, desktop-ux.md / `.mdc`, frontend AGENTS, contracts (wails-api, lifecycle, domain-model), ADR-0008 supersede note
- **RD-050:** Nested `.app` discovery (Resources), EvalSymlinks dedupe, Versions/Current paths, IsUpdater; ChatGPT-shaped fixture tests
- **RD-051:** HelperGroupDTO from Go; sync ConfirmHelpers (add/remove/ErrReferenced); BuildRulePlan collapse + ErrConflictingMatch; pendingCount; catalog.changed
- **RD-052:** EditorDialogs host (AttachModal, SetEnabled, dirty Question); dialog routes outside AppShell; DialogFooter platform button order; tray focuses open editor
- **RD-053:** Profile editor window (search, Configured only, collapsible HelperGroupDTO); profiles page list-only + delete confirm
- **RD-054:** Proxy editor window (guide prefill, password never pre-filled); proxies page list-only
- **RD-055:** Helpers dialog (group rows, hints, Adds/Removes preview); auto-open after Add; Applications uses Go pendingCount; TS groupers removed
- **RD-056:** Dock badge `!` via `pkg/services/dock` + BadgeFor table tests; `go mod vendor` refreshed

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/...` | PASS |
| `go test ./internal/adapter/inbound/dock/` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `cd frontend && npm run build` | PASS |
| `wails3 generate bindings -ts` | PASS (EditorDialogsService + HelperGroupDTO) |
| `wails3 build` | PASS |
| Owner smoke (dialogs + Dock) | NOT RUN (owner) |

## Exact next step

Owner: `wails3 dev` → open Profile / Proxy / Helpers dialogs; confirm parent disables and re-enables; dirty close asks Discard; after Add with helpers the helpers dialog opens; uncheck referenced helper shows ErrReferenced; Dock shows `!` when applied protection is not Active.

## Suggested commit messages (one per task)

```text
document adr-0014 editor dialogs and desktop ux rules (rd-049)
```

```text
scan nested helper apps with symlink dedupe and current paths (rd-050)
```

```text
sync helper confirm with go groups and rule collapse (rd-051)
```

```text
add editor dialogs host and dialog shell routes (rd-052)
```

```text
move profile editing into a dedicated dialog window (rd-053)
```

```text
move proxy editing into a dedicated dialog window (rd-054)
```

```text
add helpers review dialog driven by helpergroupdto (rd-055)
```

```text
show dock attention badge when applied protection is not active (rd-056)
```

## Context to read

- `docs/decisions/ADR-0014-editor-dialogs-helpers-dock.md`
- `docs/engineering/desktop-ux.md`
- `docs/project/tasks/RD-056.md`
- `internal/adapter/inbound/wails/editor_dialogs.go`
- `internal/adapter/inbound/dock/`

## Notes

- Bindings regenerated; vendor includes `pkg/services/dock`.
- TypeScript `groupHelpersByBundle` / `groupProcessesByBundle` removed; UI renders Go `HelperGroupDTO` only.
- Silent Direct fallback still forbidden; ErrConflictingMatch / ErrReferenced are typed errors.
