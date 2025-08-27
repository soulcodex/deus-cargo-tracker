package bus

import (
	"sync"
)

var _ Bus = (*SyncBus)(nil)
var _ Registerer = (*SyncBus)(nil)

type SyncBus struct {
	handlers map[string]Handler[any, Dto]
	lock     sync.Mutex
}

func InitSyncBus() *SyncBus {
	return &SyncBus{
		handlers: make(map[string]Handler[any, Dto]),
		lock:     sync.Mutex{},
	}
}

func (sb *SyncBus) Register(dto Dto, handler Handler[any, Dto]) error {
	defer sb.lock.Unlock()
	sb.lock.Lock()

	queryName := dto.Type()
	if _, ok := sb.handlers[queryName]; ok {
		return newHandlerAlreadyRegistered(dto, handler)
	}

	sb.handlers[queryName] = handler

	return nil
}

func (sb *SyncBus) GetHandler(dto Dto) (Handler[any, Dto], error) {
	queryName := dto.Type()
	if handler, ok := sb.handlers[queryName]; ok {
		return WrapAsAnyHandler(handler), nil
	}

	return nil, newHandlerNotRegistered(dto)
}
