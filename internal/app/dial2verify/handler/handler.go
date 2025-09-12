package handler

import (
	"dial2verify/internal/app/dial2verify/storage"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"regexp"
)

type Handler struct {
	s storage.Storage
	l *slog.Logger
}

func New(s storage.Storage, l *slog.Logger) *Handler {
	return &Handler{s: s, l: l}
}

func (h *Handler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (h *Handler) Check(c echo.Context) error {
	phone := c.Param("phone")
	if !regexp.MustCompile(`^7[0-9]{10}$`).MatchString(phone) {
		h.l.Debug("Invalid phone number format", "phone", phone)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status":  "error",
			"message": "Invalid phone number format",
		})
	}

	exists, err := h.s.CheckPhone(c.Request().Context(), phone)
	if err != nil {
		h.l.Error("Storage error", "phone", phone, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"found":  exists,
		"phone":  phone,
	})
}
