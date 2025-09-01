package querybus

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
)

type Bus = bus.Bus

func InitQueryBus() Bus {
	return bus.InitSyncBus()
}
