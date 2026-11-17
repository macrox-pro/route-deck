package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/pkg/retry"
)

func TestDoSuccessFirstTry(t *testing.T) {
	t.Parallel()
	calls := 0
	err := retry.Do(context.Background(), retry.Policy{Attempts: 3}, func(context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestDoRetriesThenSucceeds(t *testing.T) {
	t.Parallel()
	calls := 0
	err := retry.Do(context.Background(), retry.Policy{Attempts: 3, Delay: time.Millisecond}, func(context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
}

func TestDoExhaustsAttempts(t *testing.T) {
	t.Parallel()
	want := errors.New("always")
	calls := 0
	err := retry.Do(context.Background(), retry.Policy{Attempts: 3}, func(context.Context) error {
		calls++
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("err = %v, want %v", err, want)
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
}

func TestDoZeroAttemptsMeansOne(t *testing.T) {
	t.Parallel()
	calls := 0
	_ = retry.Do(context.Background(), retry.Policy{Attempts: 0}, func(context.Context) error {
		calls++
		return errors.New("x")
	})
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}

func TestDoRespectsCancelBetweenAttempts(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := retry.Do(ctx, retry.Policy{Attempts: 5, Delay: time.Hour}, func(context.Context) error {
		calls++
		cancel()
		return errors.New("fail")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}
}
