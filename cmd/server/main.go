package main

import (
	"log"
	"os"

	hdl "github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
)

func main() {
	parseFlags()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverAddr = envRunAddr
	}

	handler := hdl.NewHandler()
	handler.SetRoute()

	srv := server.Server{}
	if err := srv.Run(serverAddr, handler.Router); err != nil {
		log.Fatalf("failed occured server: %s", err.Error())
	}
}
