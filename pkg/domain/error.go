package domain

import (
	"fmt"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

// BaseError is the base domain error interface for all errors in the project.
type BaseError struct {
	*errutil.BaseError
}

func (e BaseError) Wrap(err error) BaseError {
	e.BaseError = e.BaseError.Wrap(err)
	return e
}

// NewError creates a new domain error with the given message and options.
func NewError(message string, opts ...errutil.ErrorOpts) BaseError {
	return BaseError{
		BaseError: errutil.NewError(message, opts...),
	}
}

// NewErrorf creates a new domain error with the given formatted message and options.
func NewErrorf(format string, args ...interface{}) BaseError {
	return BaseError{
		BaseError: errutil.NewError(fmt.Sprintf(format, args...)),
	}
}

// NewErrorMetadata creates a new domain error metadata.
func NewErrorMetadata() errutil.ErrorMetadata {
	return errutil.NewErrorMetadata()
}
