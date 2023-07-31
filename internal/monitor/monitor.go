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
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

type Monitor struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
	RandomValue   float64
	PollCount     int64

	TotalMemory     float64
	FreeMemory      float64
	CPUutilization1 float64
}

func NewMonitor() *Monitor {
	return &Monitor{
		PollCount:   0,
		RandomValue: rand.Float64(),
	}
}

func (m *Monitor) RunMonitor(pollInterval int, wg *sync.WaitGroup) {
	defer wg.Done()
	rtm := runtime.MemStats{}
	for {
		v, err := mem.VirtualMemory()
		if err != nil {
			logger.Log.Warn("Error get virtual memory", zap.Error(err))
			continue
		}
		c, err := cpu.Percent(time.Second, true)
		if err != nil {
			logger.Log.Warn("Error get cpu percent", zap.Error(err))
			continue
		}
		runtime.ReadMemStats(&rtm)
		m.Alloc = float64(rtm.Alloc)
		m.BuckHashSys = float64(rtm.BuckHashSys)
		m.Frees = float64(rtm.Frees)
		m.GCCPUFraction = rtm.GCCPUFraction
		m.GCSys = float64(rtm.GCSys)
		m.HeapAlloc = float64(rtm.HeapAlloc)
		m.HeapIdle = float64(rtm.HeapIdle)
		m.HeapInuse = float64(rtm.HeapInuse)
		m.HeapObjects = float64(rtm.HeapObjects)
		m.HeapReleased = float64(rtm.HeapReleased)
		m.HeapSys = float64(rtm.HeapSys)
		m.LastGC = float64(rtm.LastGC)
		m.Lookups = float64(rtm.Lookups)
		m.MCacheInuse = float64(rtm.MCacheInuse)
		m.MCacheSys = float64(rtm.MCacheSys)
		m.MSpanInuse = float64(rtm.MSpanInuse)
		m.MSpanSys = float64(rtm.MSpanSys)
		m.Mallocs = float64(rtm.Mallocs)
		m.NextGC = float64(rtm.NextGC)
		m.NumForcedGC = float64(rtm.NumForcedGC)
		m.NumGC = float64(rtm.NumGC)
		m.OtherSys = float64(rtm.OtherSys)
		m.PauseTotalNs = float64(rtm.PauseTotalNs)
		m.StackInuse = float64(rtm.StackInuse)
		m.StackSys = float64(rtm.StackSys)
		m.Sys = float64(rtm.Sys)
		m.TotalAlloc = float64(rtm.TotalAlloc)
		m.RandomValue = rand.Float64()
		m.PollCount++ // делаем инкремент каждые pollInterval секунд

		m.TotalMemory = float64(v.Total)
		m.FreeMemory = float64(v.Free)
		m.CPUutilization1 = c[0]

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

		metric := storage.Metrics{ID: typeField.Name}

		switch v := valueField.Interface().(type) {
		case int64:
			metric.MType = "counter"
			delta := v
			metric.Delta = &delta
		case float64:
			metric.MType = "gauge"
			value := v
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

func (m *Monitor) GetMetrics() []storage.Metrics {
	metrics := []storage.Metrics{}
	val := reflect.ValueOf(m).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		metric := storage.Metrics{ID: typeField.Name}

		switch v := valueField.Interface().(type) {
		case int64:
			metric.MType = "counter"
			delta := v
			metric.Delta = &delta
		case float64:
			metric.MType = "gauge"
			value := v
			metric.Value = &value
		}
		metrics = append(metrics, metric)
	}
	return metrics
}
