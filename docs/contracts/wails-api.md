# Wails API contract

Library: `github.com/wailsapp/wails/v3` v3.0.0. Docs: https://v3.wails.io/

MUST NOT invent methods. Read the installed module or the v3 docs before adding a binding.

## Current scaffold (not product API)

`main.go` registers `GreetService` with `Greet(name string) string` and emits event `time` (`string`) every second. That is template code. Product tasks replace it. Do not build features on `Greet`.

`ApplicationShouldTerminateAfterLastWindowClosed` is `true` in the scaffold. The tray phase sets it to `false`. Do not change it until that task is READY.

## Target shape

Wails services are inbound adapters:

1. Accept a DTO
2. Validate input
3. Call an application service
4. Return a DTO and a typed error

They MUST NOT open GORM and MUST NOT call `exec.Command`.

Long operations return an `operation_id`. Progress and failure use events. A successful call is not proof that routing works. The client loads a snapshot after the call.

## Target services (names fixed, methods not invented)

| Service | Responsibility |
| --- | --- |
| Snapshot | Return bridge, proxy health, protection, desired vs applied |
| Proxies | CRUD for the global proxy catalog (password via Keychain port) |
| Applications | List and select apps, helper confirmation |
| Profiles | Create, copy, switch desired profile |
| Protection | Apply and stop, returning the confirmed snapshot |

Method signatures are specified in the READY task that adds them, after checking Wails v3 service binding rules. This file does not list fake signatures.

## Types

Generated bindings are the TypeScript source for these DTOs. Hand-written duplicates are forbidden. Domain types and GORM models are not the DTO.

## Errors

Return typed errors the UI can map. MUST NOT return a bare string that means “ok” when the CLI did not confirm.
