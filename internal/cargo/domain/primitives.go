package cargodomain

import (
	"time"
)

type ItemsPrimitives struct {
	Name   string `json:"name"`
	Weight uint64 `json:"weight"`
}

func itemsToPrimitives(items Items) []ItemsPrimitives {
	primitives := make([]ItemsPrimitives, len(items))
	for i, item := range items {
		primitives[i] = ItemsPrimitives{
			Name:   item.name,
			Weight: item.weight,
		}
	}

	return primitives
}

type CargoPrimitives struct {
	ID        string
	VesselID  string
	Items     []ItemsPrimitives
	Status    string
	Weight    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func newCargoPrimitives(c *Cargo) CargoPrimitives {
	items := itemsToPrimitives(c.items)

	return CargoPrimitives{
		ID:        c.id.String(),
		VesselID:  c.vesselID.String(),
		Items:     items,
		Status:    c.status.String(),
		Weight:    c.items.Weight(),
		CreatedAt: c.createdAt,
		UpdatedAt: c.updatedAt,
		DeletedAt: c.deletedAt,
	}
}
