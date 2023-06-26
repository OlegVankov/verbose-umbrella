package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"strconv"

	"github.com/OlegVankov/verbose-umbrella/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Router  *chi.Mux
	storage storage.Storage
}

func NewHandler() Handler {
	return Handler{
		Router:  chi.NewRouter(),
		storage: storage.NewStorage(),
	}
}

func (h *Handler) SetRoute() {
	h.Router.Get("/", logger.RequestLogger(h.home))
	h.Router.Route("/value", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(h.value))
		r.Get("/gauge/{name}", logger.RequestLogger(h.valueGauge))
		r.Get("/counter/{name}", logger.RequestLogger(h.valueCounter))
	})
	h.Router.Route("/update", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(h.update2))
		r.Post("/{type}/{name}/{value}", logger.RequestLogger(h.update))
	})
}

func (h *Handler) value(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var metric storage.Metrics
	_, err := buf.ReadFrom(r.Body)
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
		delta, _ := h.storage.GetCounter(metric.ID)
		metric.Delta = (*int64)(&delta)
	case "gauge":
		value, _ := h.storage.GetGauge(metric.ID)
		metric.Value = (*float64)(&value)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	resp, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	logger.Log.Info("/value request:", zap.ByteString("body", buf.Bytes()))
	logger.Log.Info("/value response:", zap.ByteString("body", resp))
	_, _ = w.Write(resp)
}

func (h *Handler) update2(w http.ResponseWriter, req *http.Request) {
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
		h.storage.UpdateCounter(metric.ID, *metric.Delta)
		delta, _ := h.storage.GetCounter(metric.ID)
		*metric.Delta = storage.CounterToInt(delta)
	case "gauge":
		h.storage.UpdateGauge(metric.ID, *metric.Value)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	resp, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	logger.Log.Info("/update request:", zap.ByteString("body", buf.Bytes()))
	logger.Log.Info("/update response:", zap.ByteString("body", resp))
	_, _ = w.Write(resp)
}

func (h *Handler) home(w http.ResponseWriter, req *http.Request) {
	ts, err := template.ParseFiles("./html/home.page.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	err = ts.Execute(w, h.storage)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
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
		h.storage.UpdateCounter(name, value)
	case "gauge":
		value, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateGauge(name, value)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) valueGauge(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	value, ok := h.storage.GetGauge(name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)

	val := float64(value)
	_, _ = w.Write([]byte(fmt.Sprint(val)))
}

func (h *Handler) valueCounter(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	value, ok := h.storage.GetCounter(name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)

	val := int64(value)
	_, _ = w.Write([]byte(fmt.Sprint(val)))
}
