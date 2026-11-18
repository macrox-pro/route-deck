package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/pkg/retry"
)

func TestDo(t *testing.T) {
	t.Parallel()

	transient := errors.New("transient")
	always := errors.New("always")

	tests := []struct {
		name        string
		policy      retry.Policy
		setup       func(t *testing.T) (ctx context.Context, fn retry.Func, calls *atomic.Int32)
		wantErr     error
		wantAlsoErr error
		wantCalls   int32
	}{
		{
			name:   "success_first_try",
			policy: retry.Policy{Attempts: 3},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				calls := new(atomic.Int32)
				return context.Background(), func(context.Context) error {
					calls.Add(1)
					return nil
				}, calls
			},
			wantCalls: 1,
		},
		{
			name:   "retries_then_succeeds",
			policy: retry.Policy{Attempts: 3, Delay: time.Millisecond},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				calls := new(atomic.Int32)
				return context.Background(), func(context.Context) error {
					if calls.Add(1) < 3 {
						return transient
					}
					return nil
				}, calls
			},
			wantCalls: 3,
		},
		{
			name:   "exhausts_attempts",
			policy: retry.Policy{Attempts: 3},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				calls := new(atomic.Int32)
				return context.Background(), func(context.Context) error {
					calls.Add(1)
					return always
				}, calls
			},
			wantErr:   always,
			wantCalls: 3,
		},
		{
			name:   "zero_attempts_means_one",
			policy: retry.Policy{Attempts: 0},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				calls := new(atomic.Int32)
				return context.Background(), func(context.Context) error {
					calls.Add(1)
					return always
				}, calls
			},
			wantErr:   always,
			wantCalls: 1,
		},
		{
			name:   "negative_attempts_means_one",
			policy: retry.Policy{Attempts: -2},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				calls := new(atomic.Int32)
				return context.Background(), func(context.Context) error {
					calls.Add(1)
					return always
				}, calls
			},
			wantErr:   always,
			wantCalls: 1,
		},
		{
			name:   "nil_func",
			policy: retry.Policy{Attempts: 3},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				return context.Background(), nil, new(atomic.Int32)
			},
			wantErr:   retry.ErrNilFunc,
			wantCalls: 0,
		},
		{
			name:   "already_canceled_before_first_attempt",
			policy: retry.Policy{Attempts: 3, Delay: time.Millisecond},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				calls := new(atomic.Int32)
				return ctx, func(context.Context) error {
					calls.Add(1)
					t.Fatal("fn must not run when ctx is already canceled")
					return nil
				}, calls
			},
			wantErr:   context.Canceled,
			wantCalls: 0,
		},
		{
			name:   "cancel_during_delay_unwraps_both",
			policy: retry.Policy{Attempts: 5, Delay: 100 * time.Millisecond},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				ctx, cancel := context.WithCancel(context.Background())
				t.Cleanup(cancel)
				calls := new(atomic.Int32)
				return ctx, func(context.Context) error {
					calls.Add(1)
					// Cancel after Do has entered wait, not inside the post-fn check.
					time.AfterFunc(5*time.Millisecond, cancel)
					return transient
				}, calls
			},
			wantErr:     context.Canceled,
			wantAlsoErr: transient,
			wantCalls:   1,
		},
		{
			name:   "zero_delay_still_retries",
			policy: retry.Policy{Attempts: 3, Delay: 0},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				calls := new(atomic.Int32)
				return context.Background(), func(context.Context) error {
					calls.Add(1)
					return always
				}, calls
			},
			wantErr:   always,
			wantCalls: 3,
		},
		{
			name:   "stops_after_fn_when_ctx_done",
			policy: retry.Policy{Attempts: 5, Delay: time.Hour},
			setup: func(t *testing.T) (context.Context, retry.Func, *atomic.Int32) {
				t.Helper()
				ctx, cancel := context.WithCancel(context.Background())
				t.Cleanup(cancel)
				calls := new(atomic.Int32)
				return ctx, func(context.Context) error {
					calls.Add(1)
					cancel()
					return transient
				}, calls
			},
			wantErr:     context.Canceled,
			wantAlsoErr: transient,
			wantCalls:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, fn, calls := tt.setup(t)
			err := retry.Do(ctx, tt.policy, fn)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Do() error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Do() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantAlsoErr != nil && !errors.Is(err, tt.wantAlsoErr) {
				t.Fatalf("Do() error = %v, want unwrap to %v", err, tt.wantAlsoErr)
			}
			if got := calls.Load(); got != tt.wantCalls {
				t.Fatalf("calls = %d, want %d", got, tt.wantCalls)
			}
		})
	}
}
