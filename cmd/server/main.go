package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
	"go.uber.org/zap"
)

func main() {
	parseFlags()

	newHandler := handler.NewHandler()
	newHandler.SetRoute()
	srv := server.Server{}

	if err := logger.Initialize(level); err != nil {
		panic(err)
	}
	if restore {
		if err := newHandler.Storage.RestoreStorage(fileStoragePath); err != nil {
			panic(err)
		}
	}
	go newHandler.Storage.SaveStorage(fileStoragePath, storeInterval)

	go func() {
		err := srv.Run(serverAddr, newHandler.Router)
		if err != nil && err == http.ErrServerClosed {
			logger.Log.Fatal(err.Error(), zap.String("event", "start server"))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
