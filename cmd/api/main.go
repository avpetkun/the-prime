package main

import (
	"context"
	"os"

	"github.com/avpetkun/the-prime/internal/api"
	"github.com/avpetkun/the-prime/internal/cache"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/pkg/loggeru"
	"github.com/avpetkun/the-prime/pkg/metricu"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/signalu"
	"github.com/avpetkun/the-prime/pkg/tgu"
)

func main() {
	ctx, ctxCancel := signalu.WaitExitContext(context.Background())
	defer ctxCancel()

	ctx, log := loggeru.GetLogger(ctx)

	cfg, err := api.LoadConfig(os.Getenv("CONFIG"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	metrics := metricu.NewPrometheusRegistry()

	cache, err := cache.Connect(ctx, cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to cache")
	}
	defer cache.Close()

	db, err := dbx.Connect(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	defer db.Close()

	nats, err := natsu.Connect(log, cfg.Nats.Addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to nats queue")
	}
	defer nats.Stop()

	bot, err := tgu.CreateBot(log, cfg.Bot)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create telegram bot")
	}

	service := api.NewService(cfg, log, db, cache, nats, bot)

	if err = service.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to start api service")
	}

	httpAPI := api.NewHTTP(cfg.HTTP, log, service, metrics)
	if err = httpAPI.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to start http api")
	}

	log.Info().Msg("service started")
	<-ctx.Done()
	log.Info().Msg("service stopped")
}
