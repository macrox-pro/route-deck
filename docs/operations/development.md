# Development

## Toolchain

| Tool | Required |
| --- | --- |
| Go | 1.26+ (`go.mod` says `go 1.26.0`) |
| Wails | v3.0.0 (`wails3`) |
| Node | npm for `frontend/` |
| golangci-lint | v2.15.0 (see install) |
| C toolchain | Needed later for `gorm.io/driver/sqlite` (CGO). Not required to compile the current scaffold |

ESLint is not installed and is not required. Do not add it.

## Vendor

This repo commits `vendor/` (ADR-0006). Owner added it with `go mod vendor`.

```bash
go mod vendor
```

After any change to `go.mod` or `go.sum`, re-run `go mod vendor` and keep `vendor/` in the owner commit. MUST NOT gitignore `vendor/`. MUST NOT delete `vendor/` as cleanup.

With `vendor/` present, Go defaults to `-mod=vendor`. Prefer reading APIs under `vendor/` before inventing calls.

## Commands that exist in this repo

From the repo root:

```bash
wails3 dev
wails3 task build
wails3 doctor
go test ./...
go mod vendor
gofmt -l .
```

`main.go` embeds `frontend/dist` (gitignored). Run `cd frontend && npm run build` once before `go test ./...` or `golangci-lint run ./...` if that directory is missing. CI builds the frontend first for the same reason.

Frontend (`frontend/`):

```bash
npm run dev
npm run build
npx untitledui@latest add <component> --yes
```

`npm run build` runs `tsc` and Vite. That is the frontend check. There is no `npm test` script.

Untitled UI (ADR-0007): add components from `frontend/` with `npx untitledui@latest add <name>`. MUST NOT run `init` (it scaffolds a new app). Agent rules: `frontend/AGENTS.md`.

`--yes` implies `--overwrite` and will replace customized files such as `src/components/base/buttons/button.tsx`. Prefer a manual copy from the Untitled UI React repo, or back up customized files before `add --yes` and restore them after.

## golangci-lint

Config: `/.golangci.yml` (`modules-download-mode: vendor`). Version pin: **v2.15.0**. `linters.exclusions.paths` and `formatters.exclusions.paths` include `^build/`: the Wails scaffold under `build/` is not product Go and is not linted. `vendor/` stays excluded by the tool default.

`golangci-lint run ./...` also reports unformatted files (`gofmt` formatter) and enum switches that miss a constant (`exhaustive` linter). Fix formatting with `golangci-lint fmt ./...`. Lint shows one issue per line, so a formatting issue can hide behind another issue on the same line until that one is fixed.

Install (official binary method). From this repo, override vendor mode — the CLI is not in `vendor/`:

```bash
GOFLAGS=-mod=mod go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
```

The binary lands in `$(go env GOPATH)/bin`. Confirm:

```bash
golangci-lint version
golangci-lint run ./...
```

Go DoD requires this command and a recorded result. FAIL blocks DONE.

This pin uses the `go install` version suffix, not a repo tool directive. Change the version only in this file and in CI together.

## Module name

`go.mod` is `module github.com/macrox-pro/route-deck`. Do not rename it inside unrelated tasks.

## Docs for agents

`AGENTS.md`, then `docs/project/STATUS.md`, then `docs/project/HANDOFF.md`.
