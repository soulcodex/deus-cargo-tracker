package cargodomain

import (
	"time"
)

type ItemsPrimitives struct {
	Name   string
	Weight uint64
}

func itemsToPrimitives(items Items) ([]ItemsPrimitives, uint64) {
	primitives := make([]ItemsPrimitives, len(items))
	for i, item := range items {
		primitives[i] = ItemsPrimitives{
			Name:   item.name,
			Weight: item.weight,
		}
	}

	return primitives, items.Weight()
}

type CargoPrimitives struct {
	ID        string
	Items     []ItemsPrimitives
	Status    string
	Weight    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newCargoPrimitives(c *Cargo) CargoPrimitives {
	items, totalWeight := itemsToPrimitives(c.items)

	return CargoPrimitives{
		ID:        c.id.String(),
		Items:     items,
		Status:    c.status.String(),
		Weight:    totalWeight,
		CreatedAt: c.createdAt,
		UpdatedAt: c.updatedAt,
	}
}
