package cargodomain

import (
	"context"
	"time"

	cargotrackingdomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/tracking"
)

type Cargo struct {
	id        CargoID
	vesselID  VesselID
	items     Items
	tracking  cargotrackingdomain.Tracking
	status    Status
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewCargo(
	id CargoID,
	vesselID VesselID,
	trackingID cargotrackingdomain.TrackingID,
	items Items,
	at time.Time,
) *Cargo {
	cargo := &Cargo{
		id:        id,
		vesselID:  vesselID,
		items:     items,
		tracking:  make(cargotrackingdomain.Tracking, 0),
		status:    StatusPending,
		createdAt: at,
		updatedAt: at,
		deletedAt: nil,
	}

	cargo.appendTracking(cargotrackingdomain.NewTrackingOnCargoCreated(trackingID, cargo.status.String(), at))

	return cargo
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
		tracking:  make(cargotrackingdomain.Tracking, 0),
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

func (c *Cargo) appendTracking(item cargotrackingdomain.TrackingItem) {
	if c.tracking == nil {
		c.tracking = make(cargotrackingdomain.Tracking, 0)
	}

	c.tracking = append(c.tracking, item)
}

func (c *Cargo) Update(_ context.Context, updates ...CargoUpdateOpt) error {
	// Prevent updates if cargo is deleted
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
