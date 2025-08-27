package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v3"
)

// Callback is the operation to retry.
type Callback[T any] func() (T, error)

// OnRetryHook is called after each failed attempt.
// If it returns true, the retry loop will stop early.
type OnRetryHook func(retryCount int, nextDelay time.Duration, lastErr error) bool

// DoWithBackoff executes callback with exponential backoff until success,
// context cancellation, or retry limits from opts.
// Returns the first successful result or the last error.
func DoWithBackoff[T any](ctx context.Context, callback Callback[T], opts ...OptionFunc) (T, error) {
	options := NewOptions(opts...)

	// Apply total timeout if set
	if options.MaxElapsedTime > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.MaxElapsedTime)
		defer cancel()
	}

	// Configure backoff policy
	expPolicy := backoff.NewExponentialBackOff()
	expPolicy.InitialInterval = options.InitialInterval
	expPolicy.MaxInterval = options.MaxInterval
	expPolicy.Multiplier = options.Multiplier
	expPolicy.MaxElapsedTime = options.MaxElapsedTime
	expPolicy.RandomizationFactor = options.RandomizationFactor
	policy := backoff.WithContext(expPolicy, ctx)

	var (
		result T
		err    error
	)

	for attempt := 0; attempt <= options.MaxRetries; attempt++ {
		result, err = callback()
		if err == nil {
			return result, nil
		}

		// if we've reached the max retries, break
		if attempt == options.MaxRetries {
			break
		}

		delay := policy.NextBackOff()

		// early exit hook
		if options.OnRetryHook != nil && options.OnRetryHook(attempt+1, delay, err) {
			break
		}

		// log the retry
		if options.Logger != nil {
			options.Logger.
				Warn().
				Ctx(ctx).
				Err(err).
				Int("retry_count", attempt+1).
				Int("max_retries", options.MaxRetries).
				Dur("delay", delay).
				Dur("elapsed_time", expPolicy.GetElapsedTime()).
				Msg("retrying operation")
		}

		// Wait or exit on context done
		select {
		case <-ctx.Done():
			return result, fmt.Errorf("context canceled: %w", ctx.Err())
		case <-time.After(delay):
		}
	}

	return result, err
}
