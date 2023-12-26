package agent

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/oktavarium/go-gauger/internal/shared"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"golang.org/x/sync/errgroup"
)

type metrics struct {
	gauges   map[string]float64
	counters map[string]int64
}

// NewMetrics конструктор хранилища метрик между отправками на сервер
func NewMetrics() metrics {
	return metrics{
		make(map[string]float64),
		make(map[string]int64),
	}
}

// readMetrics - читает все метрики через runtime.ReadMemStats
func readMetrics(m metrics) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.gauges["Alloc"] = float64(memStats.Alloc)
	m.gauges["TotalAlloc"] = float64(memStats.TotalAlloc)
	m.gauges["Sys"] = float64(memStats.Sys)
	m.gauges["Lookups"] = float64(memStats.Lookups)
	m.gauges["Frees"] = float64(memStats.Frees)
	m.gauges["Mallocs"] = float64(memStats.Mallocs)
	m.gauges["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.gauges["HeapSys"] = float64(memStats.HeapSys)
	m.gauges["HeapIdle"] = float64(memStats.HeapIdle)
	m.gauges["HeapInuse"] = float64(memStats.HeapInuse)
	m.gauges["HeapReleased"] = float64(memStats.HeapReleased)
	m.gauges["HeapObjects"] = float64(memStats.HeapObjects)
	m.gauges["StackInuse"] = float64(memStats.StackInuse)
	m.gauges["StackSys"] = float64(memStats.StackSys)
	m.gauges["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.gauges["MSpanSys"] = float64(memStats.MSpanSys)
	m.gauges["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.gauges["MCacheSys"] = float64(memStats.MCacheSys)
	m.gauges["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.gauges["GCSys"] = float64(memStats.GCSys)
	m.gauges["OtherSys"] = float64(memStats.OtherSys)
	m.gauges["NextGC"] = float64(memStats.NextGC)
	m.gauges["LastGC"] = float64(memStats.LastGC)
	m.gauges["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.gauges["NumGC"] = float64(memStats.NumGC)
	m.gauges["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.gauges["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	m.gauges["RandomValue"] = rand.Float64()

	m.counters["PollCount"]++

	return nil
}

// readMetrics - читает метрики через mem.VirtualMemory
func readPsMetrics(m metrics) error {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	m.gauges["TotalMemory"] = float64(vm.Total)
	m.gauges["FreeMemory"] = float64(vm.Free)

	cpu, err := cpu.Percent(0, true)
	if err != nil {
		return err
	}

	m.gauges["CPUutilization1"] = float64(cpu[0])

	return nil
}

// packMetrics - упаковывает все метики методом compressMetrics
func packMetrics(m metrics) ([]byte, error) {
	allMetrics := make([]shared.Metric, 0, len(m.gauges)+len(m.counters))
	for k, v := range m.gauges {
		v := v
		allMetrics = append(allMetrics, shared.NewGaugeMetric(k, &v))
	}

	for k, v := range m.counters {
		v := v
		allMetrics = append(allMetrics, shared.NewCounterMetric(k, &v))
	}

	compressedMetrics, err := compressMetrics(allMetrics)

	return compressedMetrics, err
}

func collector(
	ctx context.Context,
	collect func(metrics) error,
	eg *errgroup.Group,
	d time.Duration,
) chan []byte {

	chOut := make(chan []byte)
	metrics := NewMetrics()
	ticker := time.NewTicker(d)

	eg.Go(func() error {
		defer close(chOut)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := collect(metrics); err != nil {
					return err
				}
				packedMatrics, err := packMetrics(metrics)
				if err != nil {
					return err
				}
				chOut <- packedMatrics
			case <-ctx.Done():
				return nil
			}
		}
	})

	return chOut
}
