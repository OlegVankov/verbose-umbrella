package handler

import (
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

type Handler struct {
	Storage storage.Storage
	Key     string
}

func NewHandler(storage storage.Storage, key string) *Handler {
	return &Handler{
		Storage: storage,
		Key:     key,
	}
}
