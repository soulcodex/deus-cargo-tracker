package di

import (
	"context"

	"github.com/soulcodex/deus-cargo-tracker/configs"
	httpserver "github.com/soulcodex/deus-cargo-tracker/pkg/http-server"
	"github.com/soulcodex/deus-cargo-tracker/pkg/logger"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

func initHTTPRouter(
	_ context.Context,
	appLogger logger.ZerologLogger,
	timeProvider utils.DateTimeProvider,
	cfg *configs.Config,
) *httpserver.Router {
	routerOpts := []httpserver.RouterConfigFunc{
		httpserver.WithHost(cfg.HTTPHost),
		httpserver.WithPort(cfg.HTTPPort),
		httpserver.WithReadTimeoutSeconds(cfg.HTTPReadTimeout),
		httpserver.WithWriteTimeoutSeconds(cfg.HTTPWriteTimeout),
		httpserver.WithMiddleware(httpserver.NewPanicRecoverMiddleware(appLogger).Middleware),
		httpserver.WithMiddleware(httpserver.NewRequestLoggingMiddleware(appLogger, timeProvider).Middleware),
		httpserver.WithCORSMiddleware(),
	}

	router := httpserver.New(routerOpts...)

	return &router
}
