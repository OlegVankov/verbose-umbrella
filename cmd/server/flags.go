package main

import (
	"flag"
)

var serverAddr string

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "адрес и порт сервера")
	flag.Parse()
}
