package storage

type Gauge float64

type Counter int64

type Storage interface {
	UpdateGauge(name string, val Gauge)
	UpdateCounter(name string, val Counter)
	GetGauge(name string) (Gauge, bool)
	GetCounter(name string) (Counter, bool)
	GetGaugeAll() map[string]Gauge
	GetCounterAll() map[string]Counter
}

type MemStorage struct {
	Gauge   map[string]Gauge
	Counter map[string]Counter
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Gauge:   map[string]Gauge{},
		Counter: map[string]Counter{},
	}
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

func (m *MemStorage) GetGaugeAll() map[string]Gauge {
	return m.Gauge
}

func (m *MemStorage) GetCounterAll() map[string]Counter {
	return m.Counter
}
