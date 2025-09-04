package cargoqueries

import (
	"time"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
)

type CargoResponseItems struct {
	Name   string
	Weight uint64
}
type CargoResponse struct {
	ID        string
	VesselID  string
	Items     []CargoResponseItems
	Status    string
	Weight    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCargoResponse(p cargodomain.CargoPrimitives) CargoResponse {
	cargoItems := make([]CargoResponseItems, len(p.Items))
	for i, item := range p.Items {
		cargoItems[i] = CargoResponseItems{
			Name:   item.Name,
			Weight: item.Weight,
		}
	}

	return CargoResponse{
		ID:        p.ID,
		VesselID:  p.VesselID,
		Items:     cargoItems,
		Status:    p.Status,
		Weight:    p.Weight,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
