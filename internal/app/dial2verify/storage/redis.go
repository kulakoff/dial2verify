package storage

import (
	"context"
	"dial2verify/internal/app/dial2verify/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"time"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(cfg config.RedisConfig) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		slog.Error("Redis connection error", "error", err)
		return nil, err
	}

	slog.Debug("Success connected to Redis", "host", cfg.Host, "port", cfg.Port)

	return &RedisStorage{client: client}, nil
}

func (rc *RedisStorage) Close() error {
	return rc.client.Close()
}

func (rc *RedisStorage) CheckPhone(ctx context.Context, phone string) (bool, error) {
	redisKey := "incoming_call_" + phone
	exists, err := rc.client.Exists(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
