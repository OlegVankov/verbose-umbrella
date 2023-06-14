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
	flag.IntVar(&pollInterval, "p", 2, "частота отправки метрик на сервер")
	flag.IntVar(&reportInterval, "r", 10, "частота опроса метрик")
	flag.Parse()
}

func getEnv() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverAddr = envRunAddr
	}
	if envRI := os.Getenv("REPORT_INTERVAL"); envRI != "" {
		reportInterval, _ = strconv.Atoi(envRI)
	}
	if envPI := os.Getenv("POLL_INTERVAL"); envPI != "" {
		pollInterval, _ = strconv.Atoi(envPI)
	}

}
