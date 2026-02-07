package config

import "os"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	AccessSecret      string
	RefreshSecret     string
	AccessExpiration  string
	RefreshExpiration string
}

func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ":8080"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			AccessSecret:      getEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret:     getEnv("JWT_REFRESH_SECRET", ""),
			AccessExpiration:  getEnv("JWT_ACCESS_EXPIRATION", ""),
			RefreshExpiration: getEnv("JWT_REFRESH_EXPIRATION", ""),
		},
	}
}

func getEnv(key string, defVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defVal
}
