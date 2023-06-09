package main

import (
	"github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
	"github.com/go-chi/chi/v5"
	"log"
)

const (
	PORT = "8080"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", handler.Main)
	r.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{name}", handler.ValueGauge)
		r.Get("/counter/{name}", handler.ValueCounter)
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", handler.Update)
	})

	srv := server.Server{}
	if err := srv.Run(PORT, r); err != nil {
		log.Fatalf("failed occured server: %s", err.Error())
	}
}
