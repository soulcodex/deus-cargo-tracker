package cargodomain

import (
	"strings"

	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
)

const (
	StatusPending   Status = "pending"
	StatusInTransit Status = "in_transit"
	StatusDelivered Status = "delivered"
)

var (
	validStatuses = map[Status]struct{}{
		StatusPending:   {},
		StatusInTransit: {},
		StatusDelivered: {},
	}

	ErrInvalidStatusProvided = domainvalidation.NewError("invalid cargo status provided")
)

type Status string

func NewStatus(status string) (Status, error) {
	status = strings.ToLower(status)

	validator := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[Status](),
		domainvalidation.InMap(validStatuses),
	)

	c := Status(status)

	if err := validator.Validate(c); err != nil {
		return "", ErrInvalidStatusProvided.Wrap(err)
	}

	return c, nil
}

func (s Status) IsPending() bool {
	return s == StatusPending
}

func (s Status) String() string {
	return string(s)
}
