package vesselqueries

import (
	"time"

	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
)

type VesselResponse struct {
	ID        string
	Name      string
	Capacity  uint64
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewVesselResponse(p vesseldomain.VesselPrimitives) VesselResponse {
	return VesselResponse{
		ID:        p.ID,
		Name:      p.Name,
		Capacity:  p.Capacity,
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
