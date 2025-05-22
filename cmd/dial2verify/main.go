package main

import (
	"dial2verify/internal/app"
	"dial2verify/internal/config"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	slog.Info("APP started")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	server, err := app.New(cfg)
	if err != nil {
		slog.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}

	if err := server.Run(); err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}
