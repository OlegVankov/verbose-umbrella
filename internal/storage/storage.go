package storage

type Gauge float64

type Counter int64

type MemStorage struct {
	Gauge   map[string]Gauge
	Counter map[string]Counter
}

func (m *MemStorage) UpdateGauge(name string, val Gauge) {
	m.Gauge[name] = val
}

func (m *MemStorage) UpdateCounter(name string, val Counter) {
	m.Counter[name] += val
}
