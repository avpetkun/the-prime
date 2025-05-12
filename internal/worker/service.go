package worker

import (
	"context"
	"sync/atomic"

	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/internal/cache"
	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/pkg/fragment"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/tgu"
	"github.com/avpetkun/the-prime/pkg/tonu"
)

type Service struct {
	cfg  Config
	log  zerolog.Logger
	ton  *tonu.API
	frag *fragment.API
	bot  *tgu.Bot
	ch   cache.Cache
	db   *dbx.DB
	ns   *natsu.Stream

	allTasks []*common.FullTask

	wasInsufficientBalance atomic.Int32
}

func NewService(
	cfg Config, log zerolog.Logger,
	ch cache.Cache, db *dbx.DB, ns *natsu.Stream,
	bot *tgu.Bot, tonApi *tonu.API, fragmentApi *fragment.API,
) *Service {
	return &Service{
		cfg:  cfg,
		log:  log,
		ton:  tonApi,
		frag: fragmentApi,
		bot:  bot,
		ch:   ch,
		db:   db,
		ns:   ns,
	}
}

func (s *Service) Start(ctx context.Context) error {
	if err := s.startLoadTasksLoop(ctx); err != nil {
		return err
	}

	if err := s.listenTonTransactions(ctx); err != nil {
		return err
	}
	if err := s.startWebhookWorkers(ctx); err != nil {
		return err
	}
	if err := s.startUsersWorkers(ctx); err != nil {
		return err
	}
	if err := s.startTasksWorkers(ctx); err != nil {
		return err
	}
	if err := s.startTelegramWorkers(ctx); err != nil {
		return err
	}
	if err := s.startProductsWorkers(ctx); err != nil {
		return err
	}
	if err := s.startChecksWorkers(ctx); err != nil {
		return err
	}
	if err := s.startFragmentWorkers(ctx); err != nil {
		return err
	}

	return nil
}
