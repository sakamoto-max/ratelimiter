package routes

import (
	"github.com/sakamoto-max/ratelimiter/internal/handlers"
	"github.com/sakamoto-max/ratelimiter/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *handlers.Handler, reg *prometheus.Registry) *chi.Mux {
	r := chi.NewRouter()

	middlewares := middleware.New()

	r.Get("/health", handler.Health.Health)

	r.With(middlewares.Auth).Get("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP)

	r.Post("/register", handler.Owner.Register)

	r.With(middlewares.Auth).Post("/policy", handler.Policy.New)
	r.With(middlewares.Auth).Get("/policy", handler.Policy.Get)
	r.With(middlewares.Auth).Patch("/policy", handler.Policy.Patch)
	r.With(middlewares.Auth).Delete("/policy", handler.Policy.Delete)

	r.With(middlewares.Auth).Get("/token", handler.Token.GetToken)
	r.With(middlewares.Auth).Post("/token", handler.Token.NewToken)
	r.With(middlewares.Auth).Delete("/token", handler.Token.DeleteToken)

	return r
}
