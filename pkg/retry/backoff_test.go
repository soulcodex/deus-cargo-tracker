package retry_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soulcodex/deus-cargo-tracker/pkg/retry"
)

var _ retry.OnRetryHook = nil // ensure type exists

func TestDoWithBackoff_SucceedsImmediately(t *testing.T) {
	ctx := context.Background()
	called := 0

	res, err := retry.DoWithBackoff(ctx, func() (any, error) {
		called++
		return "ok", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
	assert.Equal(t, 1, called)
}

func TestDoWithBackoff_SucceedsAfterRetries(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	res, err := retry.DoWithBackoff(ctx, func() (any, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("fail")
		}
		return 123, nil
	}, retry.WithMaxRetries(5), retry.WithInitialInterval(1*time.Millisecond))
	require.NoError(t, err)
	assert.Equal(t, 123, res)
	assert.Equal(t, 3, attempts)
}

func TestDoWithBackoff_ExhaustsRetries(t *testing.T) {
	ctx := context.Background()
	failErr := errors.New("always fail")
	attempts := 0

	res, err := retry.DoWithBackoff(ctx, func() (any, error) {
		attempts++
		return nil, failErr
	}, retry.WithMaxRetries(2), retry.WithInitialInterval(1*time.Millisecond))
	assert.Nil(t, res)
	assert.Equal(t, failErr, err)
	assert.Equal(t, 3, attempts) // initial + 2 retries
}

func TestDoWithBackoff_OnRetryHookStopsEarly(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	failErr := errors.New("always fail")

	res, err := retry.DoWithBackoff(ctx, func() (any, error) {
		attempts++
		return nil, failErr
	}, retry.WithMaxRetries(10), retry.WithInitialInterval(1*time.Millisecond),
		retry.WithOnRetryHook(func(retryCount int, _ time.Duration, _ error) bool {
			return retryCount >= 2 // stop after 2 retries
		}),
	)
	assert.Nil(t, res)
	assert.Equal(t, failErr, err)
	assert.Equal(t, 2, attempts) // initial + 1 retries
}

func TestDoWithBackoff_LogsRetries(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	writer := bytes.NewBuffer(nil)
	logger := zerolog.New(writer).With().Logger()

	_, _ = retry.DoWithBackoff(
		ctx,
		func() (any, error) {
			attempts++
			return nil, errors.New("fail")
		},
		retry.WithMaxRetries(2),
		retry.WithInitialInterval(1*time.Millisecond),
		retry.WithLogger(logger),
	)

	logs := writer.String()
	assert.Contains(t, logs, "retrying operation")
	assert.Contains(t, logs, "\"retry_count\":1")
	assert.Contains(t, logs, "\"retry_count\":2")
	assert.NotContains(t, logs, "\"retry_count\":3") // no third attempt
}

func TestDoWithBackoff_ContextTimeout(t *testing.T) {
	// total timeout of 10ms
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := retry.DoWithBackoff(
		ctx,
		func() (any, error) {
			time.Sleep(20 * time.Millisecond) // simulate work
			return nil, errors.New("fail")
		},
		retry.WithMaxRetries(100),
		retry.WithInitialInterval(20*time.Millisecond),
	)
	elapsed := time.Since(start)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Greater(t, elapsed, 10*time.Millisecond) // ensure we waited at least the timeout
	assert.Less(t, elapsed, 200*time.Millisecond)   // ensure we don't wait too long
}
