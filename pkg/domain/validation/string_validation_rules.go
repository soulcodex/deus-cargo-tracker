package domainvalidation

import (
	"net/mail"
	"net/url"
	"regexp"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/utils"
)

var (
	// ErrLengthUnderMin is a generic error for min length validation rule.
	ErrLengthUnderMin = NewError("string value is too short")

	// ErrLengthOverMax is a generic error for max length validation rule.
	ErrLengthOverMax = NewError("string value is too long")

	// ErrInvalidEmail is a generic error for email validation rule.
	ErrInvalidEmail = NewError("string value is not a valid email address")

	// ErrInvalidURL is a generic error for URL validation rule.
	ErrInvalidURL = NewError("string value is not a valid URL")

	// ErrRegexPatternMismatch is a generic error for regex validation rule.
	ErrRegexPatternMismatch = NewError("string value does not match the given regex pattern")

	// ErrInvalidUUID is a generic error for UUID identifier validation rule.
	ErrInvalidUUID = NewError("string value is not a valid UUID identifier")

	// ErrInvalidULID is a generic error for ULID identifier validation rule.
	ErrInvalidULID = NewError("string value is not a valid ULID identifier")
)

// MinLength returns a domain validation rule that checks if a string value has a minimum length.
func MinLength(minLength int) ValidationRule[string] {
	return func(value string) *Error {
		if len(value) >= minLength {
			return nil
		}

		return ErrLengthUnderMin.
			WithRuleName("min_length").
			WithRuleValue(value).
			WithRuleExpectedValue(minLength)
	}
}

// MaxLength returns a domain validation rule that checks if a string value has a maximum length.
func MaxLength(maxLength int) ValidationRule[string] {
	return func(value string) *Error {
		if len(value) <= maxLength {
			return nil
		}

		return ErrLengthOverMax.
			WithRuleName("max_length").
			WithRuleValue(value).
			WithRuleExpectedValue(maxLength)
	}
}

// Email returns a domain validation rule that checks if a string value is a valid email address.
func Email() ValidationRule[string] {
	return func(value string) *Error {
		if _, emailErr := mail.ParseAddress(value); emailErr == nil {
			return nil
		}

		return ErrInvalidEmail.
			WithRuleName("email").
			WithRuleValue(value).
			WithRuleExpectedValue("valid_email")
	}
}

// URL returns a domain validation rule that checks if a string value is a valid URL.
func URL() ValidationRule[string] {
	return func(value string) *Error {
		u, uriErr := url.Parse(value)

		if uriErr == nil && u.Scheme != "" && u.Host != "" {
			return nil
		}

		return ErrInvalidURL.
			WithRuleName("url").
			WithRuleValue(value).
			WithRuleExpectedValue("valid_url")
	}
}

// Regex returns a domain validation rule that checks if a string value matches a given regex pattern.
func Regex(pattern string) ValidationRule[string] {
	return func(value string) *Error {
		if matched, matchErr := regexp.MatchString(pattern, value); matchErr == nil && matched {
			return nil
		}

		return ErrRegexPatternMismatch.
			WithRuleName("regex").
			WithRuleValue(value).
			WithRuleExpectedValue(pattern)
	}
}

// UUIDIdentifier returns a domain validation rule that checks if a string value is a valid UUID identifier.
func UUIDIdentifier() ValidationRule[string] {
	return func(value string) *Error {
		if err := utils.GuardUUID(value); err == nil {
			return nil
		}

		return ErrInvalidUUID.
			WithRuleName("uuid_identifier_match").
			WithRuleValue(value).
			WithRuleExpectedValue("valid_uuid")
	}
}

// ULIDIdentifier returns a domain validation rule that checks if a string value is a valid ULID identifier.
func ULIDIdentifier() ValidationRule[string] {
	return func(value string) *Error {
		if err := utils.GuardULID(value); err == nil {
			return nil
		}

		return ErrInvalidULID.
			WithRuleName("ulid_identifier_match").
			WithRuleValue(value).
			WithRuleExpectedValue("valid_ulid")
	}
}
