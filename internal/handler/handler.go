package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"regexp"
)

type Storage interface {
	CheckPhone(ctx context.Context, phone string) (bool, error)
	Close() error
}

func CheckPhoneHandler(store Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger, ok := c.Request().Context().Value("logger").(*slog.Logger)
		if !ok {
			logger = slog.Default()
		}

		phone := c.Param("phone")
		if !regexp.MustCompile(`^7[0-9]{10}$`).MatchString(phone) {
			logger.Debug("Invalid phone number format", "phone", phone)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status":  "error",
				"message": "Invalid phone number format",
			})
		}

		ctx := c.Request().Context()
		exists, err := store.CheckPhone(ctx, phone)
		if err != nil {
			logger.Error("Storage error", "phone", phone, "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status":  "error",
				"message": "Internal server error",
			})
		}

		logger.Debug("Storage response", "phone", phone, "exists", exists)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success",
			"found":  exists,
			"phone":  phone,
		})
	}
}

func PingHandler(c echo.Context) error {
	logger, ok := c.Request().Context().Value("logger").(*slog.Logger)
	if !ok {
		logger = slog.Default()
	}
	logger.Debug("PING")
	return c.String(http.StatusOK, "pong")
}
