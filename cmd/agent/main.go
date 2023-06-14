package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/OlegVankov/verbose-umbrella/internal/storage"
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

var m = Monitor{
	PollCount: 0,
}

func RunMonitor(duration int) {

	rtm := runtime.MemStats{}
	for {
		<-time.After(time.Duration(duration) * time.Second)

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
		m.PollCount += 1
	}
}

func (m *Monitor) getRoutes() []string {
	urls := []string{}

	val := reflect.ValueOf(m).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		urls = append(urls, fmt.Sprintf("http://%s/%s/%s/%s/%v",
			serverAddr, "update",
			strings.ToLower(strings.Split(typeField.Type.String(), ".")[1]),
			typeField.Name, valueField.Interface()))
	}

	return urls
}

func SendMetrics(client *http.Client, duration int) {

	for {

		<-time.After(time.Duration(duration) * time.Second)

		for _, url := range m.getRoutes() {
			req, err := http.NewRequest(http.MethodPost, url, nil)
			req.Header.Set("Host", serverAddr)
			req.Header.Set("Content-Type", "text/plain")
			if err != nil {
				log.Println("StatusCode:", req.Response.Status, err.Error())
				return
			}
			r, err := client.Do(req)
			if err != nil {
				log.Println(err.Error())
				return
			}
			r.Body.Close()
		}

	}

}

func main() {
	parseFlags()
	getEnv()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	go RunMonitor(pollInterval)
	SendMetrics(client, reportInterval)
}
