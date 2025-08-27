package bus

import (
	"context"

	dsync "github.com/soulcodex/deus-cargoes-tracker/pkg/distributed-sync"
	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

type DispatchWithOutputFunc[Input Dto, Output any] func(context.Context, Input) (Output, error)
type DispatchBlockingFunc func(context.Context, BlockingDto) error
type DispatchFunc func(context.Context, Dto) error

func DispatchWithResponse[Input Dto, Output any](bus Bus) DispatchWithOutputFunc[Input, Output] {
	return func(ctx context.Context, input Input) (Output, error) {
		var output Output

		handler, err := bus.GetHandler(input)
		if err != nil {
			return output, ErrNoHandlerForInput(input, err)
		}

		out, err := handler.Handle(ctx, input)
		if err != nil {
			return output, ErrUnprocessableHandler(input, err)
		}

		response, ok := out.(Output)
		if !ok {
			return output, newInvalidOutputReceived(handler, output, out)
		}

		return response, nil
	}
}

func DispatchBlocking(bus Bus, mutex dsync.MutexService) DispatchBlockingFunc {
	return func(ctx context.Context, input BlockingDto) error {
		handler, err := bus.GetHandler(input)
		if err != nil {
			return ErrNoHandlerForInput(input, err)
		}

		operation := func() (interface{}, error) {
			return handler.Handle(ctx, input)
		}

		_, blockingErr := mutex.Mutex(ctx, input.BlockingKey(), operation)
		if blockingErr != nil {
			return ErrUnprocessableHandler(input, blockingErr)
		}

		return nil
	}
}

func Dispatch(bus Bus) DispatchFunc {
	return func(ctx context.Context, input Dto) error {
		handler, err := bus.GetHandler(input)
		if err != nil {
			return ErrNoHandlerForInput(input, err)
		}

		_, handleErr := handler.Handle(ctx, input)
		if handleErr != nil {
			return ErrUnprocessableHandler(input, handleErr)
		}

		return nil
	}
}

func ErrNoHandlerForInput(input Dto, previous error) error {
	return errutil.NewError(
		"unable to get handler for input",
		errutil.WithMetadataKeyValue("bus.operation.input_type", input.Type())).
		Wrap(previous)
}

func ErrUnprocessableHandler(input Dto, previous error) error {
	return errutil.NewError(
		"unprocessable bus input",
		errutil.WithMetadataKeyValue("bus.operation.input_type", input.Type())).
		Wrap(previous)
}
