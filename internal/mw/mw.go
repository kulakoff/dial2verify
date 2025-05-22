package mw

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func APIKeyAuth(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-API-Key")
			if key == "" {
				key = c.QueryParam("api_key")
			}

			if key != apiKey {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"status":  "error",
					"message": "Invalid or missing API key",
				})
			}

			return next(c)
		}
	}
}
