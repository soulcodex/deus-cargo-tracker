package di

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/soulcodex/deus-cargo-tracker/configs"
	commandbus "github.com/soulcodex/deus-cargo-tracker/pkg/bus/command"
	querybus "github.com/soulcodex/deus-cargo-tracker/pkg/bus/query"
	distributedsync "github.com/soulcodex/deus-cargo-tracker/pkg/distributed-sync"
	httpserver "github.com/soulcodex/deus-cargo-tracker/pkg/http-server"
	"github.com/soulcodex/deus-cargo-tracker/pkg/logger"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

type CommonServices struct {
	Config             *configs.Config
	DBPool             sqldb.ConnectionPool
	DBMigrator         sqldb.Migrator
	ResponseMiddleware *httpserver.JSONAPIResponseMiddleware
	Logger             logger.ZerologLogger
	RedisClient        *redis.Client
	CommandBus         commandbus.Bus
	QueryBus           querybus.Bus
	Mutex              distributedsync.MutexService
	Router             *httpserver.Router
	UUIDProvider       utils.UUIDProvider
	ULIDProvider       utils.ULIDProvider
	TimeProvider       utils.DateTimeProvider
}

func MustInitCommonServices(ctx context.Context) *CommonServices {
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	appLogger := logger.NewZerologLogger(
		ctx,
		"cargo-tracker",
		logger.WithLogLevel(cfg.LogLevel),
		logger.WithAppVersion(cfg.AppVersion),
	)

	timeProvider := utils.NewSystemTimeProvider()

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(redisOpts)

	router := initHTTPRouter(ctx, appLogger, timeProvider, cfg)
	dbPool := initPostgresDBPool(ctx, cfg)
	dbMigrator := initSQLMigrator(ctx, cfg, dbPool)

	queryBus := querybus.InitQueryBus()
	commandBus := commandbus.InitCommandBus()
	mutexService := distributedsync.NewRedisMutexService(redisClient, appLogger)
	uuidProvider := utils.NewRandomUUIDProvider()
	ulidProvider := utils.NewRandomULIDProvider()
	responseMiddleware := httpserver.NewJSONAPIResponseMiddleware(appLogger)

	return &CommonServices{
		Config:             cfg,
		Logger:             appLogger,
		DBPool:             dbPool,
		DBMigrator:         dbMigrator,
		ResponseMiddleware: responseMiddleware,
		RedisClient:        redisClient,
		QueryBus:           queryBus,
		CommandBus:         commandBus,
		Mutex:              mutexService,
		Router:             router,
		UUIDProvider:       uuidProvider,
		ULIDProvider:       ulidProvider,
		TimeProvider:       timeProvider,
	}
}

func MustInitCommonServicesWithEnvFiles(ctx context.Context, envFiles ...string) *CommonServices {
	err := godotenv.Overload(envFiles...)
	if err != nil {
		panic(err)
	}

	return MustInitCommonServices(ctx)
}
