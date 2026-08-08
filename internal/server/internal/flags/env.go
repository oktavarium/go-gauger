package flags

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

func (c *Config) loadEnv() error {
	var config Config
	if err := env.Parse(&config); err != nil {
		return fmt.Errorf("error on parsing env parameters: %w", err)
	}

	if len(config.Address) != 0 {
		c.Address = config.Address
	}

	if len(config.LogLevel) != 0 {
		c.LogLevel = config.LogLevel
	}

	if len(config.FilePath) != 0 {
		c.FilePath = config.FilePath
	}

	if len(config.DatabaseDSN) != 0 {
		c.DatabaseDSN = config.DatabaseDSN
	}

	if len(config.HashKey) != 0 {
		c.HashKey = config.HashKey
	}

	if len(config.CryptoKey) != 0 {
		c.CryptoKey = config.CryptoKey
	}

	if config.StoreIntervalInt != 0 {
		c.StoreIntervalInt = config.StoreIntervalInt
	}

	if config.StoreInterval != 0 {
		c.StoreInterval = config.StoreInterval
	}

	if len(config.TrustedSubnet) != 0 {
		c.TrustedSubnet = config.TrustedSubnet
	}

	if len(config.GrpcAddress) != 0 {
		c.GrpcAddress = config.GrpcAddress
	}

	return nil
}
