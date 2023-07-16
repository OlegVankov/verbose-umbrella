package storage

import (
	"context"
)

type Storage interface {
	UpdateGauge(ctx context.Context, name string, val float64)
	UpdateCounter(ctx context.Context, name string, val int64) (int64, error)
	GetGauge(ctx context.Context, name string) (float64, bool)
	GetCounter(ctx context.Context, name string) (int64, bool)
	GetGaugeAll(ctx context.Context) map[string]float64
	GetCounterAll(ctx context.Context) map[string]int64
	PingStorage(ctx context.Context) error
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
