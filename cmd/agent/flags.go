package main

import (
	"flag"
)

var (
	flagRunAddr    string
	pollInterval   int
	reportInterval int
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&pollInterval, "p", 2, "частота отправки метрик на сервер")
	flag.IntVar(&reportInterval, "r", 10, "частота опроса метрик")
	flag.Parse()
}
