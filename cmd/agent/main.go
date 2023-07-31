package main

import (
	"sync"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/monitor"
)

func main() {
	parseFlags()
	if err := logger.Initialize(level); err != nil {
		panic(err)
	}
	mtr := monitor.NewMonitor()
	logger.Log.Info("Agent", zap.String("running", serverAddr))

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go mtr.RunMonitor(pollInterval, wg)
	// monitor.SendMetrics(mtr, serverAddr, reportInterval)
	for i := 0; i < rateLimit; i++ {
		wg.Add(1)
		go monitor.SendBatch(mtr, serverAddr, reportInterval, key, wg)
	}

	wg.Wait()
	logger.Log.Info("Agent", zap.String("stopped", serverAddr))
}
