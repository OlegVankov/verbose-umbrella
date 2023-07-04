package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"go.uber.org/zap"
)

var (
	serverAddr     string
	pollInterval   int
	reportInterval int
	level          string
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "адрес и порт сервера принимающего метрики")
	flag.IntVar(&pollInterval, "p", 2, "частота опроса метрик")
	flag.IntVar(&reportInterval, "r", 10, "частота отправки метрик на сервер")
	flag.StringVar(&level, "l", "info", "уровень логирования")
	flag.Parse()
}

func getEnv() {
	if envServerAddr := os.Getenv("ADDRESS"); envServerAddr != "" {
		serverAddr = envServerAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		var err error
		reportInterval, err = strconv.Atoi(envReportInterval)
		if err != nil {
			logger.Log.Fatal("Get ENV REPORT_INTERVAL", zap.Error(err))
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		var err error
		pollInterval, err = strconv.Atoi(envPollInterval)
		if err != nil {
			logger.Log.Fatal("Get ENV POLL_INTERVAL", zap.Error(err))
		}
	}
}
