package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

func (h *Handler) updateJSON(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	var metric storage.Metrics
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "counter":
		delta := storage.CounterToInt(h.Storage.UpdateCounter(metric.ID, *metric.Delta))
		metric.Delta = &delta
	case "gauge":
		h.Storage.UpdateGauge(metric.ID, *metric.Value)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	logger.Log.Warn("/update", zap.ByteString("request", buf.Bytes()), zap.ByteString("response", resp))
	json.NewEncoder(w).Encode(metric)
}

func (h *Handler) update(w http.ResponseWriter, req *http.Request) {
	typeMetric := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")

	switch typeMetric {
	case "counter":
		value, err := strconv.ParseInt(chi.URLParam(req, "value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.Storage.UpdateCounter(name, value)
	case "gauge":
		value, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.Storage.UpdateGauge(name, value)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
}
