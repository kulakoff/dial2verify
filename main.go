package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var rdb *redis.Client

func initRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPasswd := os.Getenv("REDIS_PASSWORD")

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPasswd,
		DB:       0,
	})

	// check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		slog.Error("Redis connection error: %v", err)
	}

}

func checkPhone(c echo.Context) error {
	phone := c.Param("phone")
	ctx := c.Request().Context()

	exists, err := rdb.Exists(ctx, "incoming_call_"+phone).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if exists == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"found": exists == 1,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"found": exists == 1,
	})
}

func main() {
	initRedis()

	//api
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/checkPhone/:phone", checkPhone)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(addr))
}
