package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	serverAddr     string
	pollInterval   int
	reportInterval int
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "адрес и порт сервера принимающего метрики")
	flag.IntVar(&pollInterval, "p", 2, "частота опроса метрик")
	flag.IntVar(&reportInterval, "r", 10, "частота отправки метрик на сервер")
	flag.Parse()
}

func getEnv() {
	if envServerAddr := os.Getenv("ADDRESS"); envServerAddr != "" {
		serverAddr = envServerAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		reportInterval, _ = strconv.Atoi(envReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollInterval, _ = strconv.Atoi(envPollInterval)
	}

}
