package main

import (
	"flag"
	"os"
)

var serverAddr string

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "адрес и порт сервера")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverAddr = envRunAddr
	}
}
