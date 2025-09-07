package vesselqueries

import (
	"context"
	"fmt"

	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
)

type FetchVesselByIDQuery struct {
	ID string
}

func (q *FetchVesselByIDQuery) Type() string {
	return "fetch_vessel_by_id_query"
}

type FetchVesselByIDQueryHandler struct {
	repository vesseldomain.VesselRepository
}

func NewFetchVesselByIDQueryHandler(repository vesseldomain.VesselRepository) *FetchVesselByIDQueryHandler {
	return &FetchVesselByIDQueryHandler{
		repository: repository,
	}
}

func (h *FetchVesselByIDQueryHandler) Handle(ctx context.Context, q *FetchVesselByIDQuery) (VesselResponse, error) {
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
