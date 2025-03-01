package app

import (
	"flag"
	"net"
	"os"
)

type Config struct {
	AddrServer    string
	AddrDatabase  string
	LogLevel      string
	StoreInterval uint64
	FilenameSave  string
	RestoreData   bool
	SecretKey     string
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.AddrServer, "a", "localhost:8080", "server and port to run server")
	// flag.StringVar(&cfg.AddrDatabase, "d", "", "address to postgres base")
	flag.StringVar(&cfg.AddrDatabase, "d", "host=localhost user=pUser password=temp_pass dbname=db_temp sslmode=disable", "address to postgres base")
	flag.StringVar(&cfg.LogLevel, "r", "info", "adrees accrual system `ip:port`")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.AddrServer = envRunAddr
	}

	if envDatabaseAddr := os.Getenv("DATABASE_DSN"); envDatabaseAddr != "" {
		cfg.AddrDatabase = envDatabaseAddr
	}

	if envSecretKey := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envSecretKey != "" {
		cfg.SecretKey = envSecretKey
	}

	flag.Parse()

	_, _, err := net.SplitHostPort(cfg.AddrServer)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
