package domainvalidation

import (
	"fmt"
)

var (
	// ErrOutOfBounds is a generic error for out of bounds validation rule.
	ErrOutOfBounds = NewError("given value is out of the expected bounds")

	// ErrUnderMin is a generic error for min validation rule.
	ErrUnderMin = NewError("given value is under the expected minimum")

	// ErrOverMax is a generic error for max validation rule.
	ErrOverMax = NewError("given value is over the expected maximum")
)

// Number is a generic constraint for comparing numeric values.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 |
		~int64 | ~uint | ~uint8 | ~uint16 |
		~uint32 | ~uint64 | ~float32 | ~float64
}

// WithinBounds returns a validation rule that checks if the given value is
// within the expected numeric range.
func WithinBounds[T Number](minValue, maxValue T) ValidationRule[T] {
	return func(value T) *Error {
		if value >= minValue && value <= maxValue {
			return nil
		}

		return ErrOutOfBounds.
			WithRuleName("within_bounds").
			WithRuleValue(value).
			WithRuleExpectedValue(fmt.Sprintf("[%v, %v]", minValue, maxValue))
	}
}

// Min returns a validation rule that checks if the given value is
// less than the expected minimum value.
func Min[T Number](minValue T) ValidationRule[T] {
	return func(value T) *Error {
		if value > minValue {
			return nil
		}

		return ErrUnderMin.
			WithRuleName("min").
			WithRuleValue(value).
			WithRuleExpectedValue(minValue)
	}
}

// Max returns a validation rule that checks if the given value is
// over the expected maximum value.
func Max[T Number](maxValue T) ValidationRule[T] {
	return func(value T) *Error {
		if value < maxValue {
			return nil
		}

		return ErrOverMax.
			WithRuleName("max").
			WithRuleValue(value).
			WithRuleExpectedValue(maxValue)
	}
}
