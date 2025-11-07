package config

import (
	"os"
)

type Config struct {
	DatabasePath string
	ServerAddr   string
}

func LoadConfig() *Config {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "blog.db" // Default database path
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080" // Default server address
	}

	return &Config{
		DatabasePath: dbPath,
		ServerAddr:   serverAddr,
	}
}
