# Auto-run

Cursor setting: Agents → Approvals & Execution.

Use **Auto-review**. Do not use **Run Everything**.

The allowlist is convenience, not a security boundary.

## Suggested allowlist

```text
go test
go vet
gofmt
golangci-lint run
npm run build
wails3 doctor
```

There is no `npm test` script. Do not allowlist `npm test` until one exists.

## MUST NOT auto-run

- `git commit`, `git push`, `git reset`, `git clean`
- `rm -rf`
- Keychain writes (`security add-generic-password` and similar)
- Anything that talks to a production host

Agents still MUST NOT commit, even if a command is allowlisted. `git commit` stays off the allowlist.
