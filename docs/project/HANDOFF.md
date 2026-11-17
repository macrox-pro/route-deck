# Handoff

Timestamp: 2026-11-17 20:05 UTC
Agent: agent (RD-058 buildmeta version/commit)
Task ID: RD-058
Branch / commit: working tree (agent MUST NOT commit)

## Done

- `internal/buildmeta`: `Version`/`Commit`, `DisplayVersion`, `PlistVersion` + table tests
- darwin Taskfile: `-ldflags -X` via `APP_VERSION` / `APP_COMMIT`; stamp Info.plist
- `build-macos.yml`: resolve tag → version (strip `v`) or `dev`; short SHA; pass into package
- SettingsDTO + Settings About; sidebar badge from `appVersionLabel`
- Bindings regenerated; contracts/build/release docs updated

## Verification

| Command | Result |
| --- | --- |
| `go test ./internal/buildmeta/` | PASS |
| `go test ./internal/adapter/inbound/wails/` | PASS |
| `golangci-lint run ./...` | PASS (0 issues) |
| `cd frontend && npm run build` | PASS |
| `wails3 generate bindings -ts` | PASS |
| Remote tagged `build-macos` | NOT RUN (owner) |

## Exact next step

Owner: commit, then optionally  
`wails3 task darwin:package APP_VERSION=0.1.0 APP_COMMIT=$(git rev-parse --short=7 HEAD)`  
and check Settings → About + sidebar badge; or push tag `v0.1.0`.

## Suggested commit message

```text
inject tag version and short commit into macos build and settings (rd-058)
```

## Context to read

- `docs/project/tasks/RD-058.md`
- `internal/buildmeta/`
- `.github/workflows/build-macos.yml`
- `frontend/src/pages/settings.tsx`
- `frontend/src/components/application/app-shell/app-shell.tsx`

## Notes

- DEV/`wails3 dev` stays `dev` (production ldflags path only).
- Info.plist uses `0.0.0` when Version is `dev`.
