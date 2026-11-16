# Go

Baseline: [Effective Go](https://go.dev/doc/effective_go) and [Code Review Comments](https://go.dev/wiki/CodeReviewComments). Toolchain: Go 1.26+ ([release notes](https://go.dev/doc/go1.26)), Wails v3.0.0.

## MUST

- Name packages, types, and functions for what they are. If a comment is required to explain the name, rename it.
- Package names: short, lowercase, no underscores, one concept.
- Wrap errors with `%w`. Use sentinel or typed errors where callers branch.
- Pass `context.Context` into I/O, repository, and process calls. Set a timeout at the edge.
- Give every goroutine an owner and a stop path.
- Check APIs under `vendor/` (then installed module docs) before calling a library.
- Keep `vendor/` in sync: after `go.mod` / `go.sum` changes, run `go mod vendor` (ADR-0006).
- Run `golangci-lint run ./...` with the repo `.golangci.yml` before Go DoD. Put that command in the plan Verification list and in its own todo. A recorded FAIL blocks DONE. Do not start the next task while the command is FAIL, unless the READY task quotes an owner waiver for that command. `build/` is excluded in that config: it is the Wails scaffold, not product Go.
- Prefer `new(expr)` (Go 1.26) when it makes an optional pointer obvious.
- Keep every `.go` file `gofmt`-clean. `golangci-lint run ./...` reports formatting through the `gofmt` formatter. `golangci-lint fmt ./...` rewrites files in place.
- Compute a conversion, call, or expression once per function. Bind it to a local when it is used twice.
- Parse a closed string enum by converting once into a local, switching on that local, and returning it. The `case` lists every constant of the type. The `exhaustive` linter enforces this; `default` does not count as coverage.
- Add a new enum constant, its parser `case`, and its test case in the same change.
- Re-read the diff against this file before Go DoD. Lint PASS is necessary, not sufficient.

## MUST NOT

- Compare `err.Error()` strings.
- Use `panic` for expected failures.
- Synchronize with `time.Sleep`.
- Store business config in mutable package-level variables.
- Add an interface for every struct. Interfaces belong at ports.
- Add a dependency for a trivial helper without refreshing `vendor/`.
- Gitignore or delete `vendor/` as cleanup.
- Create `Util`, `Helper`, or `Manager` packages without a domain noun.
- Extract a function used once unless the block is too large to test in place. Then the name MUST say what it does.
- Add `helpers.go` grab-bags.
- Comment to restate the code. Comments are English and rare.
- Claim lint passed without pasting the command result.
- Silence `exhaustive` with `//exhaustive:ignore`, `nolint`, or `default-signifies-exhaustive`. A switch that handles a subset lists the other constants in an explicit reject `case`.
- Disable a linter or formatter in `.golangci.yml` to make a change pass.

## Examples

Correct: `return fmt.Errorf("start proxybridge: %w", err)`

Forbidden: `if err.Error() == "proxy down"`.

Correct: keep a 15-line validation inline until a second caller exists.

Forbidden: `func helper(x Profile) error` used only from one function.

Correct: closed enum parser.

```go
func ParseProxyKind(s string) (ProxyKind, error) {
	kind := ProxyKind(s)
	switch kind {
	case ProxyKindSOCKS5, ProxyKindHTTP:
		return kind, nil
	default:
		return "", fmt.Errorf("parse proxy kind %q: %w", s, ErrUnknownProxyKind)
	}
}
```

Forbidden: `switch ProxyKind(s) { case ...: return ProxyKind(s), nil }` (repeated conversion), or a `case` list that omits a constant of the type.
