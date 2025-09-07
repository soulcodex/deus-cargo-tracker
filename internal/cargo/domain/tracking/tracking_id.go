package cargotrackingdomain

import (
	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

var (
	ErrInvalidTrackingIDProvided = errutil.NewError("invalid tracking id provided")
)

type TrackingID string

func NewTrackingID(id string) (TrackingID, error) {
	trackingID := TrackingID(id)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[string](),
		domainvalidation.ULIDIdentifier(),
	)

	if err := validation.Validate(id); err != nil {
		return "", ErrInvalidTrackingIDProvided.Wrap(err)
	}

	return trackingID, nil
}

func (t TrackingID) String() string {
	return string(t)
}
