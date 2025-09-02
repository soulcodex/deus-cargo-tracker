package vesseldomain

import (
	"github.com/soulcodex/deus-cargo-tracker/pkg/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const vesselAlreadyExistsErrorMsg = "vessel already exists."

type VesselAlreadyExistsError struct {
	domain.BaseError
}

func NewVesselAlreadyExistsError(id VesselID) *VesselAlreadyExistsError {
	return &VesselAlreadyExistsError{
		BaseError: domain.NewError(
			vesselAlreadyExistsErrorMsg,
			errutil.WithMetadataKeyValue("domain.vesel.id", id.String()),
		),
	}
}
