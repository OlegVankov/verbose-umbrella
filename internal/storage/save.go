package storage

import (
	"encoding/json"
	"os"
	"time"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"go.uber.org/zap"
)

func (m *MemStorage) SaveStorage(fileStoragePath string, storeInterval int) {
	for {
		<-time.After(time.Duration(storeInterval) * time.Second)

		file, err := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			logger.Log.Error("error open file", zap.String("file name", fileStoragePath), zap.Error(err))
			return
		}

		for k, v := range m.GetGaugeAll() {
			value := GaugeToFloat(v)
			metric := Metrics{
				ID:    k,
				MType: "gauge",
				Value: &value,
			}
			data, err := json.Marshal(&metric)
			if err != nil {
				logger.Log.Warn("Save storage marshal JSON", zap.Error(err))
				continue
			}
			file.Write(data)
			file.Write([]byte("\n"))
		}

		for k, v := range m.GetCounterAll() {
			value := CounterToInt(v)
			metric := Metrics{
				ID:    k,
				MType: "counter",
				Delta: &value,
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
