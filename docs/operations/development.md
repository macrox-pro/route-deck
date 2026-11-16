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

## Commands that exist in this repo

From the repo root:

```bash
wails3 dev
wails3 task build
wails3 doctor
go test ./...
gofmt -l .
```

Frontend (`frontend/`):

```bash
npm run dev
npm run build
```

`npm run build` runs `tsc` and Vite. That is the frontend check. There is no `npm test` script.

## golangci-lint

Config: `/.golangci.yml`. Version pin: **v2.7.2**.

Install (official binary method):

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
```

The binary lands in `$(go env GOPATH)/bin`. Confirm:

```bash
golangci-lint version
golangci-lint run ./...
```

Go DoD requires this command and a recorded result. FAIL blocks DONE.

This pin uses the `go install` version suffix, not a repo tool directive. Change the version only in this file and in CI together.

## Module name

`go.mod` is still `module changeme` until RD-002. Do not rename it inside unrelated tasks.

## Docs for agents

`AGENTS.md`, then `docs/project/STATUS.md`, then `docs/project/HANDOFF.md`.
