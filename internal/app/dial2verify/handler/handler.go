package handler

import (
	"dial2verify/internal/app/dial2verify/storage"
	"dial2verify/pkg/response"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

var phonePattern = regexp.MustCompile(`^7[0-9]{10}$`)

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
	start := time.Now()
	checkPhoneRequestsTotal.Inc()
	defer checkPhoneDurationSeconds.Observe(time.Since(start).Seconds())

	phone := c.Param("phone")
	if !phonePattern.MatchString(phone) {
		checkPhoneInvalidTotal.Inc()
		h.l.Debug("Invalid phone number format", "phone", phone)
		return c.JSON(http.StatusBadRequest,
			response.Error("Invalid phone number format"))
	}

	exists, err := h.s.CheckPhone(c.Request().Context(), phone)
	if err != nil {
		checkPhoneErrorsTotal.Inc()
		h.l.Error("Storage error", "phone", phone, "error", err)
		return c.JSON(http.StatusInternalServerError,
			response.Error("Internal server error"))
	}

	if exists {
		checkPhoneFoundTrueTotal.Inc()
	} else {
		checkPhoneFoundFalseTotal.Inc()
	}

	return c.JSON(http.StatusOK, response.SuccessPhoneCheck(phone, exists))
}
