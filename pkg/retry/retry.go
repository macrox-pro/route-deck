// Package retry runs an operation a fixed number of times with a delay between attempts.
package retry

import (
	"context"
	"fmt"
	"time"
)

// Policy controls how Do retries.
type Policy struct {
	// Attempts is the total number of tries, including the first. Values below 1 are treated as 1.
	Attempts int
	// Delay is waited after a failed attempt before the next try. Ignored after the last attempt.
	Delay time.Duration
}

// Do runs fn until it returns nil or Attempts is exhausted.
// Context cancel between attempts stops early; the last error is wrapped with ctx.Err().
func Do(ctx context.Context, policy Policy, fn func(context.Context) error) error {
	if fn == nil {
		return fmt.Errorf("retry: nil func")
	}
	attempts := policy.Attempts
	if attempts < 1 {
		attempts = 1
	}

	var last error
	for i := 1; i <= attempts; i++ {
		if err := ctx.Err(); err != nil {
			if last != nil {
				return fmt.Errorf("retry: %w: %w", err, last)
			}
			return fmt.Errorf("retry: %w", err)
		}
		last = fn(ctx)
		if last == nil {
			return nil
		}
		if i == attempts {
			break
		}
		if policy.Delay <= 0 {
			continue
		}
		timer := time.NewTimer(policy.Delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("retry: %w: %w", ctx.Err(), last)
		case <-timer.C:
		}
	}
	return last
}
