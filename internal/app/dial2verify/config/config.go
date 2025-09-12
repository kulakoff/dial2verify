package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Config struct {
	Redis RedisConfig
	API   APIConfig
}
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}
type APIConfig struct {
	Port string
	Key  string
}

func Load() (*Config, error) {
	slog.Debug("Load config start")
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	cfg := Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "storage"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		API: APIConfig{
			Port: getEnv("API_PORT", "8080"),
			Key:  getEnv("API_KEY", ""),
		},
	}

	if cfg.API.Key == "" {
		return &Config{}, fmt.Errorf("API_KEY is required")
	}

	return &cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
