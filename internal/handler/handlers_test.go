package handler

import (
	"context"
	"dial2verify/internal/config"
	mv "dial2verify/internal/mw"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	// make echo
	e := echo.New()

	// config route
	e.GET("/ping", PingHandler)

	// make HTTP request
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	// call request
	e.ServeHTTP(rec, req)

	// check result
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())
}

type mockStorage struct {
}

func (m *mockStorage) CheckPhone(ctx context.Context, phone string) (bool, error) {
	return false, nil
}

func (m *mockStorage) Close() error {
	return nil
}

func TestCheckPhoneHandler_InvalidFormat(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{Key: "test-api-key"},
	}

	e := echo.New()

	// config routes
	api := e.Group("/api")
	api.Use(mv.APIKeyAuth(cfg.API.Key))
	api.GET("/checkPhone/:phone", CheckPhoneHandler(&mockStorage{}))

	// make api request with invalid phone number
	req := httptest.NewRequest(http.MethodGet, "/api/checkPhone/123456789", nil)
	req.Header.Set("X-API-Key", "test-api-key")
	rec := httptest.NewRecorder()

	// call request
	e.ServeHTTP(rec, req)

	// check result
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		"status":  "error",
		"message": "Invalid phone number format",
	}, resp)

}
