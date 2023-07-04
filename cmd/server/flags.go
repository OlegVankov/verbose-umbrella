package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"go.uber.org/zap"
)

var (
	serverAddr      string
	fileStoragePath string
	storeInterval   int
	restore         bool
	level           string
)

func parseFlags() {
	flag.StringVar(&serverAddr, "a", "localhost:8080", "адрес и порт сервера")
	flag.IntVar(&storeInterval, "i", 300, "время (в секундах), по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "файл, куда сохраняются текущие значения")
	flag.BoolVar(&restore, "r", true, "загружать или нет сохранённые значения из файла при старте сервера")
	flag.StringVar(&level, "l", "info", "уровень логирования")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		serverAddr = envRunAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		var err error
		storeInterval, err = strconv.Atoi(envStoreInterval)
		if err != nil {
			logger.Log.Fatal("Get ENV STORE_INTERVAL", zap.Error(err))
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		var err error
		restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			logger.Log.Fatal("Get ENV RESTORE", zap.Error(err))
		}
	}
}
