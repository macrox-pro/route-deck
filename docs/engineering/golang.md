# Go

Baseline: [Effective Go](https://go.dev/doc/effective_go) and [Code Review Comments](https://go.dev/wiki/CodeReviewComments). Toolchain: Go 1.26+ ([release notes](https://go.dev/doc/go1.26)), Wails v3.0.0.

## MUST

- Name packages, types, and functions for what they are. If a comment is required to explain the name, rename it.
- Package names: short, lowercase, no underscores, one concept.
- Wrap errors with `%w`. Use sentinel or typed errors where callers branch.
- Pass `context.Context` into I/O, repository, and process calls. Set a timeout at the edge.
- Give every goroutine an owner and a stop path.
- Check the installed module docs before calling a library.
- Run `golangci-lint run ./...` with the repo `.golangci.yml` before Go DoD.
- Prefer `new(expr)` (Go 1.26) when it makes an optional pointer obvious.

## MUST NOT

- Compare `err.Error()` strings.
- Use `panic` for expected failures.
- Synchronize with `time.Sleep`.
- Store business config in mutable package-level variables.
- Add an interface for every struct. Interfaces belong at ports.
- Add a dependency for a trivial helper.
- Create `Util`, `Helper`, or `Manager` packages without a domain noun.
- Extract a function used once unless the block is too large to test in place. Then the name MUST say what it does.
- Add `helpers.go` grab-bags.
- Comment to restate the code. Comments are English and rare.
- Claim lint passed without pasting the command result.

## Examples

Correct: `return fmt.Errorf("start proxybridge: %w", err)`

Forbidden: `if err.Error() == "proxy down"`.

Correct: keep a 15-line validation inline until a second caller exists.

Forbidden: `func helper(x Profile) error` used only from one function.
