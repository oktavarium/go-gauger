package handlers

import (
	"net/http"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
)

// GetHandle получает все доступные в данный момент метрики
func (h *Handler) GetHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data, err := h.storage.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		logger.LogError("GetHandler", err)
	}
}
