package handler

import (
	api "dial2verify/internal/api/storage"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func CheckPhoneHandler(redisClient *api.RedisClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		phone := c.Param("phone")
		if phone == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status":  "error",
				"message": "Phone number is required",
			})
		}

		ctx := c.Request().Context()
		exists, err := redisClient.CheckPhone(ctx, phone)
		if err != nil {
			slog.Error("Redis error", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status":  "error",
				"message": "Internal server error",
			})
		}
		slog.Info("REDIS resp", "phone", phone, "exists", exists)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success",
			"found":  exists,
			"phone":  phone,
		})
	}
}
