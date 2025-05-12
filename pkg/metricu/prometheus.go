package metricu

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheusRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()

	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	return reg
}

func RegisterFiberPrometheus(app *fiber.App, path string, reg *prometheus.Registry) {
	app.Get(path, adaptor.HTTPHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
}
