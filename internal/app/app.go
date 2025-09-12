package app

import (
	"dial2verify/internal/config"
	"dial2verify/internal/handler"
	"dial2verify/internal/mw"
	"dial2verify/internal/storage"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"log/slog"
	"os"
)

type App struct {
	echo  *echo.Echo
	cfg   *config.Config
	store storage.Storage
}

func New(cfg *config.Config) (*App, error) {
	// init Echo
	e := echo.New()
	e.Use(echoMW.Recover())
	e.Use(echoMW.RequestID())
	e.Use(echoMW.Logger())

	// init storage
	store, err := storage.NewRedisStorage(cfg.Redis)
	if err != nil {
		slog.Error("Failed to init storage", "error", err)
		os.Exit(1)
	}

	e.GET("/ping", handler.PingHandler)

	group := e.Group("/api")
	group.Use(mw.APIKeyAuth(cfg.API.Key))
	group.GET("/checkPhone/:phone", handler.CheckPhoneHandler(store))

	return &App{
		echo:  e,
		cfg:   cfg,
		store: store,
	}, nil
}

func (a *App) Run() error {
	slog.Info("Server started", "port", a.cfg.API.Port)
	defer a.store.Close()

	err := a.echo.Start(":" + a.cfg.API.Port)
	if err != nil {
		slog.Warn("Error running server", "error", err)
	}
	return nil
}
