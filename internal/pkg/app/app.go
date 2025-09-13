package app

import (
	"dial2verify/internal/app/dial2verify/config"
	"dial2verify/internal/app/dial2verify/handler"
	"dial2verify/internal/app/dial2verify/mw"
	storage2 "dial2verify/internal/app/dial2verify/storage"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"log/slog"
	"os"
)

type App struct {
	echo    *echo.Echo
	cfg     *config.Config
	store   storage2.Storage
	logger  *slog.Logger
	handler *handler.Handler
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// init Echo
	e := echo.New()
	e.Use(echoMW.Recover())
	e.Use(echoMW.RequestID())
	e.Use(echoMW.Logger())

	// init storage
	store, err := storage2.NewRedisStorage(cfg.Redis)
	if err != nil {
		logger.Error("Failed to init storage", "error", err)
		os.Exit(1)
	}

	// init handler
	h := handler.New(store, logger)

	e.GET("/api/health", h.Ping)
	e.GET("/api/checkPhone/:phone", h.Check, mw.Auth(cfg.API.Key))

	return &App{
		echo:   e,
		cfg:    cfg,
		store:  store,
		logger: logger,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("Server started", "port", a.cfg.API.Port)
	defer a.store.Close()

	err := a.echo.Start(":" + a.cfg.API.Port)
	if err != nil {
		a.logger.Warn("Error running server", "error", err)
	}
	return nil
}
