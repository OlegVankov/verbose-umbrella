package main

import (
	"flag"
	"os"
)

var (
	serverAddr string
	level      string
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "127.0.0.1:8080", "адрес и порт сервера")
	flag.StringVar(&level, "l", "info", "уровень логирования")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverAddr = envRunAddr
	}
}
