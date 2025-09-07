package cargoinfra

import (
	"context"

	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	vesselqueries "github.com/soulcodex/deus-cargo-tracker/internal/vessel/application/queries"
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
	querybus "github.com/soulcodex/deus-cargo-tracker/pkg/bus/query"
	"github.com/soulcodex/deus-cargo-tracker/pkg/errutil"
)

var (
	_ cargodomain.CargoVesselChecker = (*QueryBusVesselChecker)(nil)

	ErrVesselNotFound = errutil.NewError("vessel not found")
)

type QueryBusVesselChecker struct {
	queryBus querybus.Bus
}

func NewQueryBusVesselChecker(queryBus querybus.Bus) *QueryBusVesselChecker {
	return &QueryBusVesselChecker{queryBus: queryBus}
}

func (q *QueryBusVesselChecker) Check(ctx context.Context, vesselID cargodomain.VesselID) error {
	query := &vesselqueries.FetchVesselByIDQuery{ID: vesselID.String()}
	err := bus.Dispatch(q.queryBus)(ctx, query)
	if err != nil {
		return ErrVesselNotFound.Wrap(err)
	}

	return nil
}
