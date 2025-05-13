package main

import (
	"dial2verify/internal/api/handler"
	"dial2verify/internal/api/middleware"
	"dial2verify/internal/api/storage"
	"dial2verify/internal/config"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found")
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	slog.Debug("CONFIG", "cong", cfg)

	//api
	e := echo.New()
	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.RequestID())

	redisClient := storage.NewRedisClient(cfg.Redis)
	defer redisClient.Close()

	//auth group
	protected := e.Group("/api")
	protected.Use(middleware.APIKeyAuth(cfg.API.Key))

	protected.GET("/checkPhone/:phone", handler.CheckPhoneHandler(redisClient))

	e.GET("/api/ping", func(c echo.Context) error {
		slog.Info("PING")
		//todo add c.logger("PING")
		c.Logger().Info()
		return c.String(http.StatusOK, "pong")
	})

	addr := ":" + cfg.API.Port
	e.Logger.Fatal(e.Start(addr))
}
