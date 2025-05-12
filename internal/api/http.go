package api

import (
	"context"
	"fmt"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"

	"github.com/avpetkun/the-prime/pkg/metricu"
)

type HTTP struct {
	cfg HTTPConfig
	log zerolog.Logger
	svc *Service

	public  *fiber.App
	private *fiber.App
}

func NewHTTP(cfg HTTPConfig, log zerolog.Logger, svc *Service, metrics *prometheus.Registry) *HTTP {
	api := &HTTP{
		cfg: cfg,
		log: log,
		svc: svc,
	}

	fiberConfig := fiber.Config{DisableStartupMessage: true}

	public := fiber.New(fiberConfig)
	public.Use(cors.New())
	{
		v1 := public.Group("/api/v1")

		v1.Use(api.userAuthMiddleware)
		v1.Post("/init", api.handlerUserInit)
		v1.Get("/overview", api.handlerGetOverview())
		v1.Get("/tasks/events", api.handlerGetTasksEvents)
		v1.Post("/tasks/:task_id/:sub_id/start", api.handlerTaskStart)
		v1.Post("/tasks/:task_id/:sub_id/claim", api.handlerTaskClaim)
		v1.Post("/tasks/:task_id/invoice/stars", api.handlerTaskStarsInvoice)
		v1.Post("/tasks/:task_id/invoice/ton", api.handlerTaskTonInvoice)
		v1.Post("/products/:product_id/claim", api.handlerProductClaim)
		v1.Post("/invite-message", api.handlerGetInviteMessage)

		v1.Use(api.adminRequiredMiddleware)
		v1.Get("/overview-admin", api.handlerAdminGetOverview())
		v1.Post("/products", api.handlerAdminSaveProduct)
		v1.Delete("/products/:product_id", api.handlerAdminDeleteProduct)
		v1.Post("/tasks", api.handlerAdminSaveTask)
		v1.Delete("/tasks/:task_id", api.handlerAdminDeleteTask)
		v1.Post("/reward-user", api.handlerAdminRewardUser)
	}
	api.public = public

	private := fiber.New(fiberConfig)
	private.Use(cors.New())
	private.Use(pprof.New())
	{
		metricu.RegisterFiberPrometheus(private, "/metrics", metrics)
	}
	api.private = private

	return api
}

func (h *HTTP) Start(ctx context.Context) error {
	h.log.Info().Msgf("[api] start private api on port %d", h.cfg.PrivatePort)

	privateSocket, err := net.Listen("tcp", fmt.Sprintf(":%d", h.cfg.PrivatePort))
	if err != nil {
		return err
	}

	h.log.Info().Msgf("[api] start public api on port %d", h.cfg.PublicPort)

	publicSocket, err := net.Listen("tcp", fmt.Sprintf(":%d", h.cfg.PublicPort))
	if err != nil {
		return err
	}

	go h.private.Listener(privateSocket)
	go h.public.Listener(publicSocket)
	go func() {
		<-ctx.Done()
		publicSocket.Close()
		privateSocket.Close()
	}()
	return nil
}
