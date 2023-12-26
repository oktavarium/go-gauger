package memory

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"github.com/oktavarium/go-gauger/internal/shared"
	"go.uber.org/zap"
)

func (s *storage) restore() error {
	data, err := s.archive.Restore()
	if err != nil {
		return fmt.Errorf("error on restoring archive: %w", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var metrics shared.Metric
		err := json.Unmarshal(scanner.Bytes(), &metrics)
		if err != nil {
			return fmt.Errorf("error on restoring archive: %w", err)
		}
		switch metrics.MType {
		case string(shared.GaugeType):
			if err := s.SaveGauge(context.Background(), metrics.ID, *metrics.Value); err != nil {
				logger.Logger().Error("error",
					zap.String("func", "restore"),
					zap.Error(err),
				)
			}
		case string(shared.CounterType):
			if _, err := s.UpdateCounter(context.Background(), metrics.ID, *metrics.Delta); err != nil {
				logger.Logger().Error("error",
					zap.String("func", "restore"),
					zap.Error(err),
				)
			}
		}
	}

	return nil
}
