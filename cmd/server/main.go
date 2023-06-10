package main

import (
	"github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
	"log"
)

func main() {
	parseFlags()
	srv := server.Server{}
	if err := srv.Run(flagRunAddr, handler.MetricsRouter()); err != nil {
		log.Fatalf("failed occured server: %s", err.Error())
	}
}
