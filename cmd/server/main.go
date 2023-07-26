package main

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/server"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
	"github.com/OlegVankov/verbose-umbrella/internal/storage/memory"
	"github.com/OlegVankov/verbose-umbrella/internal/storage/pg"
)

func main() {
	parseFlags()

	if err := logger.Initialize(level); err != nil {
		panic(err)
	}

	stor, err := setStorage(context.Background())
	if err != nil {
		logger.Log.Fatal("storage", zap.Error(err))
	}

	if err := server.Run(serverAddr, stor, key); err != nil {
		logger.Log.Fatal("server", zap.Error(err))
	}

	logger.Log.Info("server gracefully shutdown complete")
}

func setStorage(ctx context.Context) (storage.Storage, error) {
	if databaseDSN != "" {
		db, err := sqlx.Open("pgx", databaseDSN)
		if err != nil {
			return nil, err
		}
		store := pg.NewStorage(db)
		if err := store.Bootstrap(ctx); err != nil {
			return nil, err
		}
		return store, nil
	}
	if fileStoragePath != "" {
		store := memory.NewStorage()
		if err := store.RestoreStorage(fileStoragePath); err != nil {
			return nil, err
		}
		go store.SaveStorage(fileStoragePath, storeInterval)
		return store, nil
	}
	return nil, errors.New("error set storage")
}
