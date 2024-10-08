package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juandr89/delivery-notifier-buyer/middleware"
	"github.com/stretchr/testify/assert"
)

func TestApiKeyMiddleware(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		validApiKey := "valid-api-key"
		middleware := middleware.ApiKeyMiddleware(validApiKey)

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Authorized"))
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-api-key", validApiKey)

		rr := httptest.NewRecorder()

		handler := middleware(mockHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Authorized", rr.Body.String())
	})

	t.Run("Unauthorized", func(t *testing.T) {
		validApiKey := "valid-api-key"
		middleware := middleware.ApiKeyMiddleware(validApiKey)

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Authorized"))
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("x-api-key", "invalid-api-key")

		rr := httptest.NewRecorder()

		handler := middleware(mockHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Acceso no autorizado\n", rr.Body.String())
	})

	t.Run("MissingApiKey", func(t *testing.T) {
		validApiKey := "valid-api-key"
		middleware := middleware.ApiKeyMiddleware(validApiKey)

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Authorized"))
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		rr := httptest.NewRecorder()

		handler := middleware(mockHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
