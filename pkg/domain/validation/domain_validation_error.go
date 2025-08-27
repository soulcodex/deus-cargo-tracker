package domainvalidation

import (
	"github.com/soulcodex/deus-cargoes-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

// domainValidationErrorMetadataKey is the key used to store the validation error items in the metadata.
const domainValidationErrorMetadataKey = "domain_validation"

// Error is the base domain validation error type.
type Error struct {
	domain.BaseError
}

func (e *Error) Wrap(err error) *Error {
	e.BaseError = e.BaseError.Wrap(err)
	return e
}

func (e *Error) WithRuleName(ruleName string) *Error {
	metadata := e.Metadata()
	verror := e.getValidationRuleErrorFromMetadata(&metadata)
	verror.RuleName = ruleName
	metadata.Set(domainValidationErrorMetadataKey, verror)
	return e
}

func (e *Error) WithRuleValue(value any) *Error {
	metadata := e.Metadata()
	verror := e.getValidationRuleErrorFromMetadata(&metadata)
	verror.Value = value
	metadata.Set(domainValidationErrorMetadataKey, verror)
	return e
}

func (e *Error) WithRuleExpectedValue(expectedValue any) *Error {
	metadata := e.Metadata()
	verror := e.getValidationRuleErrorFromMetadata(&metadata)
	verror.ExpectedValue = expectedValue
	metadata.Set(domainValidationErrorMetadataKey, verror)
	return e
}

func (e *Error) getValidationRuleErrorFromMetadata(metadata *errutil.ErrorMetadata) *ValidationRuleError {
	if metadata.IsEmpty() {
		*metadata = domain.NewErrorMetadata()
	}

	ruleError, ok := metadata.Get(domainValidationErrorMetadataKey).(*ValidationRuleError)
	if !ok {
		return &ValidationRuleError{}
	}

	return ruleError
}

type ValidationRuleError struct {
	RuleName      string
	Value         any
	ExpectedValue any
}

// NewValidationRuleError creates a new domain validation error item with the given values.
func NewValidationRuleError(ruleName string, value any, expectedValue any) ValidationRuleError {
	return ValidationRuleError{
		RuleName:      ruleName,
		Value:         value,
		ExpectedValue: expectedValue,
	}
}

// NewError creates a new domain validation error with the given message and validation rule error.
func NewError(msg string) *Error {
	return &Error{
		BaseError: domain.NewError(msg),
	}
}
