package memory

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

type Storage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewStorage() *Storage {
	return &Storage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
}

func (m *Storage) UpdateGauge(ctx context.Context, name string, val float64) {
	m.Gauge[name] = val
}

func (m *Storage) UpdateCounter(ctx context.Context, name string, val int64) (int64, error) {
	m.Counter[name] += val
	return m.Counter[name], nil
}

func (m *Storage) GetGauge(ctx context.Context, name string) (float64, bool) {
	val, ok := m.Gauge[name]
	return val, ok
}

func (m *Storage) GetCounter(ctx context.Context, name string) (int64, bool) {
	val, ok := m.Counter[name]
	return val, ok
}

func (m *Storage) GetGaugeAll(ctx context.Context) map[string]float64 {
	return m.Gauge
}

func (m *Storage) GetCounterAll(ctx context.Context) map[string]int64 {
	return m.Counter
}

func (m *Storage) RestoreStorage(fileStoragePath string) error {
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var metric storage.Metrics
		err := json.Unmarshal(scanner.Bytes(), &metric)
		if err != nil {
			logger.Log.Warn("restore storage: json unmarshal",
				zap.Error(err))
			continue
		}
		switch metric.MType {
		case "counter":
			m.UpdateCounter(context.TODO(), metric.ID, *metric.Delta)
		case "gauge":
			m.UpdateGauge(context.TODO(), metric.ID, *metric.Value)
		}
	}

	logger.Log.Info("restore storage", zap.String("file", fileStoragePath))
	return nil
}

func (m *Storage) SaveStorage(fileStoragePath string, storeInterval int) {
	for {
		<-time.After(time.Duration(storeInterval) * time.Second)

		file, err := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			logger.Log.Error("error open file", zap.String("file name", fileStoragePath), zap.Error(err))
			return
		}

		for k, v := range m.GetGaugeAll(context.TODO()) {
			metric := storage.Metrics{
				ID:    k,
				MType: "gauge",
				Value: &v,
			}
			data, err := json.Marshal(&metric)
			if err != nil {
				logger.Log.Warn("Save storage marshal JSON", zap.Error(err))
				continue
			}
			file.Write(data)
			file.Write([]byte("\n"))
		}

		for k, v := range m.GetCounterAll(context.TODO()) {
			metric := storage.Metrics{
				ID:    k,
				MType: "counter",
				Delta: &v,
			}
			data, err := json.Marshal(&metric)
			if err != nil {
				logger.Log.Warn("Save storage marshal JSON", zap.Error(err))
				continue
			}
			file.Write(data)
			file.Write([]byte("\n"))
		}

		file.Close()

		logger.Log.Info("save storage", zap.String("file", fileStoragePath))
	}
}
func (m *Storage) PingStorage(ctx context.Context) error {
	return errors.New("not db")
}
