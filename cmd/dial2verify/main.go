package main

import (
	"dial2verify/internal/app/dial2verify/config"
	"dial2verify/internal/pkg/app"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	server, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}

	if err := server.Run(); err != nil {
		logger.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}
