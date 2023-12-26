package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGaugeMetric(t *testing.T) {
	id := "test_id"
	val := 0.28
	newMetric := NewGaugeMetric(id, &val)

	require.Equal(t, id, newMetric.ID)
	require.Equal(t, val, *newMetric.Value)
	require.Equal(t, GaugeType, newMetric.MType)
}

func TestNewCounterMetric(t *testing.T) {
	id := "test_id"
	var val int64 = 28
	newMetric := NewCounterMetric(id, &val)

	require.Equal(t, id, newMetric.ID)
	require.Equal(t, val, *newMetric.Delta)
	require.Equal(t, CounterType, newMetric.MType)
}

func TestNewEmptyGaugeMetric(t *testing.T) {
	newMetric := NewEmptyGaugeMetric()

	require.Empty(t, newMetric.ID)
	require.Nil(t, newMetric.Value)
	require.Equal(t, GaugeType, newMetric.MType)
}

func TestNewEmptyCounterMetric(t *testing.T) {
	newMetric := NewEmptyCounterMetric()

	require.Empty(t, newMetric.ID)
	require.Nil(t, newMetric.Delta)
	require.Equal(t, CounterType, newMetric.MType)
}
