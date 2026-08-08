package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func (s *storage) GetGauger(ctx context.Context, name string) (float64, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	val, ok := s.gauge[name]
	return val, ok
}

func (s *storage) GetCounter(ctx context.Context, name string) (int64, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	val, ok := s.counter[name]
	return val, ok
}

func (s *storage) GetAll(ctx context.Context) ([]byte, error) {
	allMetrics, err := s.GetAllMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("error on getting all metrics: %w", err)
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)

	for _, v := range allMetrics {
		err := encoder.Encode(&v)
		if err != nil {
			return nil, fmt.Errorf("error on encoding data: %w", err)
		}
	}
	return buffer.Bytes(), nil
}

func (s *storage) GetAllMetrics(ctx context.Context) ([]shared.Metric, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	allMetrics := make([]shared.Metric, 0, len(s.gauge)+len(s.counter))

	for k, v := range s.gauge {
		v := v
		allMetrics = append(allMetrics, shared.NewGaugeMetric(k, &v))

	}
	for k, v := range s.counter {
		v := v
		allMetrics = append(allMetrics, shared.NewCounterMetric(k, &v))
	}

	return allMetrics, nil
}
