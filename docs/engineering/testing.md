# Testing

## Unit tests

Default shape is table-driven with `t.Run`.

Unit tests cover **exported behavior**: public functions and methods. Nested logic, private helpers, and field layout are exercised only through those entry points.

A plain struct with no methods and no constructor is not a unit-test target. Field lists and type shapes are reviewed in code; do not write reflect inventories or “closed field” tests for them.

Use unit tests for domain invariants, profile validation, rule compilation, CLI parsers, mappers, and state transitions — always via the public API that owns that behavior.

MUST NOT start the real ProxyBridge CLI in unit tests. Use a fake process and fixtures.

MUST NOT add `github.com/google/go-cmp` unless a READY task says so. Compare fields or use `errors.Is`.

MUST NOT write tests only to raise coverage, assert struct field order/types, or mirror the type definition without behavior.

A bugfix SHOULD add a failing table case first when the defect is in pure logic.

## Integration tests

GORM + SQLite and the process manager with a fake process are integration tests. Name them so `go test` packages make the split obvious. They may use a temp SQLite file. They still MUST NOT require the system ProxyBridge extension.

## DoD

Record the exact command and PASS, FAIL, or NOT RUN in `docs/project/HANDOFF.md`. Unrun tests are not done.

## Skeleton

```go
func TestValidateProxyRefs(t *testing.T) {
	tests := []struct {
		name    string
		in      Profile
		wantErr error
	}{
		{name: "all_refs_resolve", in: validProfile, wantErr: nil},
		{name: "missing_proxy_id", in: brokenRefProfile, wantErr: ErrUnresolvedProxy},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProxyRefs(tt.in)
			if tt.wantErr == nil && err != nil {
				t.Fatalf("ValidateProxyRefs() error = %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("ValidateProxyRefs() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
```

## Forbidden

A test named `TestItWorks` with no cases, a unit test that execs `proxybridge`, a test that only lists struct fields with `reflect`, or any other test written only for coverage or ceremony.
