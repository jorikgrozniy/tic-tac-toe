package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

func New() *Config {
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ":8080"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://user:password@host:port/database"),
		},
	}
}

func getEnv(key string, defVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defVal
}
