package server

import (
	"context"
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func (s *Server) Run(address string, handler http.Handler) error {
	s.srv = &http.Server{
		Addr:           address,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	logger.Log.Info("Running server",
		zap.String("address", s.srv.Addr),
	)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
