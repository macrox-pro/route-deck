# Secrets

## Store

Proxy passwords go to the macOS Keychain through a port implemented in the Keychain adapter. SQLite stores the item reference (account + service), not the password.

Domain structs that are saved MUST NOT contain the password. The profile serializer loads it at apply time and writes it into the stdin JSON because the ProxyBridge format requires `password` in clear text.

## Logs and errors

MUST NOT log:

- passwords or tokens
- the full profile JSON
- Keychain item data
- packet payloads

Structured logging destinations, levels, redaction, and stable attribute keys: `docs/engineering/logging.md` (ADR-0008).

## Export

Default export strips `password`. Putting secrets in an export requires an explicit user action in a future task. That task is not approved yet.

## Tests

Fixtures use fake passwords such as `test-password` only in memory. MUST NOT commit a real proxy password. MUST NOT write secrets into the repo SQLite files.

## Process

Agents MUST NOT run Keychain write commands while applying docs or refactors. Auto-run MUST NOT allowlist `security add-generic-password`.
