package config

import (
	"flag"
	"time"
)

type Config struct {
	addr                     string
	maxClients               int
	broadcastPositionTimeout time.Duration
}

func ParseConfig() Config {
	var c Config
	flag.StringVar(&c.addr, "a", "localhost:8888", "address for running server")
	flag.IntVar(&c.maxClients, "m", 2, "max clients for server")
	flag.DurationVar(&c.broadcastPositionTimeout, "d", 50*time.Millisecond, "timeout to broadcast all clients positions")
	flag.Parse()
	return c
}

func (c Config) GetAddress() string {
	return c.addr
}

func (c Config) GetMaxClients() int {
	return c.maxClients
}

func (c Config) GetTimeout() time.Duration {
	return c.broadcastPositionTimeout
}
