# Handoff

Timestamp: 2026-11-17 18:55 UTC
Agent: agent (RD-048 profiles app editor)
Task ID: RD-048
Branch / commit: working tree (agent MUST NOT commit)

## Done

- Profiles editor replaced with Applications-style list: all catalog apps, icon, name, Action (None/Proxy/Direct/Block) on the right
- Confirmed helpers nested under each app with independent Action
- **Apply app action to helpers** copies main Action (+ proxy fields) onto helpers; does not auto-cascade on main change
- Flat policy rows and Add-all block removed; Save still writes per-process PolicyDTO (None omitted)
- Review helpers pre-checks helpers already in ProcessIdentity (bundleId)
- Same bundle id grouped to one checkbox / one profile row (ProxyBridge matches bundle/executable, not path); ConfirmHelpers skips duplicate bundle ids

## Verification

| Command | Result |
| --- | --- |
| `cd frontend && npm run build` | PASS |
| `go test ./internal/application/service/ -run TestConfirmHelpersSkipsDuplicateBundleID` | PASS |
| `golangci-lint run ./...` | PASS (0 issues; cache write warnings in sandbox) |
| Owner smoke | NOT RUN (owner) |

## Exact next step

Owner: `wails3 dev` → Profiles → set ChatGPT Block → Apply app action to helpers → Save → Applications Review helpers… shows previously confirmed helpers checked.

## Suggested commit message

```text
replace profiles policy rows with app-grouped editor and preselect helpers (rd-048)
```

## Context to read

- `docs/project/tasks/RD-048.md`
- `frontend/src/pages/profiles.tsx`
- `frontend/src/pages/applications.tsx`

## Notes

- ConfirmHelpers remains additive for new bundle ids; uncheck does not remove catalogued helpers.
- Path is not used in processNames (contract); grouping by bundle id is required for correct rules.
- RD-046/RD-047 may still be uncommitted with RD-048.
