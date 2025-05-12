package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/timeu"
)

type Service struct {
	cfg Config
	log zerolog.Logger
	ns  *natsu.Stream

	httpSocket net.Listener
}

func NewService(cfg Config, log zerolog.Logger, ns *natsu.Stream) *Service {
	return &Service{cfg: cfg, log: log, ns: ns}
}

func (s *Service) Close() error {
	return s.httpSocket.Close()
}

func (s *Service) Start(ctx context.Context) error {
	wrappedPublish := func(subject string, message any) error {
		data, err := json.Marshal(message)
		if err != nil {
			return err
		}
		for {
			err = s.ns.PublishData(ctx, subject, data)
			if err == nil {
				s.log.Info().
					Bytes("message", data).
					Str("subject", subject).
					Msg("partner event sent")
				return nil
			}
			s.log.Error().Err(err).
				Bytes("message", data).
				Str("subject", subject).
				Msg("failed to publish message")
			if timeu.SleepContext(ctx, time.Second) {
				return nil
			}
		}
	}

	fap := fiber.New(fiber.Config{DisableStartupMessage: true})
	fap.Use(cors.New())

	v1 := fap.Group("/api/v1")

	v1.Get("/reward/:token", handlerReward(wrappedPublish))
	v1.Get("/tappads/:token", handlerRewardTappAds(wrappedPublish))

	s.log.Info().Msgf("start public api on port %d", s.cfg.HttpPort)
	var err error
	s.httpSocket, err = net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.HttpPort))
	if err != nil {
		return err
	}

	go fap.Listener(s.httpSocket)

	return nil
}
