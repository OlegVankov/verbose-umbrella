package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

func main() {
	parseFlags()

	if err := logger.Initialize(level); err != nil {
		panic(err)
	}

	conn, err := storage.ConnectDB(databaseDSN)
	if err != nil {
		logger.Log.Fatal("Connection DB", zap.Error(err))
	}

	newHandler := handler.NewHandler(conn)
	newHandler.SetRoute()
	srv := server.Server{}

	if restore {
		if err := newHandler.Storage.RestoreStorage(fileStoragePath); err != nil {
			logger.Log.Fatal("Restore storage", zap.Error(err))
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
