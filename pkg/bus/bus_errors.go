package bus

import (
	"fmt"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

const (
	busDtoSemConvKey     = "messaging.bus.dto"
	busHandlerSemConvKey = "messaging.bus.handler"

	busOutputExpectedSemConvKey = "messaging.bus.output.expected"
	busOutputReceivedSemConvKey = "messaging.bus.output.received"
)

type HandlerAlreadyRegisteredError struct {
	*errutil.BaseError
}

func newHandlerAlreadyRegistered(dto Dto, handler Handler[any, Dto]) *HandlerAlreadyRegisteredError {
	return &HandlerAlreadyRegisteredError{
		BaseError: errutil.NewError(
			"bus handler already registered",
			errutil.WithMetadataKeyValue(busDtoSemConvKey, fmt.Sprintf("%T", dto)),
			errutil.WithMetadataKeyValue(busHandlerSemConvKey, fmt.Sprintf("%T", handler)),
		),
	}
}

type HandlerNotRegisteredError struct {
	*errutil.BaseError
}

func newHandlerNotRegistered(dto Dto) *HandlerNotRegisteredError {
	return &HandlerNotRegisteredError{
		BaseError: errutil.NewError(
			"bus handler not registered",
			errutil.WithMetadataKeyValue(busDtoSemConvKey, fmt.Sprintf("%T", dto)),
		),
	}
}

type InvalidDtoProvidedError struct {
	*errutil.BaseError
}

func NewInvalidDtoProvided[Output, Input any](dto Dto, handler Handler[Output, Input]) *InvalidDtoProvidedError {
	return &InvalidDtoProvidedError{
		BaseError: errutil.NewError(
			"invalid dto provided",
			errutil.WithMetadataKeyValue(busDtoSemConvKey, fmt.Sprintf("%T", dto)),
			errutil.WithMetadataKeyValue(busHandlerSemConvKey, fmt.Sprintf("%T", handler)),
		),
	}
}

type InvalidOutputReceivedError struct {
	*errutil.BaseError
}

func newInvalidOutputReceived(handler Handler[any, Dto], expected, received interface{}) *InvalidOutputReceivedError {
	return &InvalidOutputReceivedError{
		BaseError: errutil.NewError(
			"invalid output received by handler",
			errutil.WithMetadataKeyValue(busHandlerSemConvKey, fmt.Sprintf("%T", handler)),
			errutil.WithMetadataKeyValue(busOutputExpectedSemConvKey, fmt.Sprintf("%T", expected)),
			errutil.WithMetadataKeyValue(busOutputReceivedSemConvKey, fmt.Sprintf("%T", received)),
		),
	}
}
