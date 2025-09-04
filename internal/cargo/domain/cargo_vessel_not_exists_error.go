package cargodomain

import (
	"errors"

	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const cargoVesselNotExistsErrorMsg = "cargo vessel does not exist."

type CargoVesselNotExistsError struct {
	domain.BaseError
}

func NewCargoVesselNotExistsError(id CargoID, vesselID VesselID) *CargoVesselNotExistsError {
	return &CargoVesselNotExistsError{
		BaseError: domain.NewError(
			cargoVesselNotExistsErrorMsg,
			errutil.WithMetadataKeyValue("domain.cargo.vessel_id", vesselID.String()),
			errutil.WithMetadataKeyValue("domain.cargo.id", id.String()),
		),
	}
}

func IsCargoVesselNotExistsError(err error) bool {
	var self *CargoVesselNotExistsError
	return errors.As(err, &self)
}
