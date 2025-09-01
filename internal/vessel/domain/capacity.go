package vesseldomain

import (
	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

const (
	vesselKilogramsMaxCapacity = 100_000 // 100 tons
)

var (
	ErrInvalidCapacityProvided = errutil.NewError("invalid capacity provided")
)

type Capacity uint64

func NewCapacity(capacity uint64) (Capacity, error) {
	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[uint64](),
		domainvalidation.Max[uint64](vesselKilogramsMaxCapacity),
	)

	if err := validation.Validate(capacity); err != nil {
		return 0, ErrInvalidCapacityProvided.Wrap(err)
	}

	return Capacity(capacity), nil
}

func (c Capacity) Value() uint64 {
	return uint64(c)
}
