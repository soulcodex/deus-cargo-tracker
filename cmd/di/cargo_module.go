package di

import (
	"context"

	cargocommands "github.com/soulcodex/deus-cargo-tracker/internal/cargo/application/commands"
	cargoqueries "github.com/soulcodex/deus-cargo-tracker/internal/cargo/application/queries"
	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargoinfra "github.com/soulcodex/deus-cargo-tracker/internal/cargo/infrastructure"
	cargoentrypoint "github.com/soulcodex/deus-cargo-tracker/internal/cargo/infrastructure/entrypoint"
	cargopersistence "github.com/soulcodex/deus-cargo-tracker/internal/cargo/infrastructure/persistence"
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
)

type CargoModule struct {
	Repository cargodomain.CargoRepository
}

func NewCargoModule(_ context.Context, common *CommonServices) *CargoModule {
	cargoRepo := cargopersistence.NewPostgresCargoRepository(common.Config.PostgresSchema, common.DBPool)
	createCargoHTTPHandler := cargoentrypoint.HandlePOSTCreateCargoV1HTTP(common.CommandBus, common.ResponseMiddleware)
	fetchCargoByIDHTTPHandler := cargoentrypoint.HandleGETFetchCargoByIDV1HTTP(common.QueryBus, common.ResponseMiddleware)
	cargoVesselChecker := cargoinfra.NewQueryBusVesselChecker(common.QueryBus)
	cargoCreator := cargodomain.NewCargoCreator(cargoRepo, cargoVesselChecker)

	common.Router.Post("/cargoes", createCargoHTTPHandler)
	common.Router.Get("/cargoes/{cargo_id}", fetchCargoByIDHTTPHandler)

	bus.MustRegister(
		common.CommandBus,
		&cargocommands.CreateCargoCommand{},
		cargocommands.NewCreateCargoCommandHandler(cargoCreator, common.TimeProvider),
	)

	bus.MustRegister(
		common.QueryBus,
		&cargoqueries.FetchCargoByID{},
		cargoqueries.NewFetchCargoByIDHandler(cargoRepo),
	)

	return &CargoModule{
		Repository: cargoRepo,
	}
}
