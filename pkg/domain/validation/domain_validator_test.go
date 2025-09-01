package domainvalidation_test

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

func TestDomainValidatorWithStringConstraints(t *testing.T) {
	uuidProvider, ulidProvider, stringGenerator :=
		utils.NewFixedUUIDProvider(),
		utils.NewFixedULIDProvider(),
		utils.NewRandomStringGenerator()

	tests := []struct {
		name           string
		constraints    []domainvalidation.ValidationRule[string]
		input          string
		expectedErrors []*domainvalidation.Error
	}{
		{
			name: "it should fail if string is empty",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
			},
			input: "",
			expectedErrors: []*domainvalidation.Error{
				domainvalidation.ErrEmptyStringError,
			},
		},
		{
			name: "it should pass if string isn't empty",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
			},
			input:          uuidProvider.New().String(),
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should fail if not match an specific UUID format",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.UUIDIdentifier(),
			},
			input:          ulidProvider.New().String(),
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrInvalidUUID},
		},
		{
			name: "it should pass if match an specific UUID format",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.UUIDIdentifier(),
			},
			input:          uuidProvider.New().String(),
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should fail if not match an specific ULID format",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.ULIDIdentifier(),
			},
			input:          uuidProvider.New().String(),
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrInvalidULID},
		},
		{
			name: "it should pass if match an specific ULID format",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.ULIDIdentifier(),
			},
			input:          ulidProvider.New().String(),
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should fail if the string is shorten than expected",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.MinLength(5),
			},
			input:          stringGenerator.MustGenerate(4),
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrLengthUnderMin},
		},
		{
			name: "it should fail if the string is longer than expected",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.MaxLength(3),
			},
			input:          stringGenerator.MustGenerate(4),
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrLengthOverMax},
		},
		{
			name: "it should fail if the string doesn't match the pattern",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.Regex("^PT-[a-zA-Z]*"),
			},
			input:          stringGenerator.MustGenerate(6),
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrRegexPatternMismatch},
		},
		{
			name: "it should pass if the string matches the pattern",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.Regex("^PT-[a-zA-Z]*"),
			},
			input:          fmt.Sprintf("PT-%s", stringGenerator.MustGenerate(6)),
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should pass if the string matches the pattern",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.Regex("^PT-[a-zA-Z]*"),
			},
			input:          fmt.Sprintf("PT-%s", stringGenerator.MustGenerate(6)),
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should pass if the string is a valid email",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.Email(),
			},
			input:          "john.doe@mail.com",
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should fail if the string is an invalid email",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.Email(),
			},
			input:          "fake.mail",
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrInvalidEmail},
		},
		{
			name: "it should pass if the string is a valid URL",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.URL(),
			},
			input:          fmt.Sprintf("https://%s.nice-domain.eu", strings.ToLower(stringGenerator.MustGenerate(6))),
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should fail if the string is an invalid URL",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.URL(),
			},
			input:          "fake.url",
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrInvalidURL},
		},
		{
			name: "it should pass if the string is allowed",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.InMap(map[string]struct{}{
					"accepted": {},
					"rejected": {},
				}),
			},
			input:          "accepted",
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should fail if the string is not allowed",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.NotEmpty[string](),
				domainvalidation.InMap(map[string]struct{}{
					"accepted": {},
					"rejected": {},
				}),
			},
			input:          "delayed",
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrValueNotInMap},
		},
		{
			name: "it should fail if the string is not one of the expected values",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.IsOneOf([]string{"accepted", "rejected"}),
			},
			input:          "delayed",
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrValueNotOneOfExpected},
		},
		{
			name: "it should pass if the string is one of the expected values",
			constraints: []domainvalidation.ValidationRule[string]{
				domainvalidation.IsOneOf([]string{"accepted", "rejected"}),
			},
			input:          "accepted",
			expectedErrors: []*domainvalidation.Error{},
		},
	}

	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			validator := domainvalidation.NewValidator(scenario.constraints...)
			errs := validator.Validate(scenario.input)

			if len(scenario.expectedErrors) > 0 {
				require.NotNil(t, errs)
				assert.Len(t, errs.Errors(), len(scenario.expectedErrors))
				for i, expectedError := range scenario.expectedErrors {
					require.ErrorAs(t, errs.Errors()[i], expectedError)
				}
				return
			}

			assert.Nil(t, errs)
		})
	}
}

