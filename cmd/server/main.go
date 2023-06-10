package main

import (
	"github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
	"log"
	"os"
)

func main() {
	parseFlags()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	srv := server.Server{}
	if err := srv.Run(flagRunAddr, handler.MetricsRouter()); err != nil {
		log.Fatalf("failed occured server: %s", err.Error())
	}
}
