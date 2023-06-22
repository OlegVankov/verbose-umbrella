package server

import (
	"context"
	"log"
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
	log.Printf("server starting: %s", address)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
