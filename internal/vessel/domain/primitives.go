package vesseldomain

import (
	"time"
)

type VesselPrimitives struct {
	ID        string
	Name      string
	Capacity  uint64
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func newVesselPrimitives(v *Vessel) VesselPrimitives {
	return VesselPrimitives{
		ID:        v.id.String(),
		Name:      v.name.String(),
		Capacity:  v.capacity.Value(),
		Latitude:  v.location.latitude,
		Longitude: v.location.longitude,
		CreatedAt: v.createdAt,
		UpdatedAt: v.updatedAt,
		DeletedAt: v.deletedAt,
	}
}
