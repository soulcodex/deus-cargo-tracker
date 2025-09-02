package cargodomain

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const cargoAlreadyExistsErrorMsg = "cargo already exists."

type CargoAlreadyExistsError struct {
	domain.BaseError
}

func NewCargoAlreadyExistsError(id CargoID) *CargoAlreadyExistsError {
	return &CargoAlreadyExistsError{
		BaseError: domain.NewError(
			cargoAlreadyExistsErrorMsg,
			errutil.WithMetadataKeyValue("domain.cargo.vessel_id", id.String()),
		),
	}
}
