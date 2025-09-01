package vesseldomain

import (
	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
)

var (
	minMaxLatitude  = [2]float64{-90, 90}
	minMaxLongitude = [2]float64{-180, 180}

	ErrInvalidLocationProvided = domainvalidation.NewError("invalid location provided")
)

type Location struct {
	latitude  float64
	longitude float64
}

func NewLocation(latitude, longitude float64) (Location, error) {
	c := Location{
		latitude:  latitude,
		longitude: longitude,
	}

	validator := domainvalidation.NewValidator(
		func(value *Location) *domainvalidation.Error {
			return domainvalidation.WithinBounds(minMaxLatitude[0], minMaxLatitude[1])(value.latitude)
		},
		func(value *Location) *domainvalidation.Error {
			return domainvalidation.WithinBounds(minMaxLongitude[0], minMaxLongitude[1])(value.longitude)
		},
	)

	if err := validator.Validate(&c); err != nil {
		return Location{}, ErrInvalidLocationProvided.Wrap(err)
	}

	return c, nil
}
