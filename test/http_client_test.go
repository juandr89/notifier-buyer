package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/juandr89/delivery-notifier-buyer/app_init"
	"github.com/juandr89/delivery-notifier-buyer/server"
	repository "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/repository"
	"github.com/stretchr/testify/assert"
)

func TestDoRequestWithRetry(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	}))
	opts := server.RequestOptions{
		Method:         "GET",
		URL:            mockServer.URL,
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           nil,
		RequestTimeout: 2 * time.Second,
		MaxRetries:     3,
		RetryDelay:     1 * time.Second,
	}
	t.Run("Success", func(t *testing.T) {
		defer mockServer.Close()

		resp, err := server.DoRequestWithRetry(opts)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("AllRetriesFail", func(t *testing.T) {
		defer mockServer.Close()

		resp, err := server.DoRequestWithRetry(opts)

		assert.Nil(t, resp)
		assert.Error(t, err)

		assert.Contains(t, err.Error(), "request failed after 3 attempts")
	})

}

func TestNewNotificationSender(t *testing.T) {
	t.Run("NewNotificationSenderError", func(t *testing.T) {
		config := server.Config{}
		monkey.Patch(repository.NewNotificationRepository, func(redisConfig server.RedisConfig) *repository.RedisRepository {
			return nil
		})
		defer monkey.Unpatch(repository.NewNotificationRepository)
		response := app_init.NewNotificationRepository(&config)

		assert.Nil(t, response)
	})

}
