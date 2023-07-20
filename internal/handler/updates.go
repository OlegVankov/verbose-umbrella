package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

func (h *Handler) updates(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var metrics []storage.Metrics
	if err := json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := []string{}
	for _, m := range metrics {
		switch m.MType {
		case "counter":
			delta, err := h.Storage.UpdateCounter(req.Context(), m.ID, *m.Delta)
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
			m.Delta = &delta
		case "gauge":
			h.Storage.UpdateGauge(req.Context(), m.ID, *m.Value)
		}
	}

	if len(errs) != 0 {
		logger.Log.Info(strings.Join(errs, "|"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}
