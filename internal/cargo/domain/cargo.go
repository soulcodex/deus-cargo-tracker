package cargodomain

import (
	"context"
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

func (c *Cargo) Update(
	_ context.Context,
	at time.Time,
	updates ...CargoUpdateOpt,
) error {
	if c.updatedAt.After(at) {
		return nil
	}

	if isDeleted := c.deletedAt != nil; isDeleted {
		return NewCargoNotModifiableError(c.id, c.status, isDeleted)
	}

	for _, update := range updates {
		if updateErr := update(c); updateErr != nil {
			return updateErr
		}
	}

	return nil
}
