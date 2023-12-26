package handlers

import (
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

// NewHandler - конструтор для типа управляющего всеми эндпоинтами сервиса
func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}
