package config

import (
	"flag"
	"net"
	"os"
	"time"
)

type Config struct {
	AddrServer        string
	AddrDatabase      string
	AddrAccraulSystem string
	SecretKey         string
	AliveToken        time.Duration
}

var config *Config = nil

func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	config.AliveToken = time.Hour * 24
	return config
}

func ParseFlags() (*Config, error) {
	cfg := Get()
	flag.StringVar(&cfg.AddrServer, "a", "localhost:8080", "server and port to run server")
	// flag.StringVar(&cfg.AddrDatabase, "d", "", "address to postgres base")
	flag.StringVar(&cfg.AddrDatabase, "d", "host=localhost user=pUser password=temp_pass dbname=db_temp sslmode=disable", "address to postgres base")
	flag.StringVar(&cfg.AddrAccraulSystem, "r", "localhost:8000", "adrees accrual system `ip:port`")
	flag.StringVar(&cfg.SecretKey, "s", "secret_key", "secret key for jwt")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.AddrServer = envRunAddr
	}

	if envDatabaseAddr := os.Getenv("DATABASE_DSN"); envDatabaseAddr != "" {
		cfg.AddrDatabase = envDatabaseAddr
	}

	if envAddressAccraulSystem := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAddressAccraulSystem != "" {
		cfg.AddrAccraulSystem = envAddressAccraulSystem
	}

	if envDatabaseAddr := os.Getenv("DATABASE_DSN"); envDatabaseAddr != "" {
		cfg.AddrDatabase = envDatabaseAddr
	}

	flag.Parse()

	_, _, err := net.SplitHostPort(cfg.AddrServer)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
