package handler

import (
	"github.com/go-chi/chi/v5"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

type Handler struct {
	Router  *chi.Mux
	Storage storage.Storage
}

func NewHandler() *Handler {
	return &Handler{
		Router:  chi.NewRouter(),
		Storage: storage.NewStorage(),
	}
}

func (h *Handler) SetRoute() {
	h.Router.Use(logger.RequestLogger)
	h.Router.Use(compressMiddleware)
	h.Router.Get("/", h.home)
	h.Router.Route("/value", func(r chi.Router) {
		r.Post("/", h.value)
		r.Get("/gauge/{name}", h.valueGauge)
		r.Get("/counter/{name}", h.valueCounter)
	})
	h.Router.Route("/update", func(r chi.Router) {
		r.Post("/", h.updateJSON)
		r.Post("/{type}/{name}/{value}", h.update)
	})
}
