package cargodomain

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
)

var (
	ErrInvalidCargoOptionProvided = domain.NewError("invalid cargo update value provided")
)

type CargoUpdateOpt func(*Cargo) error

func WithStatus(status string) CargoUpdateOpt {
	return func(c *Cargo) error {
		newStatus, err := NewStatus(status)
		if err != nil {
			return ErrInvalidCargoOptionProvided.Wrap(err)
		}

		if !c.status.IsTransitionAllowed(newStatus) {
			return ErrStatusTransitionNotAllowed
		}

		c.status = newStatus

		return nil
	}
}
