package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
)

func NewRouter(handler *Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(logger.RequestLogger)
	router.Use(compressMiddleware)
	router.Use(handler.checkHash)
	router.Get("/", handler.home)
	router.Get("/ping", handler.ping)
	router.Route("/value", func(r chi.Router) {
		r.Post("/", handler.value)
		r.Get("/gauge/{name}", handler.valueGauge)
		r.Get("/counter/{name}", handler.valueCounter)
	})
	router.Route("/update", func(r chi.Router) {
		r.Post("/", handler.updateJSON)
		r.Post("/{type}/{name}/{value}", handler.update)
	})
	router.Route("/updates", func(r chi.Router) {
		r.Post("/", handler.updates)
	})

	return router
}
