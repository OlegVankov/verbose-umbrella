package handler

import (
	"bufio"
	"encoding/json"
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"os"
	"time"
)

type Handler struct {
	Router  *chi.Mux
	storage storage.Storage
}

func NewHandler() *Handler {
	return &Handler{
		Router:  chi.NewRouter(),
		storage: storage.NewStorage(),
	}
}

func (h *Handler) SetRoute() {
	h.Router.Use(logger.RequestLogger)
	h.Router.Use(compressMiddleware)
	h.Router.Get("/", h.home)
	h.Router.Route("/value", func(r chi.Router) {
		r.Post("/", h.value)
		r.Get("/gauge/{name}", h.valueGauge)
		r.Get("/counter/{name}", h.valueCounter)
	})
	h.Router.Route("/update", func(r chi.Router) {
		r.Post("/", h.updateJSON)
		r.Post("/{type}/{name}/{value}", h.update)
	})
}

func (h *Handler) SaveStorage(fileStoragePath string, storeInterval int) {
	for {
		<-time.After(time.Duration(storeInterval) * time.Second)

		file, _ := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

		for k, v := range h.storage.GetGaugeAll() {
			value := storage.GaugeToFloat(v)
			metric := storage.Metrics{
				ID:    k,
				MType: "gauge",
				Value: &value,
			}
			data, _ := json.Marshal(&metric)
			_, _ = file.Write(data)
			_, _ = file.Write([]byte("\n"))
		}

		for k, v := range h.storage.GetCounterAll() {
			value := storage.CounterToInt(v)
			metric := storage.Metrics{
				ID:    k,
				MType: "counter",
				Delta: &value,
			}
			data, _ := json.Marshal(&metric)
			_, _ = file.Write(data)
			_, _ = file.Write([]byte("\n"))
		}

		_ = file.Close()

		logger.Log.Info("save storage", zap.String("file", fileStoragePath))
	}
}

func (h *Handler) RestoreStorage(fileStoragePath string) error {
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var metric storage.Metrics
		_ = json.Unmarshal(scanner.Bytes(), &metric)
		switch metric.MType {
		case "counter":
			h.storage.UpdateCounter(metric.ID, *metric.Delta)
		case "gauge":
			h.storage.UpdateGauge(metric.ID, *metric.Value)
		}
	}

	logger.Log.Info("restore storage", zap.String("file", fileStoragePath))
	return nil
}
