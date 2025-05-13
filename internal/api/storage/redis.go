package storage

import (
	"context"
	"dial2verify/internal/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg config.RedisConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		slog.Warn("Redis connection error", "error", err)
	}

	return &RedisClient{client: client}
}

func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

func (rc *RedisClient) CheckPhone(ctx context.Context, phone string) (bool, error) {
	redisKey := "incoming_call_" + phone
	exists, err := rc.client.Exists(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
