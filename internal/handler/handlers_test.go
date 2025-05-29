package handler

import (
	"context"
	"dial2verify/internal/config"
	"dial2verify/internal/mw"
	"dial2verify/internal/storage"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mockStorage для модульных тестов
type mockStorage struct {
	checkPhone func(ctx context.Context, phone string) (bool, error)
}

func (m *mockStorage) CheckPhone(ctx context.Context, phone string) (bool, error) {
	return m.checkPhone(ctx, phone)
}

func (m *mockStorage) Close() error {
	return nil
}

func setupTestEcho(t *testing.T, store Storage) *echo.Echo {
	cfg := &config.Config{
		API: config.APIConfig{
			Key: "test-api-key",
		},
	}
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), "logger", slog.Default())
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})

	e.GET("/ping", PingHandler)
	api := e.Group("/api")
	api.Use(mw.APIKeyAuth(cfg.API.Key))
	api.GET("/checkPhone/:phone", CheckPhoneHandler(store))

	return e
}

func TestPingHandler(t *testing.T) {
	e := setupTestEcho(t, nil)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())
}

func TestCheckPhoneHandler(t *testing.T) {
	tests := []struct {
		name           string
		phone          string
		apiKey         string
		storage        Storage
		expectedStatus int
		expectedResp   map[string]interface{}
	}{
		{
			name:   "01 Successful verification - phone exists",
			phone:  "79123456789",
			apiKey: "test-api-key",
			storage: &mockStorage{
				checkPhone: func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedResp: map[string]interface{}{
				"status": "success",
				"found":  true,
				"phone":  "79123456789",
			},
		},
		{
			name:   "02 Phone not found or too old",
			phone:  "79123456789",
			apiKey: "test-api-key",
			storage: &mockStorage{
				checkPhone: func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedResp: map[string]interface{}{
				"status": "success",
				"found":  false,
				"phone":  "79123456789",
			},
		},
		{
			name:   "03 Invalid API key",
			phone:  "79123456789",
			apiKey: "wrong-key",
			storage: &mockStorage{
				checkPhone: func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResp: map[string]interface{}{
				"status":  "error",
				"message": "Invalid or missing API key",
			},
		},
		{
			name:   "04 Missing API key",
			phone:  "79123456789",
			apiKey: "",
			storage: &mockStorage{
				checkPhone: func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResp: map[string]interface{}{
				"status":  "error",
				"message": "Invalid or missing API key",
			},
		},
		{
			name:   "05 Invalid phone format",
			phone:  "123456789",
			apiKey: "test-api-key",
			storage: &mockStorage{
				checkPhone: func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"status":  "error",
				"message": "Invalid phone number format",
			},
		},
		{
			name:   "06 Storage error",
			phone:  "79123456789",
			apiKey: "test-api-key",
			storage: &mockStorage{
				checkPhone: func(ctx context.Context, phone string) (bool, error) {
					return false, errors.New("storage error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp: map[string]interface{}{
				"status":  "error",
				"message": "Internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupTestEcho(t, tt.storage)
			req := httptest.NewRequest(http.MethodGet, "/api/checkPhone/"+tt.phone, nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}

func TestCheckPhoneHandler_Integration(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // Используем отдельную БД для тестов
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		t.Skip("Redis not available, skipping integration tests")
	}

	cfg := config.RedisConfig{Host: "localhost", Port: "6379"}
	store, err := storage.NewRedisStorage(cfg)
	if err != nil {
		t.Fatalf("Failed to init storage: %v", err)
	}
	defer store.Close()

	e := setupTestEcho(t, store)
	ctx := context.Background()

	tests := []struct {
		name           string
		phone          string
		apiKey         string
		redisSetup     func()
		expectedStatus int
		expectedResp   map[string]interface{}
	}{
		{
			name:   "Successful verification - phone exists",
			phone:  "79123456789",
			apiKey: "test-api-key",
			redisSetup: func() {
				redisClient.Set(ctx, "incoming_call_79123456789", time.Now().Unix(), 0)
			},
			expectedStatus: http.StatusOK,
			expectedResp: map[string]interface{}{
				"status": "success",
				"found":  true,
				"phone":  "79123456789",
			},
		},
		{
			name:   "Phone not found",
			phone:  "79123456789",
			apiKey: "test-api-key",
			redisSetup: func() {
				// Ключ не создаётся
			},
			expectedStatus: http.StatusOK,
			expectedResp: map[string]interface{}{
				"status": "success",
				"found":  false,
				"phone":  "79123456789",
			},
		},
		//{
		//	name:   "Phone too old",
		//	phone:  "79123456789",
		//	apiKey: "test-api-key",
		//	redisSetup: func() {
		//		redisClient.Set(ctx, "incoming_call_79123456789", time.Now().Add(-5*time.Minute).Unix(), 0)
		//	},
		//	expectedStatus: http.StatusOK,
		//	expectedResp: map[string]interface{}{
		//		"status": "success",
		//		"found":  false,
		//		"phone":  "79123456789",
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем Redis
			redisClient.FlushDB(ctx)

			// Настраиваем Redis
			if tt.redisSetup != nil {
				tt.redisSetup()
			}

			// Создаём запрос
			req := httptest.NewRequest(http.MethodGet, "/api/checkPhone/"+tt.phone, nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rec := httptest.NewRecorder()

			// Выполняем запрос
			e.ServeHTTP(rec, req)

			// Проверяем статус
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// Проверяем ответ
			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}

	// Очистка
	t.Cleanup(func() {
		redisClient.FlushDB(ctx)
		redisClient.Close()
	})
}
