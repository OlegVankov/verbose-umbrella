package main

import (
	hdl "github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
	"go.uber.org/zap"
)

func main() {
	parseFlags()

	handler := hdl.NewHandler()
	handler.SetRoute()

	srv := server.Server{}

	if err := logger.Initialize(level); err != nil {
		panic(err)
	}

	if err := srv.Run(serverAddr, handler.Router); err != nil {
		logger.Log.Fatal(err.Error(), zap.String("event", "start server"))
	}
}
