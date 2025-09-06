package cargodomain

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const cargoNotModifiableErrorMsg = "cargo is not modifiable"

type CargoNotModifiableError struct {
	domain.BaseError
}

func NewCargoNotModifiableError(id CargoID, status Status, isDeleted bool) *CargoNotModifiableError {
	return &CargoNotModifiableError{
		BaseError: domain.NewError(
			cargoNotModifiableErrorMsg,
			errutil.WithMetadataKeyValue("domain.cargo.id", id.String()),
			errutil.WithMetadataKeyValue("domain.cargo.status", status.String()),
			errutil.WithMetadataKeyValue("domain.cargo.deleted", isDeleted),
		),
	}
}
