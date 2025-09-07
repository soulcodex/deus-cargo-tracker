package cargodomain

import (
	"context"
	"fmt"
	"time"

	cargotrackingdomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/tracking"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

type CargoItemInput struct {
	Name   string
	Weight uint64
}
type CargoCreateInput struct {
	ID       string
	VesselID string
	Items    []struct {
		Name   string
		Weight uint64
	}
	At time.Time
}
type CargoCreator struct {
	repository  CargoRepository
	idProvider  utils.ULIDProvider
	vesselCheck CargoVesselChecker
}

func NewCargoCreator(repository CargoRepository, checker CargoVesselChecker, idProvider utils.ULIDProvider) *CargoCreator {
	return &CargoCreator{
		repository:  repository,
		vesselCheck: checker,
		idProvider:  idProvider,
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

	cargoItems, err := newItemsFromRaw(input.Items)
	if err != nil {
		return nil, fmt.Errorf("error creating cargo items: %w", err)
	}

	trackingID, err := cargotrackingdomain.NewTrackingID(cc.idProvider.New().String())
	if err != nil {
		return nil, fmt.Errorf("error creating tracking id: %w", err)
	}

	cargo := NewCargo(id, vesselID, trackingID, cargoItems, input.At)

	if saveErr := cc.repository.Save(ctx, cargo); saveErr != nil {
		return nil, fmt.Errorf("error saving cargo: %w", saveErr)
	}

	return cargo, nil
}
