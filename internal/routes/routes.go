package routes

import (
	"github.com/sakamoto-max/ratelimiter/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *handlers.Handler, reg *prometheus.Registry) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", handler.Health.Health)

	r.Get("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP)

	r.Post("/register", handler.Owner.Register)
	r.Post("/policy", handler.Policy.New)
	r.Get("/policy",handler.Policy.Get)
	r.Patch("/policy",handler.Policy.Patch)
	r.Delete("/policy",handler.Policy.Delete)

	return r
}
