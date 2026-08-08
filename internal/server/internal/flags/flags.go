package flags

import (
	"errors"
	"flag"
	"time"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
)

// config - структура хранения настроек сервиса
type Config struct {
	Address          string        `env:"ADDRESS" json:"address"`           // адрес и порт работы сервиса метрик
	GrpcAddress      string        `env:"GRPC_ADDRESS" json:"grpc_address"` // адрес и порт grpc-сервера
	LogLevel         string        `env:"LOGLEVEL"`                         // уровень логирования
	StoreIntervalInt int           `env:"STORE_INTERVAL"`                   // интервал сброса метрик в файл
	StoreInterval    time.Duration `json:"store_interval"`
	FilePath         string        `env:"FILE_STORAGE_PATH" json:"store_file"` // путь к файлу хранилища
	Restore          bool          `env:"RESTORE" json:"restore"`              // требуется ли восстановление при старте сервиса
	DatabaseDSN      string        `env:"DATABASE_DSN" json:"database_dsn"`    // DSN подключения к сервису posgtresql
	HashKey          string        `env:"KEY"`                                 // ключ аутентификации
	CryptoKey        string        `env:"CRYPTO_KEY" json:"crypto_key"`        // файл с приватным ключом сервера
	Config           string        `env:"CONFIG"`                              // файл с конфигурацией
	TrustedSubnet    string        `env:"TRUSTE_SUBNET" json:"trusted_subnet"` // доверенная сеть
}

// loadConfig - загружает конфигурацию - из флагов и переменных окружения
func LoadConfig() (Config, error) {
	var config Config
	config.parseCli()
	if err := config.loadEnv(); err != nil {
		logger.Logger().Error("error",
			zap.String("func", "LoadConfig"),
			zap.Error(err))
	}
	if err := config.parseConfigFile(config.Config); err != nil {
		logger.Logger().Error("error",
			zap.String("func", "LoadConfig"),
			zap.Error(err))
	}

	if len(flag.Args()) > 0 {
		return config, errors.New("unrecognised flags")
	}

	if config.StoreInterval == 0 && config.StoreIntervalInt != 0 {
		config.StoreInterval = time.Duration(config.StoreIntervalInt) * time.Second
	}

	return config, nil
}
