package storage

import (
	"context"

	"github.com/oktavarium/go-gauger/internal/shared"
)

// Storage - интерфейс абстрактного хранилища
type Storage interface {
	// SaveGauge - сохраняет метрику типа gauge
	SaveGauge(context.Context, string, float64) error
	// UpdateCounter - обновляент метрику типа counter
	UpdateCounter(context.Context, string, int64) (int64, error)
	// GetGauger - получает метрику типа gauge
	GetGauger(context.Context, string) (float64, bool)
	// GetCounter - получает метрику типа counter
	GetCounter(context.Context, string) (int64, bool)
	// GetAll - получает все метрики в бинарном виде
	GetAll(context.Context) ([]byte, error)
	// GetAllMetrics - получает все метрики
	GetAllMetrics(context.Context) ([]shared.Metric, error)
	// Ping - проверяет доступность хранилища
	Ping(context.Context) error
	// BatchUpdate - обновление метрик пачкой
	BatchUpdate(context.Context, []shared.Metric) error
}
