// Package retry runs an operation a fixed number of times with a delay between attempts.
package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrNilFunc is returned when Do is called with a nil operation.
var ErrNilFunc = errors.New("retry: nil func")

// Func is an operation that may be retried. It receives the same context passed to Do.
type Func func(ctx context.Context) error

// Policy controls how Do retries.
type Policy struct {
	// Attempts is the total number of tries, including the first. Values below 1 are treated as 1.
	Attempts int
	// Delay is waited after a failed attempt before the next try. Ignored after the last attempt
	// and when non-positive.
	Delay time.Duration
}

// Do runs fn until it returns nil or Attempts is exhausted.
//
// Context cancel or deadline before an attempt, after a failed attempt, or during Delay stops
// early. When a prior attempt failed, the returned error unwraps to both ctx.Err() and that
// last error (errors.Is / errors.As work for either). Exhausted attempts return the last
// operation error unchanged.
func Do(ctx context.Context, policy Policy, fn Func) error {
	if fn == nil {
		return ErrNilFunc
	}

	attempts := max(policy.Attempts, 1)

	var last error
	for attempt := 1; attempt <= attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return aborted(err, last)
		}

		last = fn(ctx)
		if last == nil {
			return nil
		}

		if attempt == attempts {
			break
		}

		if err := ctx.Err(); err != nil {
			return aborted(err, last)
		}

		if err := wait(ctx, policy.Delay); err != nil {
			return aborted(err, last)
		}
	}

	return last
}

func aborted(ctxErr, last error) error {
	if last != nil {
		return fmt.Errorf("retry: %w: %w", ctxErr, last)
	}

	return fmt.Errorf("retry: %w", ctxErr)
}

// wait blocks until delay elapses or ctx is done. Non-positive delay is a no-op.
// Own function so defer timer.Stop is not tied to Do's frame across loop iterations.
func wait(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
