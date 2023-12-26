package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/shared"
)

// compressMetrics - сжимает метрики через gzip
func compressMetrics(metrics []shared.Metric) ([]byte, error) {
	var compressedJSON bytes.Buffer
	wr := gzip.NewWriter(&compressedJSON)

	encoder := json.NewEncoder(wr)
	if err := encoder.Encode(metrics); err != nil {
		return nil, fmt.Errorf("error occured on encoding metric: %w", err)
	}

	wr.Close()
	return compressedJSON.Bytes(), nil
}
