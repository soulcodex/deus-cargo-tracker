package retry_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/retry"
)

func TestDo_SuccessFirstAttempt(t *testing.T) {
	res, err := retry.Do(func() (int, error) {
		return 42, nil
	}, 3)

	require.NoError(t, err)
	assert.Equal(t, 42, res)
}

func TestDo_SuccessAfterRetries(t *testing.T) {
	attempts := 0
	res, err := retry.Do(func() (string, error) {
		attempts++
		if attempts < 2 {
			return "", errors.New("try again")
		}
		return "ok", nil
	}, 3)

	require.NoError(t, err)
	assert.Equal(t, "ok", res)
	assert.Equal(t, 2, attempts)
}

func TestDo_FailureAllAttempts(t *testing.T) {
	attempts := 0
	masterErr := errors.New("always fail")

	res, err := retry.Do(func() (int, error) {
		attempts++
		return 0, masterErr
	}, 3)

	require.Error(t, err)
	assert.Equal(t, masterErr, err)
	assert.Equal(t, 0, res) // zero-value of int
	assert.Equal(t, 3, attempts)
}

func TestDo_InvalidAttempts(t *testing.T) {
	_, err := retry.Do(func() (int, error) {
		return 1, nil
	}, 0)

	assert.ErrorIs(t, err, retry.ErrInvalidRetryCount)
}

func TestDo_CustomType(t *testing.T) {
	type item struct{ A int }
	attempts := 0

	res, err := retry.Do(func() (item, error) {
		attempts++
		if attempts < 3 {
			return item{}, errors.New("not yet")
		}
		return item{A: 5}, nil
	}, 5)

	require.NoError(t, err)
	assert.Equal(t, item{A: 5}, res)
	assert.Equal(t, 3, attempts)
}
