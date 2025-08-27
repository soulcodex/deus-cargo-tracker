package bus

import (
	"context"
)

type HandlerFunc[Output, Input any] func(context.Context, Input) (Output, error)

type Handler[Output, Input any] interface {
	Handle(context.Context, Input) (Output, error)
}
