package domainvalidation

var (
	// ErrEmptyStringError is a generic error for empty string validation rule.
	ErrEmptyStringError = NewError("given value should not be empty")

	// ErrValueNotOneOfExpected is a generic error for one of validation rule.
	ErrValueNotOneOfExpected = NewError("given value is not one of the expected values")

	// ErrValueNotInMap is a generic error for in map validation rule.
	ErrValueNotInMap = NewError("given value is not in given map")

	// ErrNotInRange is a generic error for range validation rule.
	ErrNotInRange = NewError("given value is not in the expected range")
)

// NotEmpty is a generic validation rule that checks if the given value is empty.
func NotEmpty[T comparable]() ValidationRule[T] {
	return func(value T) *Error {
		var emptyValue T
		if value != emptyValue {
			return nil
		}

		return ErrEmptyStringError.
			WithRuleName("is_empty").
			WithRuleValue(value).
			WithRuleExpectedValue("not_empty")
	}
}

// NotZero is a generic validation rule that checks if the given value is zero.
func NotZero[T interface{ IsZero() bool }]() ValidationRule[T] {
	return func(value T) *Error {
		if !value.IsZero() {
			return nil
		}

		return ErrNotInRange.
			WithRuleName("is_not_zero").
			WithRuleValue(value).
			WithRuleExpectedValue("not_zero")
	}
}

// IsOneOf is a generic validation rule that checks if the given value is one of the expected values.
func IsOneOf[S ~[]E, E comparable](values S) ValidationRule[E] {
	return func(value E) *Error {
		for _, v := range values {
			if v == value {
				return nil
			}
		}

		return ErrValueNotOneOfExpected.
			WithRuleName("is_one_of").
			WithRuleValue(value).
			WithRuleExpectedValue(values)
	}
}

// InMap is a generic validation rule that checks if the given value is in the given map.
func InMap[S ~map[K]struct{}, K comparable](values S) ValidationRule[K] {
	return func(value K) *Error {
		if _, ok := values[value]; ok {
			return nil
		}

		return ErrValueNotInMap.
			WithRuleName("in_map").
			WithRuleValue(value).
			WithRuleExpectedValue(values)
	}
}
