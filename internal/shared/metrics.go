package shared

import (
	pbapi "github.com/oktavarium/go-gauger/api"
)

func ConvertDBMetricsToMetrics(metrics []Metric) []*pbapi.Metric {
	allMetrics := make([]*pbapi.Metric, 0, len(metrics))
	for _, m := range metrics {
		allMetrics = append(allMetrics, ConvertDBMetricToMetric(m))
	}

	return allMetrics
}

func ConvertDBMetricToMetric(metric Metric) *pbapi.Metric {
	m := &pbapi.Metric{
		Id:   metric.ID,
		Type: metric.MType,
	}

	switch metric.MType {
	case GaugeType:
		m.Value = float64(*metric.Value)
	case CounterType:
		m.Value = float64(*metric.Delta)
	}

	return m
}

func ConvertMetricsToDBMetrics(metrics []*pbapi.Metric) []Metric {
	allMetrics := make([]Metric, 0, len(metrics))
	for _, m := range metrics {
		allMetrics = append(allMetrics, ConvertMetricToDBMetric(m))
	}

	return allMetrics
}

func ConvertMetricToDBMetric(metric *pbapi.Metric) Metric {
	m := Metric{
		ID:    metric.GetId(),
		MType: metric.GetType(),
	}

	switch metric.GetType() {
	case GaugeType:
		m.Value = &metric.Value
	case CounterType:
		v := int64(metric.Value)
		m.Delta = &v
	}

	return m
}
