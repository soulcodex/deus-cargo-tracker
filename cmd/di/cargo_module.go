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
	updateCargoStatusHTTPHandler := cargoentrypoint.HandlePATCHUpdateCargoStatusV1HTTP(
		common.CommandBus,
		common.Mutex,
		common.ResponseMiddleware,
	)
	cargoVesselChecker := cargoinfra.NewQueryBusVesselChecker(common.QueryBus)
	cargoCreator := cargodomain.NewCargoCreator(cargoRepo, cargoVesselChecker, common.ULIDProvider)
	cargoUpdater := cargodomain.NewCargoUpdater(cargoRepo)

	common.Router.Post("/cargoes", createCargoHTTPHandler)
	common.Router.Get("/cargoes/{cargo_id}", fetchCargoByIDHTTPHandler)
	common.Router.Patch("/cargoes/{cargo_id}/update-status", updateCargoStatusHTTPHandler)

	bus.MustRegister(
		common.CommandBus,
		&cargocommands.CreateCargoCommand{},
		cargocommands.NewCreateCargoCommandHandler(cargoCreator, common.TimeProvider),
	)

	bus.MustRegister(
		common.CommandBus,
		&cargocommands.UpdateCargoStatusCommand{},
		cargocommands.NewUpdateCargoStatusCommandHandler(cargoUpdater, common.TimeProvider, common.ULIDProvider),
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
