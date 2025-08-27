package errutil

import (
	"errors"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

const (
	criticalErrorTypeStr = "critical"
)

// CriticalError is a critical error that should stop the execution of the program.
type CriticalError struct {
	*BaseError
}

// NewCriticalError creates a new CriticalError with the given message.
func NewCriticalError(msg string) *CriticalError {
	errType := criticalErrorType()
	return &CriticalError{
		BaseError: NewError(
			msg,
			WithSeverity(SeverityFatal),
			WithMetadataKeyValue(string(errType.Key), errType.Value.AsString()),
		),
	}
}

// NewCriticalErrorWithMetadata creates a new CriticalError with the given message and metadata.
func NewCriticalErrorWithMetadata(msg string, metadata ErrorMetadata) *CriticalError {
	errType := criticalErrorType()
	return &CriticalError{
		BaseError: NewError(
			msg,
			WithSeverity(SeverityFatal),
			WithMetadata(metadata),
			WithMetadataKeyValue(string(errType.Key), errType.Value.AsString()),
		),
	}
}

func criticalErrorType() attribute.KeyValue {
	return semconv.ErrorTypeKey.String(criticalErrorTypeStr)
}

// IsCriticalError reports whether err is a CriticalError.
func IsCriticalError(err error) bool {
	var ce *CriticalError
	return errors.As(err, &ce)
}
