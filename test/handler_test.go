package service_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure"
	third_party "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/third_party"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/usecases"
	mocks "github.com/juandr89/delivery-notifier-buyer/test/mocks_test"
	"github.com/stretchr/testify/assert"
)

func TestListBuyerNotifications(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		email := "test@example.com"

		expectedResponse := usecases.NotificationHistoryServiceResponse{
			History: []usecases.NotificationHistoryDetail{
				{
					NotificationSendAt: time.Now(),
					Location: usecases.Location{
						Latitude:  "40.7128",
						Longitude: "-74.0060",
					},
					ForecastCode: 1234,
				},
			},
		}

		monkey.Patch(usecases.GetBuyerNotification, func(emailparam string, repo domain.NotificationRepository) (*usecases.NotificationHistoryServiceResponse, error) {
			return &expectedResponse, nil
		})
		defer monkey.Unpatch(usecases.GetBuyerNotification)
		handler := infrastructure.NotificationHandler{
			NotificationRepository: mockRepo,
		}

		//body, _ := json.Marshal(requestGetNotification)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notifications/%s", email), nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.BuyerNotifications(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		var responseBody usecases.NotificationHistoryServiceResponse
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		//assert.Equal(t, expectedResponse, responseBody)
	})
	t.Run("NotFoundBuyerNotification", func(t *testing.T) {

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		handler := infrastructure.NewNotificationHandler(mockRepo, nil, server.Config{})

		email := "buyer@example.com"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notifications/%s", email), nil)
		rr := httptest.NewRecorder()

		monkey.Patch(usecases.GetBuyerNotification, func(emailparam string, repo domain.NotificationRepository) (*usecases.NotificationHistoryServiceResponse, error) {
			return nil, &domain.NotFoundError{Message: fmt.Sprintf("Notifications with email %s not found", email)}
		})
		defer monkey.Unpatch(usecases.GetBuyerNotification)

		handler.BuyerNotifications(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Notifications with email buyer@example.com not found")
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		handler := infrastructure.NewNotificationHandler(mockRepo, nil, server.Config{})

		email := "buyer@example.com"
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notifications/%s", email), nil)

		rr := httptest.NewRecorder()

		monkey.Patch(usecases.GetBuyerNotification, func(emailparam string, repo domain.NotificationRepository) (*usecases.NotificationHistoryServiceResponse, error) {
			return nil, errors.New("Unexpected error has ocurred")
		})
		defer monkey.Unpatch(usecases.GetBuyerNotification)

		handler.BuyerNotifications(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unexpected error has ocurred")

	})
}

func TestNotifyBuyer(t *testing.T) {
	t.Run("InvalidJSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		handler := &infrastructure.NotificationHandler{
			Config: server.Config{},
		}

		invalidJSON := `{"invalid": }`
		req := httptest.NewRequest("POST", "/notifications", bytes.NewBufferString(invalidJSON))
		w := httptest.NewRecorder()

		handler.NotifyBuyer(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid JSON data")
	})

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		mockSender := mocks.NewMockNotificationSender(ctrl)

		requestData := usecases.RequestDataNotification{
			Email: "test@example.com",
			Location: usecases.Location{
				Latitude:  "40.7128",
				Longitude: "-74.0060",
			},
		}
		requestBody, _ := json.Marshal(requestData)

		var mockNotificationResponse = usecases.NotificationServiceResponse{
			ForecastCode:        123,
			ForecastDescription: "Mocked forecast description",
			BuyerNotification:   true,
		}

		monkey.Patch(usecases.SendNotification, func(req usecases.RequestDataNotification, forecastService third_party.IForecastService, repo domain.NotificationRepository, sender domain.NotificationSender) (*usecases.NotificationServiceResponse, error) {
			return &mockNotificationResponse, nil
		})
		defer monkey.Unpatch(usecases.SendNotification)
		handler := &infrastructure.NotificationHandler{
			NotificationRepository: mockRepo,
			NotificationSender:     mockSender,
			Config:                 server.Config{},
		}

		req := httptest.NewRequest("POST", "/notifications", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		handler.NotifyBuyer(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"forecast_code":123`)
	})

	t.Run("NewNotificationHandlerSuccess", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		mockSender := mocks.NewMockNotificationSender(ctrl)

		cfg := server.Config{
			SMTPConfig: server.SMTPConfig{},
			RedisConfig: server.RedisConfig{
				Host: "localhost",
				Port: 6379,
			},
		}

		handler := infrastructure.NewNotificationHandler(mockRepo, mockSender, cfg)

		assert.NotNil(t, handler)
		assert.Equal(t, mockRepo, handler.NotificationRepository)
		assert.Equal(t, mockSender, handler.NotificationSender)
		assert.Equal(t, cfg, handler.Config)
	})

	t.Run("SendNotificationError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		requestData := usecases.RequestDataNotification{
			Email: "test@example.com",
			Location: usecases.Location{
				Latitude:  "40.7128",
				Longitude: "-74.0060",
			},
		}
		requestBody, _ := json.Marshal(requestData)

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		mockSender := mocks.NewMockNotificationSender(ctrl)

		cfg := server.Config{
			SMTPConfig: server.SMTPConfig{},
			RedisConfig: server.RedisConfig{
				Host: "localhost",
				Port: 6379,
			},
		}
		monkey.Patch(usecases.SendNotification, func(requestDataNotification usecases.RequestDataNotification, forecastService third_party.IForecastService, repository domain.NotificationRepository, sender domain.NotificationSender) (*usecases.NotificationServiceResponse, error) {
			return nil, errors.New("failed to send notification")
		})
		defer monkey.Unpatch(usecases.SendNotification)
		handler := infrastructure.NewNotificationHandler(mockRepo, mockSender, cfg)

		req := httptest.NewRequest("POST", "/notifications", bytes.NewReader(requestBody))
		w := httptest.NewRecorder()

		handler.NotifyBuyer(w, req)

		assert.Contains(t, w.Body.String(), "failed to send notification")

	})
}
