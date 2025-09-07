package cargodomain

import (
	"time"

	cargotrackingdomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/tracking"
	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
)

var (
	ErrInvalidCargoOptionProvided = domain.NewError("invalid cargo update value provided")
)

type CargoUpdateOpt func(*Cargo) error

func WithStatus(trackingID, status string, at time.Time) CargoUpdateOpt {
	return func(c *Cargo) error {
		newStatus, err := NewStatus(status)
		if err != nil {
			return ErrInvalidCargoOptionProvided.Wrap(err)
		}

		if c.status.Equals(newStatus) {
			return ErrStatusUnchanged
		}

		if !c.status.IsTransitionAllowed(newStatus) {
			return ErrStatusTransitionNotAllowed
		}

		newTrackingID, err := cargotrackingdomain.NewTrackingID(trackingID)
		if err != nil {
			return ErrInvalidCargoOptionProvided.Wrap(err)
		}

		tracking := cargotrackingdomain.NewTrackingOnCargoStatusChanged(
			newTrackingID,
			at,
			c.status.String(),
			newStatus.String(),
		)
		c.appendTracking(tracking)

		c.status = newStatus
		c.updatedAt = at

		return nil
	}
}
