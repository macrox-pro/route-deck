# Boundaries

## Allowed imports

| From | May import | MUST NOT import |
| --- | --- | --- |
| `internal/domain/...` | stdlib | `gorm.io/gorm`, `github.com/wailsapp/wails`, `os/exec` |
| `internal/application/...` | domain, stdlib | GORM, Wails, ProxyBridge JSON types |
| `adapter/inbound/wails` | application services, DTO | `*gorm.DB`, `exec.Command` |
| `adapter/outbound/persistence` | GORM, domain mapping | Wails, ProxyBridge process |
| `adapter/outbound/proxybridge` | process port, profile serializer | GORM |
| React | generated Wails bindings | Go business rules copied into TS |

## Ports

Introduce an interface only at a boundary: database, Keychain, process, clock if tests need it, or a test double.

MUST NOT add an interface per function. MUST NOT add a universal `Manager` type. MUST NOT add a DI framework.

## Transactions

Application services own transaction boundaries. A transaction MUST NOT stay open across a ProxyBridge start, stop, or health check.

## UI

Wails services validate input, call a use case, return DTO and typed errors. Long work returns an `operation_id` and reports progress by events. Success of an RPC is not proof of routing.

## Examples

Correct: application calls `ProxyStore.Get` and `SecretStore.Load`, then a compiler returns a profile, then a process port starts the CLI.

Forbidden: a Wails service opens `gorm.Open` and runs `exec.Command("sh", "-c", ...)`.
