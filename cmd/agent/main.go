package main

import (
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/monitor"
	"go.uber.org/zap"
)

func main() {
	parseFlags()
	getEnv()
	if err := logger.Initialize(level); err != nil {
		panic(err)
	}
	mtr := monitor.NewMonitor()
	logger.Log.Info("Running agent", zap.String("sender", serverAddr))
	go mtr.RunMonitor(pollInterval)
	monitor.SendMetrics(mtr, serverAddr, reportInterval)
}