func TestDomainValidatorWithInt64Constraints(t *testing.T) {
	randomNumSeed := rand.NewPCG(uint64(time.Now().UnixMilli()), uint64(time.Now().UnixMilli()))
	numGenerator := rand.New(randomNumSeed)

	tests := []struct {
		name           string
		constraints    []domainvalidation.ValidationRule[int64]
		input          int64
		expectedErrors []*domainvalidation.Error
	}{
		{
			name: "it should pass if we provide a valid int64 range",
			constraints: []domainvalidation.ValidationRule[int64]{
				domainvalidation.WithinBounds[int64](1, 10),
			},
			input:          int64(1 + numGenerator.IntN(9)),
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should fail if we provide an invalid int64 out of range",
			constraints: []domainvalidation.ValidationRule[int64]{
				domainvalidation.WithinBounds[int64](1, 5),
			},
			input:          int64(10 + numGenerator.IntN(9)),
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrOutOfBounds},
		},
		{
			name: "it should fail if we provide an int64 value under given min limit",
			constraints: []domainvalidation.ValidationRule[int64]{
				domainvalidation.Min[int64](10),
			},
			input:          int64(numGenerator.IntN(9)),
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrUnderMin},
		},
		{
			name: "it should pass if we provide an int64 value under given min limit",
			constraints: []domainvalidation.ValidationRule[int64]{
				domainvalidation.Min[int64](10),
			},
			input:          int64(11 + numGenerator.IntN(9)),
			expectedErrors: []*domainvalidation.Error{},
		},
		{
			name: "it should fail if we provide an int64 value over given max limit",
			constraints: []domainvalidation.ValidationRule[int64]{
				domainvalidation.Max[int64](10),
			},
			input:          int64(10 + numGenerator.IntN(9)),
			expectedErrors: []*domainvalidation.Error{domainvalidation.ErrOverMax},
		},
		{
			name: "it should pass if we provide an int64 value over given max limit",
			constraints: []domainvalidation.ValidationRule[int64]{
				domainvalidation.Max[int64](10),
			},
			input:          int64(numGenerator.IntN(9)),
			expectedErrors: []*domainvalidation.Error{},
		},
	}

	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			validator := domainvalidation.NewValidator(scenario.constraints...)
			errs := validator.Validate(scenario.input)

			if len(scenario.expectedErrors) > 0 {
				require.NotNil(t, errs)
				assert.Len(t, errs.Errors(), len(scenario.expectedErrors))
				for i, expectedError := range scenario.expectedErrors {
					require.ErrorAs(t, errs.Errors()[i], expectedError)
				}
				return
			}

			assert.Nil(t, errs)
		})
	}
}

func TestDomainValidatorWithFloat64Constraints(t *testing.T) {
	randomNumSeed := rand.NewPCG(uint64(time.Now().UnixMilli()), uint64(time.Now().UnixMilli()))
	numGenerator := rand.New(randomNumSeed)

	tests := []struct {
		name           string
		constraints    []domainvalidation.ValidationRule[float64]
		input          float64
		expectedErrors []*domainvalidation.Error
	}{
		{
			name: "it should pass if we provide a valid float64 on range",
			constraints: []domainvalidation.ValidationRule[float64]{
				domainvalidation.WithinBounds[float64](1, 10),
			},
			input:          float64(numGenerator.IntN(9)) + 1.0,
			expectedErrors: []*domainvalidation.Error{},
		},
	}

	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			validator := domainvalidation.NewValidator(scenario.constraints...)
			errs := validator.Validate(scenario.input)

			if len(scenario.expectedErrors) > 0 {
				require.NotNil(t, errs)
				assert.Len(t, errs.Errors(), len(scenario.expectedErrors))
				for i, expectedError := range scenario.expectedErrors {
					require.ErrorAs(t, errs.Errors()[i], expectedError)
				}
				return
			}

			require.Nil(t, errs)
		})
	}
}
