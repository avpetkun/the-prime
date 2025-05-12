package main

import (
	"context"
	"os"

	"github.com/avpetkun/the-prime/internal/support"
	"github.com/avpetkun/the-prime/pkg/loggeru"
	"github.com/avpetkun/the-prime/pkg/signalu"
)

func main() {
	ctx, ctxCancel := signalu.WaitExitContext(context.Background())
	defer ctxCancel()

	ctx, log := loggeru.GetLogger(ctx)

	cfg, err := support.LoadConfig(os.Getenv("CONFIG"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	service, err := support.NewService(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create service")
	}
	if err = service.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to start service")
	}

	log.Info().Msg("service started")
	<-ctx.Done()
	log.Info().Msg("service stopped")
}
