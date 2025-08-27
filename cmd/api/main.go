package main

import (
	"context"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"

	"github.com/soulcodex/deus-cargoes-tracker/cmd/di"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	common := di.MustInitCommonServices(ctx)

	go func() {
		common.Logger.Info().
			Str("http.host", common.Config.HTTPHost).
			Int("http.port", common.Config.HTTPPort).
			Msg("starting http server")

		if listenErr := common.Router.ListenAndServe(); listenErr != nil {
			common.Logger.Fatal().Err(listenErr).Msg("error starting HTTP server")
		}
	}()

	common.Logger.Info().Msg("cargo tracker HTTP started successfully")
	<-ctx.Done()
}
