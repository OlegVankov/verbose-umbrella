package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
)

func (m *MemStorage) RestoreStorage(fileStoragePath string) error {
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var metric Metrics
		err := json.Unmarshal(scanner.Bytes(), &metric)
		if err != nil {
			logger.Log.Warn("restore storage: json unmarshal",
				zap.Error(err))
			continue
		}
		switch metric.MType {
		case "counter":
			m.UpdateCounter(metric.ID, *metric.Delta)
		case "gauge":
			m.UpdateGauge(metric.ID, *metric.Value)
		}
	}

	logger.Log.Info("restore storage", zap.String("file", fileStoragePath))
	return nil
}
