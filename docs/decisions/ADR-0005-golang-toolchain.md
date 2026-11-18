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
- CI check gate (RD-007) is `.github/workflows/ci.yml`: `npm run build` in `frontend/` first (so `frontend/dist` exists for `//go:embed`), then `go test ./...` and `golangci-lint run ./...`. No ESLint step.
- CI macOS package (RD-057) is `.github/workflows/build-macos.yml`: `wails3 task darwin:package`, ad-hoc codesign only, artifact zip. No Developer ID / notarization in that workflow.

## Follow-up

golangci-lint v2.15.0 is pinned in `docs/operations/development.md`. Config is `.golangci.yml` (`version: "2"`).

Vendoring is decided in ADR-0006 (`vendor/` committed; refresh with `go mod vendor`).

## Open decisions (owner)

None.
