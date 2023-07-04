package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (h *Handler) value(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var metric storage.Metrics

	if err := decoder.Decode(&metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Info("Error", zap.Error(err))
		return
	}

	switch metric.MType {
	case "counter":
		val, _ := h.Storage.GetCounter(metric.ID)
		delta := int64(val)
		metric.Delta = &delta
	case "gauge":
		val, _ := h.Storage.GetGauge(metric.ID)
		value := float64(val)
		metric.Value = &value
	}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	logger.Log.Warn("value", zap.Any("metric", metric))
	if err := json.NewEncoder(w).Encode(metric); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) valueGauge(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	value, ok := h.Storage.GetGauge(name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	//w.WriteHeader(http.StatusOK)

	val := float64(value)
	w.Write([]byte(fmt.Sprint(val)))
}

func (h *Handler) valueCounter(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	value, ok := h.Storage.GetCounter(name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	//w.WriteHeader(http.StatusOK)

	val := int64(value)
	w.Write([]byte(fmt.Sprint(val)))
}
