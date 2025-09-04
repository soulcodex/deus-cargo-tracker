package cargodomain

import (
	"time"
)

type Cargo struct {
	id        CargoID
	vesselID  VesselID
	items     Items
	status    Status
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewCargo(id CargoID, vesselID VesselID, items Items, at time.Time) *Cargo {
	return &Cargo{
		id:        id,
		vesselID:  vesselID,
		items:     items,
		status:    StatusPending,
		createdAt: at,
		updatedAt: at,
		deletedAt: nil,
	}
}

func NewCargoFromPrimitives(p CargoPrimitives) *Cargo {
	items := make(Items, len(p.Items))
	for i, item := range p.Items {
		items[i] = newItem(item.Name, item.Weight)
	}

	return &Cargo{
		id:        CargoID(p.ID),
		vesselID:  VesselID(p.VesselID),
		items:     items,
		status:    Status(p.Status),
		createdAt: p.CreatedAt,
		updatedAt: p.UpdatedAt,
		deletedAt: p.DeletedAt,
	}
}

func (c *Cargo) Primitives() CargoPrimitives {
	return newCargoPrimitives(c)
}

func (c *Cargo) ID() CargoID {
	return c.id
}

func (c *Cargo) VesselID() VesselID {
	return c.vesselID
}

func (c *Cargo) AppendItem(item Item, at time.Time) error {
	if !c.status.IsPending() || c.deletedAt != nil {
		return NewCargoNotModifiableError(c.id, c.status)
	}

	items, err := NewItems(append(c.items, item)...)
	if err != nil {
		return err
	}

	c.items = items
	c.updatedAt = at

	return nil
}
