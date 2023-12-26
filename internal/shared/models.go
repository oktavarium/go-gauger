// Модуль shared содержит основной типа работы с метриками и соответствующие конструкторы
package shared

// MetricType - тип для метрик
type MetricType string

const (
	GaugeType   string = "gauge"
	CounterType string = "counter"
)

// Metric - структура, содержащая идентификатор, тип и значение соответствующей метрики
type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewGaugeMetric - конструктор метрики типа gauge
func NewGaugeMetric(id string, val *float64) Metric {
	return Metric{
		ID:    id,
		MType: GaugeType,
		Value: val,
	}
}

// NewCounterMetric - конструктор метрики типа counter
func NewCounterMetric(id string, val *int64) Metric {
	return Metric{
		ID:    id,
		MType: CounterType,
		Delta: val,
	}
}

// NewGaugeMetric - конструктор пустой метрики типа gauge
func NewEmptyGaugeMetric() Metric {
	return Metric{
		MType: GaugeType,
	}
}

// NewCounterMetric - конструктор пустой метрики типа counter
func NewEmptyCounterMetric() Metric {
	return Metric{
		MType: CounterType,
	}
}
