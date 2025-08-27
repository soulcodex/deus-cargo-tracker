package domainvalidation

// ValidationRule is a function type that checks if the given value
//
//	is equal to the expected value.
type ValidationRule[T any] func(value T) *Error

// Validator is a struct that holds a collection of domain validation rules.
//
// It provides a method to validate a given value against the rules.
type Validator[T any] struct {
	rules []ValidationRule[T]
}

func NewValidator[T any](rules ...ValidationRule[T]) *Validator[T] {
	return &Validator[T]{
		rules: rules,
	}
}

// ErrorGroup represents a collection of domain validation errors.
type ErrorGroup struct {
	errors []Error
}

func (ve *ErrorGroup) Error() string {
	baseMsg := "domain validation error occurred"
	if ve.Empty() {
		return baseMsg
	}

	for _, err := range ve.errors {
		baseMsg += ": " + err.Error()
	}

	return baseMsg
}

func (ve *ErrorGroup) Add(err Error) {
	ve.errors = append(ve.errors, err)
}

func (ve *ErrorGroup) Empty() bool {
	return len(ve.errors) == 0
}

func (ve *ErrorGroup) Errors() []Error {
	return ve.errors
}

func (dv *Validator[T]) Validate(value T) *ErrorGroup {
	allErrors := dv.validate(value)
	if allErrors.Empty() {
		return nil
	}
	return allErrors
}

func (dv *Validator[T]) validate(value T) *ErrorGroup {
	errors := NewValidationErrors()

	for _, rule := range dv.rules {
		if err := rule(value); err != nil {
			errors.Add(*err)
		}
	}

	return errors
}

// NewValidationErrors creates a new instance of ErrorGroup.
func NewValidationErrors() *ErrorGroup {
	return &ErrorGroup{
		errors: make([]Error, 0),
	}
}

// WrapErrors wraps a collection of errors into a single error.
func WrapErrors(errs []Error) *ErrorGroup {
	ve := NewValidationErrors()
	for _, err := range errs {
		ve.Add(err)
	}
	return ve
}
