# Testing

## Unit tests

Default shape is table-driven with `t.Run`.

Use it for domain invariants, profile validation, rule compilation, CLI parsers, mappers, and state transitions.

MUST NOT start the real ProxyBridge CLI in unit tests. Use a fake process and fixtures.

MUST NOT add `github.com/google/go-cmp` unless a READY task says so. Compare fields or use `errors.Is`.

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

A test named `TestItWorks` with no cases, or a unit test that execs `proxybridge`.
