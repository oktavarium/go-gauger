package flags

import "flag"

func (c *Config) parseCli() {
	flag.StringVar(&c.Address, "a", "localhost:8080",
		"address and port of server's endpoint in notaion address:port")
	flag.StringVar(&c.GrpcAddress, "g", "localhost:3333",
		"address and port of server's grpc endpoint in notaion address:port")
	flag.BoolVar(&c.UseGRPC, "u", false,
		"use grpc instead of http")
	flag.IntVar(&c.ReportIntervalInt, "r", 2,
		"report interval in seconds")
	flag.IntVar(&c.PollIntervalInt, "p", 2,
		"poll interval in seconds")
	flag.StringVar(&c.HashKey, "k", "",
		"key for hash")
	flag.StringVar(&c.CryptoKey, "crypto-key", "",
		"server public key file")
	flag.IntVar(&c.RateLimit, "l", 1,
		"requests limit")
	flag.StringVar(&c.Config, "c", "",
		"path to config")
	flag.Parse()
}
