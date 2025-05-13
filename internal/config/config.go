package config

import (
	"os"
)

type Config struct {
	Redis RedisConfig
	API   ApiConfig
}
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}
type ApiConfig struct {
	Port string
	Key  string
}

func Load() (Config, error) {

	redisCfg := RedisConfig{
		Host:     getEnv("REDIS_HOST", "storage"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
	}
	apiCfg := ApiConfig{
		Port: getEnv("APIPort", "8080"),
		Key:  getEnv("APIKey", "your-secret-api-key"),
	}
	cfg := Config{redisCfg, apiCfg}
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultValue
}
