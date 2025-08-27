package bus

import (
	"errors"
)

var (
	errUnexpectedBusProvided = errors.New("unexpected bus provided expected registerer")
)

type Registerer interface {
	Bus
	Register(command Dto, handler Handler[any, Dto]) error
}

func MustRegister[Output, Input any](bus Bus, dto Dto, handler Handler[Output, Input]) {
	registerer, isRegisterer := bus.(Registerer)
	if !isRegisterer {
		panic(errUnexpectedBusProvided)
	}

	if err := registerer.Register(dto, WrapAsAnyHandler(handler)); err != nil {
		panic(err)
	}
}
