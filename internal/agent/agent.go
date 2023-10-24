package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func Run() error {
	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error on loading config: %w", err)
	}

	eg, egCtx := errgroup.WithContext(context.Background())

	chMetrics := collector(
		egCtx,
		readMetrics,
		eg, time.Duration(flagsConfig.PollInterval*time.Second))
	chPsMetrics := collector(
		egCtx,
		readPsMetrics,
		eg, time.Duration(flagsConfig.PollInterval*time.Second))

	unitedCh := fanIn(chMetrics, chPsMetrics)
	for i := 0; i < flagsConfig.RateLimit; i++ {
		go sender(egCtx, flagsConfig.Address, flagsConfig.HashKey,
			flagsConfig.ReportInterval, unitedCh)
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func fanIn(chs ...chan []byte) chan []byte {
	chOut := make(chan []byte, len(chs))
	var wg sync.WaitGroup
	wg.Add(len(chs))

	output := func(ch chan []byte) {
		for v := range ch {
			chOut <- v
		}
		wg.Done()
	}

	for _, ch := range chs {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(chOut)
	}()

	return chOut
}
