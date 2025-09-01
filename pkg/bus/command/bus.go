package commandbus

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
)

type Bus = bus.Bus

func InitCommandBus() Bus {
	return bus.InitSyncBus()
}
