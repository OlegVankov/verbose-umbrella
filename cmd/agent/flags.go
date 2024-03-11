package main

import (
	"flag"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
)

var (
	serverAddr     string
	pollInterval   int
	reportInterval int
	level          string
	key            string
	rateLimit      int
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "адрес и порт сервера принимающего метрики")
	flag.IntVar(&pollInterval, "p", 2, "частота опроса метрик")
	flag.IntVar(&reportInterval, "r", 10, "частота отправки метрик на сервер")
	flag.IntVar(&rateLimit, "l", 5, "количество одновременно исходящих запросов")
	flag.StringVar(&level, "i", "info", "уровень логирования")
	flag.StringVar(&key, "k", "", "ключ для вычисления хеша")
	flag.Parse()
	getEnv()
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
	if envKey := os.Getenv("KEY"); envKey != "" {
		key = envKey
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		var err error
		rateLimit, err = strconv.Atoi(envRateLimit)
		if err != nil {
			logger.Log.Fatal("Get ENV RATE_LIMIT", zap.Error(err))
		}
	}
}
