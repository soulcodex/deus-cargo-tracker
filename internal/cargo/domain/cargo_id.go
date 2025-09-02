package cargodomain

import (
	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

var (
	ErrInvalidCargoIDProvided = errutil.NewError("invalid cargo id provided")
)

type CargoID string

func NewCargoID(id string) (CargoID, error) {
	cargoID := CargoID(id)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[string](),
		domainvalidation.ULIDIdentifier(),
	)

	if err := validation.Validate(id); err != nil {
		return "", ErrInvalidCargoIDProvided.Wrap(err)
	}

	return cargoID, nil
}

func (c CargoID) String() string {
	return string(c)
}
