package handler

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
)

func (h *Handler) ping(w http.ResponseWriter, req *http.Request) {
	logger.Log.Info("ping request")
	if err := h.Storage.PingStorage(req.Context()); err != nil {
		logger.Log.Info("ping result", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Log.Info("ping response ok")
}
