# ADR-0005: Go toolchain and lint

Status: Accepted
Date: 2026-11-16
Owners: MACROX RouteDeck

## Context

Wails v3.0.0 requires Go 1.26 or newer. The scaffold `go.mod` declares `go 1.26.0` and `github.com/wailsapp/wails/v3 v3.0.0`.

## Decision

- Go 1.26+ for all builds. See https://go.dev/doc/go1.26.
- Wails CLI and library stay on v3.0.0 until a task upgrades them.
- Go quality gate is golangci-lint with `.golangci.yml`. The deprecated standalone `golint` binary is not the gate.
- Frontend ESLint is not required. The scaffold has none. Do not add it. Typecheck with `tsc` in the Vite scripts.
- Tests are standard-library table-driven tests. See `docs/engineering/testing.md`.
- Agents verify library APIs against installed docs, not memory.

## Alternatives

- ESLint as a merge gate. Rejected for v1.
- `golint` only. Rejected. It is deprecated and narrower than golangci-lint.

## Consequences

- `golangci-lint run ./...` is part of Go DoD.
- CI (RD-007) should call the same lint config, `go test ./...`, and `npm run build` in `frontend/`. No ESLint step.

## Follow-up

golangci-lint v2.15.0 is pinned in `docs/operations/development.md`. Config is `.golangci.yml` (`version: "2"`).

## Open decisions (owner)

None.
