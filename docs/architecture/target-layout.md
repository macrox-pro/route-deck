# Target layout

Do not create empty directories until a READY task needs them. Today the Wails entry is `main.go` and `greetservice.go` at the repo root.

```text
main.go                         # Wails composition root until bootstrap exists
internal/
  domain/                       # entities and invariants only
  application/
    service/
    port/
    workflow/
  adapter/
    inbound/wails/
    outbound/proxybridge/
    outbound/persistence/sqlite/
    outbound/macos/
    outbound/process/
    outbound/keychain/
  platform/darwin/
  bootstrap/                    # only when wiring leaves main.go
frontend/src/
docs/
```

No top-level `migrations/` directory. Schema changes go through GORM `AutoMigrate` on known models. See `persistence.md`.

## Rules

- One concept per package. Package names are short, lowercase, no underscores.
- MUST NOT add `internal/pkg/utils` or `helpers.go` grab-bags.
- Keep `main` as the composition root until a task moves wiring into `internal/bootstrap`.
- Frontend business rules stay out of `frontend/src`.
