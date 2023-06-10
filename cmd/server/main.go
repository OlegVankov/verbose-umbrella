package main

import (
	"github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
	"log"
)

const (
	PORT = "8080"
)

func main() {
	srv := server.Server{}
	if err := srv.Run(PORT, handler.MetricsRouter()); err != nil {
		log.Fatalf("failed occured server: %s", err.Error())
	}
}
