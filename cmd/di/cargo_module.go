package di

import (
	"context"

	cargocommands "github.com/soulcodex/deus-cargo-tracker/internal/cargo/application/commands"
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
	cargoesRepo := cargopersistence.NewPostgresCargoRepository(common.Config.PostgresSchema, common.DBPool)
	createCargoHTTPHandler := cargoentrypoint.HandlePOSTCreateCargoV1HTTP(common.CommandBus, common.ResponseMiddleware)
	cargoVesselChecker := cargoinfra.NewQueryBusVesselChecker(common.QueryBus)
	cargoCreator := cargodomain.NewCargoCreator(cargoesRepo, cargoVesselChecker)

	common.Router.Post("/cargoes", createCargoHTTPHandler)

	bus.MustRegister(
		common.CommandBus,
		&cargocommands.CreateCargoCommand{},
		cargocommands.NewCreateCargoCommandHandler(cargoCreator, common.TimeProvider),
	)

	return &CargoModule{
		Repository: cargoesRepo,
	}
}
