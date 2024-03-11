package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/handler"
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

func Run(address string, storage storage.Storage, key string) error {
	logger.Log.Info("server", zap.String("starting", "..."))

	hdl := handler.NewHandler(storage, key)
	server := &http.Server{
		Addr:           address,
		Handler:        handler.NewRouter(hdl),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		logger.Log.Info("server", zap.String("listen address", address))
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("HTTP server ListenAndServe", zap.Error(err))
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c

	logger.Log.Info("server", zap.String("Graceful shutdown starter with signal", sig.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
