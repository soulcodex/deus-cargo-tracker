package cargodomain

import (
	"errors"

	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const cargoNotFoundErrorMsg = "cargo doesn't exist."

type CargoNotExistsError struct {
	domain.BaseError
}

func NewCargoNotExistsError(id CargoID) *CargoNotExistsError {
	return &CargoNotExistsError{
		BaseError: domain.NewError(
			cargoNotFoundErrorMsg,
			errutil.WithMetadataKeyValue("domain.cargo.id", id.String()),
		),
	}
}

func IsCargoNotExistsError(err error) bool {
	var self *CargoNotExistsError
	return errors.As(err, &self)
}
