package flags

// compare сравнивает оригинальный конфиг с предоставленным и обновляет значения
func (c *Config) compare(another Config) {
	if len(c.Address) == 0 && len(another.Address) != 0 {
		c.Address = another.Address
	}

	if len(c.HashKey) == 0 && len(another.HashKey) != 0 {
		c.HashKey = another.HashKey
	}

	if len(c.CryptoKey) == 0 && len(another.CryptoKey) != 0 {
		c.CryptoKey = another.CryptoKey
	}

	if c.PollIntervalInt == 0 && another.PollIntervalInt != 0 {
		c.PollIntervalInt = another.PollIntervalInt
	}

	if c.ReportIntervalInt == 0 && another.ReportIntervalInt != 0 {
		c.ReportIntervalInt = another.ReportIntervalInt
	}

	if c.PollInterval == 0 && another.PollInterval != 0 {
		c.PollInterval = another.PollInterval
	}

	if c.ReportInterval == 0 && another.ReportInterval != 0 {
		c.ReportInterval = another.ReportInterval
	}

	if len(c.GrpcAddress) == 0 && len(another.GrpcAddress) != 0 {
		c.GrpcAddress = another.GrpcAddress
	}

	c.UseGRPC = another.UseGRPC
}
