package vesselqueries

import (
	"context"
	"fmt"

	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
)

type FetchVesselByID struct {
	ID string
}

func (q *FetchVesselByID) Type() string {
	return "fetch_vessel_by_id"
}

type FetchVesselByIDHandler struct {
	repository vesseldomain.VesselRepository
}

func NewFetchVesselByIDHandler(repository vesseldomain.VesselRepository) *FetchVesselByIDHandler {
	return &FetchVesselByIDHandler{
		repository: repository,
	}
}

func (h *FetchVesselByIDHandler) Handle(ctx context.Context, q *FetchVesselByID) (VesselResponse, error) {
	vesselID, err := vesseldomain.NewVesselID(q.ID)
	if err != nil {
		return VesselResponse{}, fmt.Errorf("invalid vessel id: %w", err)
	}

	vessel, err := h.repository.Find(ctx, vesselID)
	if err != nil {
		return VesselResponse{}, fmt.Errorf("error fetching vessel: %w", err)
	}

	return NewVesselResponse(vessel.Primitives()), nil
}
