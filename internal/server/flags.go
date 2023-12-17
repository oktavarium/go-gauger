package server

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Address       string        `env:"ADDRESS"`           // адрес и порт работы сервиса метрик
	LogLevel      string        `env:"LOGLEVEL"`          // уровень логирования
	StoreInterval time.Duration `env:"STORE_INTERVAL"`    // интервал сброса метрик в файл
	FilePath      string        `env:"FILE_STORAGE_PATH"` // путь к файлу хранилища
	Restore       bool          `env:"RESTORE"`           // требуется ли восстановление при старте сервиса
	DatabaseDSN   string        `env:"DATABASE_DSN"`      // DSN подключения к сервису posgtresql
	HashKey       string        `env:"KEY"`               // ключ аутентификации
}

func loadConfig() (config, error) {
	var flagsConfig config
	flag.StringVar(&flagsConfig.Address, "a", "localhost:8080",
		"address and port of server in notaion address:port")
	flag.StringVar(&flagsConfig.LogLevel, "l", "info",
		"log level")
	flag.DurationVar(&flagsConfig.StoreInterval, "i", 300*time.Second,
		"store interval")
	flag.StringVar(&flagsConfig.FilePath, "f", "/tmp/metrics-db.json",
		"file storage path")
	flag.BoolVar(&flagsConfig.Restore, "r", true,
		"restore metrics")
	flag.StringVar(&flagsConfig.DatabaseDSN, "d", "",
		"database connection string")
	flag.StringVar(&flagsConfig.HashKey, "k", "",
		"key for hash")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return flagsConfig, fmt.Errorf("error on parsing env parameters: %w", err)
	}

	if len(flag.Args()) > 0 {
		return flagsConfig, errors.New("unrecognised flags")
	}

	return flagsConfig, nil
}
