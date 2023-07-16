package handler

import (
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

type Handler struct {
	Storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		Storage: storage,
	}
}
