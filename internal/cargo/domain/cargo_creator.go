package cargodomain

import (
	"context"
	"fmt"
	"time"
)

type CargoItemInput struct {
	Name   string
	Weight uint64
}
type CargoCreateInput struct {
	ID    string
	Items []CargoItemInput
}
type CargoCreator struct {
	repository CargoRepository
}

func NewCargoCreator(repository CargoRepository) *CargoCreator {
	return &CargoCreator{
		repository: repository,
	}
}

func (cc *CargoCreator) Create(ctx context.Context, input CargoCreateInput, at time.Time) (*Cargo, error) {
	id, err := NewCargoID(input.ID)
	if err != nil {
		return nil, err
	}

	existing, findErr := cc.repository.Find(ctx, id)
	if findErr != nil {
		return nil, fmt.Errorf("error checking existing cargo: %w", findErr)
	}

	if existing != nil {
		return nil, NewCargoAlreadyExistsError(id)
	}

	items := make(Items, len(input.Items))
	for i, itemInput := range input.Items {
		items[i] = newItem(itemInput.Name, itemInput.Weight)
	}

	cargoItems, err := NewItems(items...)
	if err != nil {
		return nil, fmt.Errorf("error creating cargo items: %w", err)
	}

	cargo := NewCargo(id, cargoItems, at)

	if saveErr := cc.repository.Save(ctx, cargo); saveErr != nil {
		return nil, fmt.Errorf("error saving cargo: %w", saveErr)
	}

	return cargo, nil
}
