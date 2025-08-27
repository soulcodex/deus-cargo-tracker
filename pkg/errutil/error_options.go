package errutil

// ErrorOpts is a function type that can be used to set options for a new error.
type ErrorOpts func(*BaseError)

// WithSeverity sets the severity of the error.
func WithSeverity(sev ErrorSeverity) ErrorOpts {
	return func(e *BaseError) {
		if sev.IsValid() {
			e.severity = sev
			return
		}

		e.severity = SeverityInfo
	}
}

// WithMetadata sets the metadata of the error.
func WithMetadata(metadata ErrorMetadata) ErrorOpts {
	return func(e *BaseError) {
		if metadata.IsEmpty() {
			return
		}

		e.metadata = metadata
	}
}

// WithCause sets the underlying cause (Unwrap target).
func WithCause(err error) ErrorOpts {
	return func(e *BaseError) {
		if err != nil {
			e.cause = err
		}
	}
}

// WithMetadataKeyValue sets a key-value pair in the metadata of the error.
func WithMetadataKeyValue(key string, value interface{}) ErrorOpts {
	return func(e *BaseError) {
		if e.metadata.IsEmpty() {
			e.metadata.metadata = make(map[string]interface{})
		}

		e.metadata.Set(key, value)
	}
}
