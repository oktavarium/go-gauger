package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"github.com/oktavarium/go-gauger/internal/shared"
)

// UpdatesHandle - обновляет несколько метрик путем bactchUpdate
func (h *Handler) UpdatesHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var metrics []shared.Metric
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&metrics)
	if err != nil {
		logger.LogError("UpdatesHandle", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(metrics) == 0 {
		err = fmt.Errorf("empty metrics")
		logger.LogError("UpdatesHandle", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.BatchUpdate(r.Context(), metrics)
	if err != nil {
		logger.LogError("UpdatesHandle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metricStab := shared.Metric{}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(metricStab)
	if err != nil {
		logger.LogError("UpdatesHandle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
