package vesseldomain

import (
	"time"
)

type Vessel struct {
	id        VesselID
	name      Name
	capacity  Capacity
	location  Location
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewVesselFromPrimitives(v VesselPrimitives) *Vessel {
	return &Vessel{
		id:        VesselID(v.ID),
		name:      Name(v.Name),
		capacity:  Capacity(v.Capacity),
		location:  Location{latitude: v.Latitude, longitude: v.Longitude},
		createdAt: v.CreatedAt,
		updatedAt: v.UpdatedAt,
		deletedAt: v.DeletedAt,
	}
}

func (v *Vessel) Primitives() VesselPrimitives {
	return newVesselPrimitives(v)
}

func (v *Vessel) ID() VesselID {
	return v.id
}
