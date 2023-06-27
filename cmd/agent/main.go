package main

import (
	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	mtr "github.com/OlegVankov/verbose-umbrella/internal/monitor"
	"go.uber.org/zap"
)

func main() {
	parseFlags()
	getEnv()
	if err := logger.Initialize(level); err != nil {
		panic(err)
	}
	monitor := mtr.NewMonitor()
	logger.Log.Info("Running agent", zap.String("sender", serverAddr))
	go monitor.RunMonitor(pollInterval)
	mtr.SendMetrics(monitor, serverAddr, reportInterval)
}
