package main

import (
	"net/http"
	"time"

	mtr "github.com/OlegVankov/verbose-umbrella/internal/monitor"
)

func main() {
	parseFlags()
	getEnv()
	monitor := mtr.NewMonitor()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	go monitor.RunMonitor(pollInterval)
	mtr.SendMetrics(client, monitor, serverAddr, reportInterval)
}
