# ADR-0006: Commit Go module vendor tree

Status: Accepted
Date: 2026-11-16
Owners: MACROX RouteDeck

## Context

The owner ran `go mod vendor` and keeps a `vendor/` tree in the repo so builds and agents can resolve Wails and other modules without relying on a live module download.

## Decision

- Commit `vendor/` with the module. Do not gitignore it.
- After any `go.mod` / `go.sum` dependency change, run `go mod vendor` and include the refreshed tree in the same change set the owner commits.
- Prefer verifying library APIs against `vendor/` (then installed module docs). MUST NOT invent APIs from memory.
- Go 1.14+ with a present `vendor/` defaults to `-mod=vendor`. Agents MUST NOT delete `vendor/` to “clean up” the tree.

## Alternatives

- No vendor; always download modules. Rejected by owner for this repo.
- Vendor only in CI cache. Rejected: local and agent loops need the same tree.

## Consequences

- Clones are larger. Dependency bumps touch many files under `vendor/`.
- `golangci-lint` and `go test` / `go build` should use the vendored modules.
- RD-002 module rename will refresh `vendor/` paths under the new module consumers; still re-run `go mod vendor` after the rename.

## Follow-up

Documented in `docs/operations/development.md`, `docs/engineering/golang.md`, and `AGENTS.md` stack rules.

## Open decisions (owner)

None.
