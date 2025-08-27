package bus

import (
	"context"
)

type AnyHandlerWrapper struct {
	handler HandlerFunc[any, Dto]
}

func (ahw AnyHandlerWrapper) Handle(ctx context.Context, input Dto) (any, error) {
	return ahw.handler(ctx, input)
}

func WrapAsAnyHandler[Output, Input any](handler Handler[Output, Input]) Handler[any, Dto] {
	return AnyHandlerWrapper{handler: func(ctx context.Context, input Dto) (any, error) {
		castedInput, inputMatch := input.(Input)
		if !inputMatch {
			return nil, NewInvalidDtoProvided(input, handler)
		}

		return handler.Handle(ctx, castedInput)
	}}
}
