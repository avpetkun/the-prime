package main

import (
	"context"
	"os"

	"github.com/avpetkun/the-prime/internal/webhook"
	"github.com/avpetkun/the-prime/pkg/loggeru"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/signalu"
)

func main() {
	ctx, ctxCancel := signalu.WaitExitContext(context.Background())
	defer ctxCancel()

	ctx, log := loggeru.GetLogger(ctx)

	cfg, err := webhook.LoadConfig(os.Getenv("CONFIG"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	nats, err := natsu.Connect(log, cfg.NatsAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to nats queue")
	}
	defer nats.Stop()

	service := webhook.NewService(cfg, log, nats)
	if err = service.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to create webhook service")
	}
	defer service.Close()

	log.Info().Msg("service started")
	<-ctx.Done()
	log.Info().Msg("service stopped")
}
