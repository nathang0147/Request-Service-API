package router

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	platformmiddleware "github.com/nathang0147/Request-Service-API/internal/platform/middleware"
)

type RouteRegistrar interface {
	RegisterRoutes(chi.Router)
}

func New(logger *zap.Logger, registrars ...RouteRegistrar) http.Handler {
	router := chi.NewRouter()
	router.Use(platformmiddleware.RequestID)
	router.Use(platformmiddleware.Logging(logger))
	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok"))
	})

	for _, registrar := range registrars {
		if registrar != nil {
			registrar.RegisterRoutes(router)
		}
	}

	return router
}
