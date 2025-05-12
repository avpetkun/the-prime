package main

import (
	"context"
	"os"

	"github.com/avpetkun/the-prime/internal/cache"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/internal/worker"
	"github.com/avpetkun/the-prime/pkg/fragment"
	"github.com/avpetkun/the-prime/pkg/loggeru"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/signalu"
	"github.com/avpetkun/the-prime/pkg/tgu"
	"github.com/avpetkun/the-prime/pkg/tonu"
)

func main() {
	ctx, ctxCancel := signalu.WaitExitContext(context.Background())
	defer ctxCancel()

	ctx, log := loggeru.GetLogger(ctx)

	cfg, err := worker.LoadConfig(os.Getenv("CONFIG"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

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

	tonApi, err := tonu.ConnectGlobal(ctx, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to ton")
	}
	fragmentApi, err := fragment.NewAPI(tonApi, log, cfg.Fragment.AuthCookie, cfg.Fragment.WalletSeed)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create fragment api")
	}

	service := worker.NewService(cfg, log, cache, db, nats, bot, tonApi, fragmentApi)

	bot.LoadOffset = func(ctx context.Context) (offset int, err error) {
		return db.GetKeyNum(ctx, "tg_bot:last_offset")
	}
	bot.SaveOffset = func(ctx context.Context, offset int) error {
		return db.SetKeyNum(ctx, "tg_bot:last_offset", offset)
	}
	bot.OnStarsTx = service.OnTelegramStarsTx
	bot.OnPrivateText = service.OnTelegramPrivateMessage
	bot.OnSelfMember = service.OnTelegramSelfMember

	if err = service.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to start worker")
	}

	if err = bot.Listen(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to listen tg bot")
	}

	log.Info().Msg("service started")
	<-ctx.Done()
	log.Info().Msg("service stopped")
}
