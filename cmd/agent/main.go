package main

import (
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
	go mtr.RunMonitor(pollInterval)
	// monitor.SendMetrics(mtr, serverAddr, reportInterval)
	monitor.SendBatch(mtr, serverAddr, reportInterval, key)
	logger.Log.Info("Agent", zap.String("stopped", serverAddr))
}
