package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oktavarium/go-gauger/internal/server/internal/flags"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
)

// Run - запускает сервис обработки метрик
func Run() error {
	flagsConfig, err := flags.LoadConfig()
	if err != nil {
		return fmt.Errorf("error on loading config: %w", err)
	}

	if err = logger.Init(flagsConfig.LogLevel); err != nil {
		return fmt.Errorf("error init logger: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	gs, err := gaugeserver.NewGaugerServer(
		ctx,
		flagsConfig.Address,
		flagsConfig.GrpcAddress,
		flagsConfig.FilePath,
		flagsConfig.Restore,
		flagsConfig.StoreInterval,
		flagsConfig.DatabaseDSN,
		flagsConfig.HashKey,
		flagsConfig.CryptoKey,
		flagsConfig.TrustedSubnet)
	if err != nil {
		return fmt.Errorf("error on creating gaugeserver: %w", err)
	}

	return gs.ListenAndServe()
}
