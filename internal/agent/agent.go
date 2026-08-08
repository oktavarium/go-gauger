package agent

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/oktavarium/go-gauger/internal/agent/internal/flags"
	"github.com/oktavarium/go-gauger/internal/shared"
	"golang.org/x/sync/errgroup"
)

// Run - запускает агент сбора и отправки метрик на сервер
func Run() error {
	flagsConfig, err := flags.LoadConfig()
	if err != nil {
		return fmt.Errorf("error on loading config: %w", err)
	}

	var publicKey *rsa.PublicKey
	if len(flagsConfig.CryptoKey) != 0 {
		publicKeyData, err := os.ReadFile(flagsConfig.CryptoKey)
		if err != nil {
			return fmt.Errorf("error on reading public key file: %w", err)
		}

		pkPEM, _ := pem.Decode(publicKeyData)
		if pkPEM.Type != "RSA PUBLIC KEY" {
			return fmt.Errorf("wrong key type: %w", err)
		}

		publicKey, err = x509.ParsePKCS1PublicKey(pkPEM.Bytes)
		if err != nil {
			return fmt.Errorf("error parsing public key: %w", err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	eg, egCtx := errgroup.WithContext(ctx)

	chMetrics := collector(
		egCtx,
		readMetrics,
		eg,
		flagsConfig.PollInterval,
	)
	chPsMetrics := collector(
		egCtx,
		readPsMetrics,
		eg,
		flagsConfig.PollInterval,
	)

	unitedCh := fanIn(chMetrics, chPsMetrics)
	for i := 0; i < flagsConfig.RateLimit; i++ {
		eg.Go(func() error {
			err := sender(
				egCtx,
				flagsConfig.Address,
				flagsConfig.GrpcAddress,
				flagsConfig.UseGRPC,
				flagsConfig.HashKey,
				publicKey,
				flagsConfig.ReportInterval,
				unitedCh)
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

// fanIn - метод мультиплексирования входящих данных от множества
func fanIn(chs ...<-chan []shared.Metric) <-chan []shared.Metric {
	chOut := make(chan []shared.Metric, len(chs))
	var wg sync.WaitGroup
	wg.Add(len(chs))

	output := func(ch <-chan []shared.Metric) {
		defer wg.Done()
		for v := range ch {
			chOut <- v
		}
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
