package cargoqueries

import (
	"context"
	"fmt"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
)

type FetchCargoByID struct {
	ID string
}

func (q *FetchCargoByID) Type() string {
	return "fetch_cargo_by_id"
}

type FetchCargoByIDHandler struct {
	repository cargodomain.CargoRepository
}

func NewFetchCargoByIDHandler(repository cargodomain.CargoRepository) *FetchCargoByIDHandler {
	return &FetchCargoByIDHandler{
		repository: repository,
	}
}

func (h *FetchCargoByIDHandler) Handle(ctx context.Context, q *FetchCargoByID) (CargoResponse, error) {
	cargoID, err := cargodomain.NewCargoID(q.ID)
	if err != nil {
		return CargoResponse{}, fmt.Errorf("invalid cargo id: %w", err)
	}

	cargo, err := h.repository.Find(ctx, cargoID)
	if err != nil {
		return CargoResponse{}, fmt.Errorf("error fetching cargo: %w", err)
	}

	return NewCargoResponse(cargo.Primitives()), nil
}
