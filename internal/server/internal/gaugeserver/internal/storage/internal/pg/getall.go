package pg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func (s *storage) GetAll(ctx context.Context) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)

	allMetrics, err := s.GetAllMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("error on getting metrics: %w", err)
	}

	for _, v := range allMetrics {
		err := encoder.Encode(&v)
		if err != nil {
			return nil, fmt.Errorf("error on encoding data: %w", err)
		}
	}
	return buffer.Bytes(), nil
}

func (s *storage) GetAllMetrics(ctx context.Context) ([]shared.Metric, error) {
	allMetrics := make([]shared.Metric, 0)

	rows, err := s.Query(ctx, "SELECT name, value FROM gauge")
	if err != nil {
		return nil, fmt.Errorf("error occured on selecting all gauge: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		m := shared.NewEmptyGaugeMetric()
		if err = rows.Scan(&m.ID, &m.Value); err != nil {
			return nil, fmt.Errorf("error occured on scanning gauge: %w", err)
		}
		allMetrics = append(allMetrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occured on selecting all gauge: %w", err)
	}

	rows, err = s.Query(ctx, "SELECT name, value FROM counter")
	if err != nil {
		return nil, fmt.Errorf("error occured on selecting all counters: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		m := shared.NewEmptyCounterMetric()
		if err = rows.Scan(&m.ID, &m.Delta); err != nil {
			return nil, fmt.Errorf("error occured on scanning counter: %w", err)
		}
		allMetrics = append(allMetrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occured on selecting all counter: %w", err)
	}

	return allMetrics, nil
}
