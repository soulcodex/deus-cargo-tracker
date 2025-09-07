package cargoqueries

import (
	"time"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
)

type CargoResponseItem struct {
	Name   string
	Weight uint64
}

type CargoTrackingResponseItem struct {
	ID           string
	EntryType    string
	StatusBefore *string
	StatusAfter  *string
	CreatedAt    time.Time
}
type CargoResponse struct {
	ID        string
	VesselID  string
	Items     []CargoResponseItem
	Tracking  []CargoTrackingResponseItem
	Status    string
	Weight    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCargoResponse(p cargodomain.CargoPrimitives) CargoResponse {
	cargoItems := make([]CargoResponseItem, len(p.Items))
	for i, item := range p.Items {
		cargoItems[i] = CargoResponseItem{
			Name:   item.Name,
			Weight: item.Weight,
		}
	}

	trackingItems := make([]CargoTrackingResponseItem, len(p.Tracking))
	for i, item := range p.Tracking {
		trackingItems[i] = CargoTrackingResponseItem{
			ID:           item.ID,
			EntryType:    item.EntryType,
			StatusBefore: item.StatusBefore,
			StatusAfter:  item.StatusAfter,
			CreatedAt:    item.CreatedAt,
		}
	}

	return CargoResponse{
		ID:        p.ID,
		VesselID:  p.VesselID,
		Items:     cargoItems,
		Tracking:  trackingItems,
		Status:    p.Status,
		Weight:    p.Weight,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
