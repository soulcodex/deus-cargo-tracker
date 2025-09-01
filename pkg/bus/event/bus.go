package eventbus

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
)

type Bus = bus.Bus

func InitEventBus() Bus {
	return bus.InitSyncBus()
}
