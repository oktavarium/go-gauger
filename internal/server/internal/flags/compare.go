package flags

// compare сравнивает оригинальный конфиг с предоставленным и обновляет значения
func (c *Config) compare(another Config) {
	if len(c.Address) == 0 && len(another.Address) != 0 {
		c.Address = another.Address
	}

	if len(c.LogLevel) == 0 && len(another.LogLevel) != 0 {
		c.LogLevel = another.LogLevel
	}

	if len(c.FilePath) == 0 && len(another.FilePath) != 0 {
		c.FilePath = another.FilePath
	}

	if len(c.DatabaseDSN) == 0 && len(another.DatabaseDSN) != 0 {
		c.DatabaseDSN = another.DatabaseDSN
	}

	if len(c.HashKey) == 0 && len(another.HashKey) != 0 {
		c.HashKey = another.HashKey
	}

	if len(c.CryptoKey) == 0 && len(another.CryptoKey) != 0 {
		c.CryptoKey = another.CryptoKey
	}

	if c.StoreIntervalInt == 0 && another.StoreIntervalInt != 0 {
		c.StoreIntervalInt = another.StoreIntervalInt
	}

	if c.StoreInterval == 0 && another.StoreInterval != 0 {
		c.StoreInterval = another.StoreInterval
	}

	if len(c.TrustedSubnet) == 0 && len(another.TrustedSubnet) != 0 {
		c.TrustedSubnet = another.TrustedSubnet
	}

	if len(c.GrpcAddress) == 0 && len(another.GrpcAddress) != 0 {
		c.GrpcAddress = another.GrpcAddress
	}
}
