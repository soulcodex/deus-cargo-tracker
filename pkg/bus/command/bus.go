package querybus

import (
	"github.com/soulcodex/deus-cargoes-tracker/pkg/bus"
)

type Bus = bus.Bus

func InitCommandBus() Bus {
	return bus.InitSyncBus()
}
