package agent

import (
	"context"
	"testing"
	"time"

	"github.com/oktavarium/go-gauger/internal/shared"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

var compressedData = []uint8([]byte{
	0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0xff, 0x8a, 0xae, 0x56, 0xca, 0x4c,
	0x51, 0xb2, 0x52, 0x52, 0xd2, 0x51, 0x2a,
	0xa9, 0x2c, 0x48, 0x55, 0xb2, 0x52, 0x4a,
	0xce, 0x2f, 0xcd, 0x2b, 0x49, 0x2d, 0x52,
	0xaa, 0xd5, 0xc1, 0x90, 0x4a, 0x4f, 0x2c,
	0x4d, 0x4f, 0x55, 0xaa, 0x8d, 0xe5, 0x2,
	0x4, 0x0, 0x0, 0xff, 0xff, 0x3a, 0x11,
	0x91, 0xb9, 0x36, 0x0, 0x0, 0x0,
})

func TestCompressMetrics(t *testing.T) {
	metrics := make([]shared.Metric, 0)
	metrics = append(metrics, shared.NewEmptyCounterMetric(), shared.NewEmptyGaugeMetric())
	compressed, err := compressMetrics(metrics)

	require.Equal(t, compressedData, compressed)
	require.NoError(t, err)
}

func TestHashData(t *testing.T) {
	key := "key"
	data := "data"
	want := "5031fe3d989c6d1537a013fa6e739da23463fdaec3b70137d828e36ace221bd0"

	hash, err := hashData([]byte(key), []byte(data))

	require.Equal(t, want, hash)
	require.NoError(t, err)
}

func TestFanIn(t *testing.T) {
	in1 := make(chan []byte, 1)
	in2 := make(chan []byte, 1)

	out := fanIn(in1, in2)
	d1 := []byte("test")
	d2 := []byte("test")
	in1 <- d1
	in2 <- d2

	require.Equal(t, d1, <-out)
	require.Equal(t, d2, <-out)
}

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()

	require.Empty(t, m.gauges)
	require.Empty(t, m.counters)
}

func TestReadMetrics(t *testing.T) {
	m := NewMetrics()
	err := readMetrics(m)
	require.NoError(t, err)

	gaugeCount := 28

	require.NotEmpty(t, m)
	require.Equal(t, gaugeCount, len(m.gauges))
}

func TestPsMetrics(t *testing.T) {
	m := NewMetrics()
	err := readPsMetrics(m)
	require.NoError(t, err)

	psCount := 3

	require.NotEmpty(t, m)
	require.Equal(t, psCount, len(m.gauges))
}

func TestCollector(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	out := collector(ctx, readMetrics, &errgroup.Group{}, 1*time.Second)
	time.Sleep(3 * time.Second)
	cancel()
	data := <-out

	require.NotEmpty(t, data)
}

func TestLoadConfig(t *testing.T) {
	cfg, err := loadConfig()
	require.NoError(t, err)
	require.Equal(t, "http://localhost:8080", cfg.Address)
	require.Equal(t, "", cfg.HashKey)
	require.Equal(t, 2*time.Second, cfg.PollInterval)
	require.Equal(t, 1, cfg.RateLimit)
	require.Equal(t, 2*time.Second, cfg.ReportInterval)
}
