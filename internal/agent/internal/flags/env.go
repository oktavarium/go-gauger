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

	if len(config.HashKey) != 0 {
		c.HashKey = config.HashKey
	}

	if len(config.CryptoKey) != 0 {
		c.CryptoKey = config.CryptoKey
	}

	if config.PollIntervalInt != 0 {
		c.PollIntervalInt = config.PollIntervalInt
	}

	if config.ReportIntervalInt != 0 {
		c.ReportIntervalInt = config.ReportIntervalInt
	}

	if config.PollInterval != 0 {
		c.PollInterval = config.PollInterval
	}

	if config.ReportInterval != 0 {
		c.ReportInterval = config.ReportInterval
	}

	if len(config.GrpcAddress) != 0 {
		c.GrpcAddress = config.GrpcAddress
	}

	return nil
}
