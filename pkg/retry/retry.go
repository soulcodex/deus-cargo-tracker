package retry

import (
	"errors"
	"fmt"
)

// ErrInvalidRetryCount indicates that the retry count provided is less than 1.
var ErrInvalidRetryCount = errors.New("retry count must be >= 1")

// Do executes the provided function fn up to the specified number of attempts.
func Do[T any](fn func() (T, error), attempts int) (T, error) {
	var zero T
	if attempts < 1 {
		return zero, fmt.Errorf("%w: %d", ErrInvalidRetryCount, attempts)
	}

	var result T
	var err error
	for range attempts {
		result, err = fn()
		if err == nil {
			return result, nil
		}
	}
	return zero, err
}
