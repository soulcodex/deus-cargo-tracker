package cargodomain

import (
	"time"
)

type Cargo struct {
	id        CargoID
	items     Items
	status    Status
	createdAt time.Time
	updatedAt time.Time
}

func NewCargo(id CargoID, items Items, at time.Time) *Cargo {
	return &Cargo{
		id:        id,
		items:     items,
		status:    StatusPending,
		createdAt: at,
		updatedAt: at,
	}
}

func (c *Cargo) Primitives() CargoPrimitives {
	return newCargoPrimitives(c)
}

func (c *Cargo) ID() CargoID {
	return c.id
}

func (c *Cargo) AppendItem(item Item, at time.Time) error {
	if !c.status.IsPending() {
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
