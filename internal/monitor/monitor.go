package monitor

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
	"go.uber.org/zap"
)

type Monitor struct {
	Alloc         storage.Gauge
	BuckHashSys   storage.Gauge
	Frees         storage.Gauge
	GCCPUFraction storage.Gauge
	GCSys         storage.Gauge
	HeapAlloc     storage.Gauge
	HeapIdle      storage.Gauge
	HeapInuse     storage.Gauge
	HeapObjects   storage.Gauge
	HeapReleased  storage.Gauge
	HeapSys       storage.Gauge
	LastGC        storage.Gauge
	Lookups       storage.Gauge
	MCacheInuse   storage.Gauge
	MCacheSys     storage.Gauge
	MSpanInuse    storage.Gauge
	MSpanSys      storage.Gauge
	Mallocs       storage.Gauge
	NextGC        storage.Gauge
	NumForcedGC   storage.Gauge
	NumGC         storage.Gauge
	OtherSys      storage.Gauge
	PauseTotalNs  storage.Gauge
	StackInuse    storage.Gauge
	StackSys      storage.Gauge
	Sys           storage.Gauge
	TotalAlloc    storage.Gauge
	RandomValue   storage.Gauge
	PollCount     storage.Counter
}

func NewMonitor() *Monitor {
	return &Monitor{
		PollCount:   0,
		RandomValue: storage.Gauge(rand.Float64()),
	}
}

func (m *Monitor) RunMonitor(pollInterval int) {
	rtm := runtime.MemStats{}
	for {
		runtime.ReadMemStats(&rtm)
		m.Alloc = storage.Gauge(rtm.Alloc)
		m.BuckHashSys = storage.Gauge(rtm.BuckHashSys)
		m.Frees = storage.Gauge(rtm.Frees)
		m.GCCPUFraction = storage.Gauge(rtm.GCCPUFraction)
		m.GCSys = storage.Gauge(rtm.GCSys)
		m.HeapAlloc = storage.Gauge(rtm.HeapAlloc)
		m.HeapIdle = storage.Gauge(rtm.HeapIdle)
		m.HeapInuse = storage.Gauge(rtm.HeapInuse)
		m.HeapObjects = storage.Gauge(rtm.HeapObjects)
		m.HeapReleased = storage.Gauge(rtm.HeapReleased)
		m.HeapSys = storage.Gauge(rtm.HeapSys)
		m.LastGC = storage.Gauge(rtm.LastGC)
		m.Lookups = storage.Gauge(rtm.Lookups)
		m.MCacheInuse = storage.Gauge(rtm.MCacheInuse)
		m.MCacheSys = storage.Gauge(rtm.MCacheSys)
		m.MSpanInuse = storage.Gauge(rtm.MSpanInuse)
		m.MSpanSys = storage.Gauge(rtm.MSpanSys)
		m.Mallocs = storage.Gauge(rtm.Mallocs)
		m.NextGC = storage.Gauge(rtm.NextGC)
		m.NumForcedGC = storage.Gauge(rtm.NumForcedGC)
		m.NumGC = storage.Gauge(rtm.NumGC)
		m.OtherSys = storage.Gauge(rtm.OtherSys)
		m.PauseTotalNs = storage.Gauge(rtm.PauseTotalNs)
		m.StackInuse = storage.Gauge(rtm.StackInuse)
		m.StackSys = storage.Gauge(rtm.StackSys)
		m.Sys = storage.Gauge(rtm.Sys)
		m.TotalAlloc = storage.Gauge(rtm.TotalAlloc)
		m.RandomValue = storage.Gauge(rand.Float64())
		m.PollCount++ // делаем инкремент каждые pollInterval секунд
		<-time.After(time.Duration(pollInterval) * time.Second)
	}
}

func (m *Monitor) GetRoutes(serverAddr string) []string {
	urls := []string{}
	val := reflect.ValueOf(m).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		uri := strings.ToLower(strings.Split(typeField.Type.String(), ".")[1])
		urls = append(urls, fmt.Sprintf("http://%s/%s/%s/%s/%v",
			serverAddr, "update", uri, typeField.Name, valueField.Interface()))
	}
	return urls
}

func (m *Monitor) GetBody() []*bytes.Buffer {
	body := []*bytes.Buffer{}
	val := reflect.ValueOf(m).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		mType := strings.ToLower(strings.Split(typeField.Type.String(), ".")[1])

		metric := storage.Metrics{
			ID:    typeField.Name,
			MType: mType,
		}

		switch v := valueField.Interface().(type) {
		case storage.Counter:
			delta := storage.CounterToInt(v)
			metric.Delta = &delta
		case storage.Gauge:
			value := storage.GaugeToFloat(v)
			metric.Value = &value
		}

		body = append(body, gzipBody(&metric))
	}
	return body
}

func (m *Monitor) resetPollCount() {
	m.PollCount = 0
}

func gzipBody(m *storage.Metrics) *bytes.Buffer {
	data, err := json.Marshal(m)
	if err != nil {
		logger.Log.Warn("JSON Marshal", zap.Error(err))
	}
	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	g.Write(data)
	g.Close()
	return &buf
}
