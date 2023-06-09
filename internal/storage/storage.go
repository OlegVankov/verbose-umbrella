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

func (m *MemStorage) GetGauge(name string) (Gauge, bool) {
	val, ok := m.Gauge[name]
	return val, ok
}

func (m *MemStorage) GetCounter(name string) (Counter, bool) {
	val, ok := m.Counter[name]
	return val, ok
}
