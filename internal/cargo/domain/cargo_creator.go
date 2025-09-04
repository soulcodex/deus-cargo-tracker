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
	ID       string
	VesselID string
	Items    []CargoItemInput
	At       time.Time
}
type CargoCreator struct {
	repository  CargoRepository
	vesselCheck CargoVesselChecker
}

func NewCargoCreator(repository CargoRepository, checker CargoVesselChecker) *CargoCreator {
	return &CargoCreator{
		repository:  repository,
		vesselCheck: checker,
	}
}

func (cc *CargoCreator) Create(ctx context.Context, input CargoCreateInput) (*Cargo, error) {
	vesselID, err := NewVesselID(input.VesselID)
	if err != nil {
		return nil, err
	}

	if vesselCheckErr := cc.vesselCheck.Check(ctx, vesselID); vesselCheckErr != nil {
		return nil, fmt.Errorf("error checking vessel: %w", vesselCheckErr)
	}

	id, err := NewCargoID(input.ID)
	if err != nil {
		return nil, err
	}

	existing, findErr := cc.repository.Find(ctx, id)
	if findErr != nil && !IsCargoNotExistsError(findErr) {
		return nil, fmt.Errorf("error checking existing cargo: %w", findErr)
	}

	if existing != nil {
		return nil, NewCargoAlreadyExistsError(id, existing.vesselID)
	}

	items := make(Items, len(input.Items))
	for i, itemInput := range input.Items {
		items[i] = newItem(itemInput.Name, itemInput.Weight)
	}

	cargoItems, err := NewItems(items...)
	if err != nil {
		return nil, fmt.Errorf("error creating cargo items: %w", err)
	}

	cargo := NewCargo(id, vesselID, cargoItems, input.At)

	if saveErr := cc.repository.Save(ctx, cargo); saveErr != nil {
		return nil, fmt.Errorf("error saving cargo: %w", saveErr)
	}

	return cargo, nil
}
