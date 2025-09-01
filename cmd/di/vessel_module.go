package di

import (
	"context"

	vesselqueries "github.com/soulcodex/deus-cargo-tracker/internal/vessel/application/queries"
	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
	vesselentrypoint "github.com/soulcodex/deus-cargo-tracker/internal/vessel/infrastructure/entrypoint"
	vesselpersistence "github.com/soulcodex/deus-cargo-tracker/internal/vessel/infrastructure/persistence"
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
)

type VesselModule struct {
	Repository vesseldomain.VesselRepository
}

func NewVesselModule(_ context.Context, common *CommonServices) *VesselModule {
	vesselRepo := vesselpersistence.NewPostgresVesselRepository(common.Config.PostgresSchema, common.DBPool)

	common.Router.Get(
		"/vessels/{vessel_id}",
		vesselentrypoint.HandleGETFetchVesselByIDV1HTTP(
			common.QueryBus,
			common.ResponseMiddleware,
		),
	)

	fetchVesselByIDHandler := vesselqueries.NewFetchVesselByIDHandler(vesselRepo)

	bus.MustRegister(common.QueryBus, &vesselqueries.FetchVesselByID{}, fetchVesselByIDHandler)

	return &VesselModule{
		Repository: vesselRepo,
	}
}
