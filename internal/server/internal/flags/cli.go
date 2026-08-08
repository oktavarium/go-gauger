package flags

import "flag"

func (c *Config) parseCli() {
	flag.StringVar(&c.Address, "a", "localhost:8080",
		"address and port of server in notaion address:port")
	flag.StringVar(&c.LogLevel, "l", "info",
		"log level")
	flag.IntVar(&c.StoreIntervalInt, "i", 10,
		"store interval")
	flag.StringVar(&c.FilePath, "f", "/tmp/metrics-db.json",
		"file storage path")
	flag.BoolVar(&c.Restore, "r", true,
		"restore metrics")
	flag.StringVar(&c.DatabaseDSN, "d", "",
		"database connection string")
	flag.StringVar(&c.HashKey, "k", "",
		"key for hash")
	flag.StringVar(&c.CryptoKey, "crypto-key", "",
		"server private key file")
	flag.StringVar(&c.Config, "c", "",
		"path to config")
	flag.StringVar(&c.TrustedSubnet, "t", "",
		"trusted subnet")
	flag.StringVar(&c.GrpcAddress, "g", "localhost:3333",
		"address for grpc server")
	flag.Parse()
}
