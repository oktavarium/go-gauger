package agent

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

// config - структура хранения настроек агента
type config struct {
	Address           string `env:"ADDRESS"`         // адрес сервиса сбора метрик
	ReportIntervalInt int    `env:"REPORT_INTERVAL"` // интервал отправки метрик
	PollIntervalInt   int    `env:"POLL_INTERVAL"`   // интервал сбора метрик
	HashKey           string `env:"KEY"`             // ключ аутентификации
	RateLimit         int    `env:"RATE_LIMIT"`      // ограничение на количество поток
	ReportInterval    time.Duration
	PollInterval      time.Duration
}

// loadConfig - загружает конфигурацию - из флагов и переменных окружения
func loadConfig() (config, error) {
	var flagsConfig config
	flag.StringVar(&flagsConfig.Address, "a", "localhost:8080",
		"address and port of server's endpoint in notaion address:port")
	flag.IntVar(&flagsConfig.ReportIntervalInt, "r", 2,
		"report interval in seconds")
	flag.IntVar(&flagsConfig.PollIntervalInt, "p", 2,
		"poll interval in seconds")
	flag.StringVar(&flagsConfig.HashKey, "k", "",
		"key for hash")
	flag.IntVar(&flagsConfig.RateLimit, "l", 1,
		"requests limit")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return flagsConfig, fmt.Errorf("error on parsing env parameters: %w", err)
	}

	if len(flag.Args()) > 0 {
		return flagsConfig, errors.New("unrecognised flags")
	}

	if flagsConfig.RateLimit <= 0 {
		flagsConfig.RateLimit = 1
	}

	flagsConfig.PollInterval = time.Duration(flagsConfig.PollIntervalInt) * time.Second
	flagsConfig.ReportInterval = time.Duration(flagsConfig.ReportIntervalInt) * time.Second

	flagsConfig.Address = "http://" + flagsConfig.Address

	return flagsConfig, nil
}
