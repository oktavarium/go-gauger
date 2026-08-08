package flags

import (
	"errors"
	"flag"
	"strings"
	"time"

	"github.com/oktavarium/go-gauger/internal/agent/internal/logger"
	"go.uber.org/zap"
)

// config - структура хранения настроек сервиса
type Config struct {
	Address           string        `env:"ADDRESS" json:"address"`           // адрес сервиса сбора метрик
	GrpcAddress       string        `env:"GRPC_ADDRESS" json:"grpc_address"` // grpc-адрес сервиса сбора метрик
	UseGRPC           bool          `env:"USE_GRPC"`                         // использовать grpc вместо http
	ReportIntervalInt int           `env:"REPORT_INTERVAL"`                  // интервал отправки метрик
	PollIntervalInt   int           `env:"POLL_INTERVAL"`                    // интервал сбора метрик
	HashKey           string        `env:"KEY"`                              // ключ аутентификации
	RateLimit         int           `env:"RATE_LIMIT"`                       // ограничение на количество поток
	CryptoKey         string        `env:"CRYPTO_KEY" json:"crypto_key"`     // файл с публичным ключом сервера
	ReportInterval    time.Duration `json:"report_interval"`
	PollInterval      time.Duration `json:"poll_interval"`
	Config            string        `env:"CONFIG"` // файл с конфигурацией
}

// loadConfig - загружает конфигурацию - из флагов и переменных окружения
func LoadConfig() (Config, error) {
	var config Config
	config.parseCli()
	if err := config.parseConfigFile(config.Config); err != nil {
		logger.Logger().Error("error",
			zap.String("func", "LoadConfig"),
			zap.Error(err))
	}
	if err := config.loadEnv(); err != nil {
		logger.Logger().Error("error",
			zap.String("func", "LoadConfig"),
			zap.Error(err))
	}

	if len(flag.Args()) > 0 {
		return config, errors.New("unrecognised flags")
	}

	if config.RateLimit <= 0 {
		config.RateLimit = 1
	}

	config.PollInterval = time.Duration(config.PollIntervalInt) * time.Second
	config.ReportInterval = time.Duration(config.ReportIntervalInt) * time.Second

	if !strings.HasPrefix(config.Address, "http://") {
		config.Address = "http://" + config.Address
	}

	return config, nil
}
