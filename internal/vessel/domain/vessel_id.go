package vesseldomain

import (
	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

var (
	ErrInvalidVesselIDProvided = errutil.NewError("invalid vessel id provided")
)

type VesselID string

func NewVesselID(id string) (VesselID, error) {
	vesselID := VesselID(id)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[string](),
		domainvalidation.ULIDIdentifier(),
	)

	if err := validation.Validate(id); err != nil {
		return "", ErrInvalidVesselIDProvided.Wrap(err)
	}

	return vesselID, nil
}

func (c VesselID) String() string {
	return string(c)
}
