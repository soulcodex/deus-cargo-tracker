package vesseldomain

import (
	"errors"

	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const vesselNotFoundErrorMsg = "vessel doesn't exist."

type VesselNotExistsError struct {
	domain.BaseError
}

func NewVesselNotExistsError(id VesselID) *VesselNotExistsError {
	return &VesselNotExistsError{
		BaseError: domain.NewError(
			vesselNotFoundErrorMsg,
			errutil.WithMetadataKeyValue("domain.error.vessel_id", id.String()),
		),
	}
}

func IsVesselNotExistsError(err error) bool {
	var self *VesselNotExistsError
	return errors.As(err, &self)
}
