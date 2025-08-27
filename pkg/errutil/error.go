package errutil

import "fmt"

var (
	_ Error             = (*BaseError)(nil)
	_ ErrorWithSeverity = (*BaseError)(nil)
)

// Error is the base interface for all errutil errors.
type Error interface {
	error
	// Metadata returns additional context for the error.
	Metadata() ErrorMetadata
}

// BaseError carries a message, optional cause, severity, and metadata.
type BaseError struct {
	message  string
	cause    error
	severity ErrorSeverity
	metadata ErrorMetadata
}

// NewError creates a BaseError with the given message and applies any ErrorOpts.
func NewError(message string, opts ...ErrorOpts) *BaseError {
	be := &BaseError{
		cause:    nil,
		message:  message,
		severity: SeverityInfo,
		metadata: NewErrorMetadata(),
	}
	for _, opt := range opts {
		opt(be)
	}
	return be
}

// Error returns the message, appending ": <cause.Error()>" if a cause exists.
func (e *BaseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// Unwrap returns the underlying cause (for errors.Is / errors.As).
func (e *BaseError) Unwrap() error {
	return e.cause
}

// Metadata returns the associated metadata.
func (e *BaseError) Metadata() ErrorMetadata {
	return e.metadata
}

// Severity returns how severe this error is.
func (e *BaseError) Severity() ErrorSeverity {
	return e.severity
}

// Wrap attaches a cause to the error in place.
// If err is nil, Wrap is a no-op.
func (e *BaseError) Wrap(err error) *BaseError {
	if err != nil {
		e.cause = err
	}
	return e
}
