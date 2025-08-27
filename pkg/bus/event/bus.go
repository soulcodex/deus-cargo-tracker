package eventbus

import (
	"github.com/soulcodex/deus-cargoes-tracker/pkg/bus"
)

type Bus = bus.Bus

func InitEventBus() Bus {
	return bus.InitSyncBus()
}
