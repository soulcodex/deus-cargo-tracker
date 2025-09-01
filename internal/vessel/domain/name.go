package vesseldomain

import (
	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const (
	vesselNameMaxLength = 100
)

var (
	ErrInvalidVesselNameProvided = errutil.NewError("invalid vessel name provided")
)

type Name string

func NewName(name string) (Name, error) {
	vesselName := Name(name)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[string](),
		domainvalidation.MaxLength(vesselNameMaxLength),
	)

	if err := validation.Validate(name); err != nil {
		return "", ErrInvalidVesselNameProvided.Wrap(err)
	}

	return vesselName, nil
}

func (n Name) String() string {
	return string(n)
}
