package mw

import (
	"dial2verify/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Auth(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-API-Key")
			if key == "" {
				key = c.QueryParam("api_key")
			}

			if key != apiKey {
				return c.JSON(http.StatusUnauthorized, response.Error("Invalid or missing API key"))
			}

			return next(c)
		}
	}
}
