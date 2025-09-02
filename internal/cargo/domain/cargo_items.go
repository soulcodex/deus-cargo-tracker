package cargodomain

import (
	domainvalidation "github.com/soulcodex/deus-cargo-tracker/pkg/domain/validation"
)

var (
	ErrInvalidItemsProvided = domainvalidation.NewError("invalid items provided")
)

const (
	MaxItemsPerCargo int    = 10
	MinItemsPerCargo int    = 1
	MaxItemWeight    uint64 = 10000 // in grams
	MinItemWeight    uint64 = 1     // in grams
	MaxCargoWeight   uint64 = 50000 // in grams
	MinCargoWeight   uint64 = 100   // in grams
	MinItemNameLen   int    = 1
	MaxItemNameLen   int    = 255
)

type Items []Item

func NewItems(items ...Item) (Items, error) {
	validator := itemsValidator()

	if err := validator.Validate(items); err != nil {
		return nil, ErrInvalidItemsProvided.Wrap(err)
	}

	return items, nil
}

func (i Items) Len() int {
	return len(i)
}

func (i Items) Weight() uint64 {
	var total uint64

	for _, item := range i {
		total += item.weight
	}

	return total
}

func itemsValidator() *domainvalidation.Validator[Items] {
	return domainvalidation.NewValidator(
		func(i Items) *domainvalidation.Error {
			err := domainvalidation.WithinBounds(MinItemsPerCargo, MaxItemsPerCargo)(i.Len())
			if err != nil {
				return domainvalidation.NewError("number of items is invalid").Wrap(err)
			}

			return nil
		},
		func(i Items) *domainvalidation.Error {
			err := domainvalidation.WithinBounds(MinCargoWeight, MaxCargoWeight)(i.Weight())
			if err != nil {
				return domainvalidation.NewError("total weight of items is invalid").Wrap(err)
			}

			return nil
		},
		func(i Items) *domainvalidation.Error {
			for _, item := range i {
				if err := domainvalidation.NewValidator(
					domainvalidation.NotEmpty[string](),
					domainvalidation.MinLength(MinItemNameLen),
					domainvalidation.MaxLength(MaxItemNameLen),
				).Validate(item.name); err != nil {
					return domainvalidation.NewError("item name is invalid").Wrap(err)
				}

				if err := domainvalidation.NewValidator(
					domainvalidation.WithinBounds(MinItemWeight, MaxItemWeight),
				).Validate(item.weight); err != nil {
					return domainvalidation.NewError("item weight is invalid").Wrap(err)
				}
			}
			return nil
		},
	)
}

type Item struct {
	name   string
	weight uint64
}

func newItem(name string, weight uint64) Item {
	return Item{
		name:   name,
		weight: weight,
	}
}
