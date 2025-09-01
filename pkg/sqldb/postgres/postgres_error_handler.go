package postgres

import (
	"errors"

	"github.com/lib/pq"
)

const (
	UniqueViolationErrorCode = pq.ErrorCode("23505")
)

type ErrorHandlers map[pq.ErrorCode]ErrorHandlerFunc
type ErrorHandlerFunc func(resource interface{}, err *pq.Error) error

type ErrorHandler struct {
	handlers ErrorHandlers
}

func NewErrorHandler(handlers map[pq.ErrorCode]ErrorHandlerFunc) *ErrorHandler {
	return &ErrorHandler{handlers: handlers}
}

func (eh *ErrorHandler) Handle(resource interface{}, err *pq.Error) error {
	if handler, exists := eh.handlers[err.Code]; exists {
		return handler(resource, err)
	}

	return err
}

func IsPostgresError(err error) (*pq.Error, bool) {
	var self *pq.Error
	return self, errors.As(err, &self)
}
