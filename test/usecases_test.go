package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	third_party "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/third_party"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/usecases"
	mocks "github.com/juandr89/delivery-notifier-buyer/test/mocks_test"
	"github.com/stretchr/testify/assert"
)

func TestRequireBuyerNotification(t *testing.T) {

	t.Run("NotificationCodeMatches", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockNotificationRepository(ctrl)

		mockCtx := context.TODO()
		code := 1234.56
		notificationCodes := []string{strconv.FormatFloat(code, 'f', -1, 64)}

		mockRepo.EXPECT().GetNotificationCodes(gomock.Any()).Return(notificationCodes, nil).Times(1)
		result, err := usecases.RequireBuyerNotification(mockCtx, mockRepo, code)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, *result)

	})

	t.Run("NotificationCodeDoesNotMatch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		mockCtx := context.TODO()

		code := 1234.56
		notificationCodes := []string{"9999.99", "8888.88"}
		mockRepo.EXPECT().GetNotificationCodes(gomock.Any()).Return(notificationCodes, nil).Times(1)

		result, err := usecases.RequireBuyerNotification(mockCtx, mockRepo, code)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, *result)
	})

	t.Run(" ErrorRetrievingNotificationCodes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		mockCtx := context.TODO()

		code := 1234.56
		repoError := errors.New("database error")

		mockRepo.EXPECT().GetNotificationCodes(gomock.Any()).Return(nil, repoError).Times(1)

		result, err := usecases.RequireBuyerNotification(mockCtx, mockRepo, code)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, repoError, err)
	})
}

func TestSendNotification(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		mockSender := mocks.NewMockNotificationSender(ctrl)
		mockForecastService := mocks.NewMockForecastService(ctrl)

		requestData := &usecases.RequestDataNotification{
			Email:    "test@example.com",
			Location: usecases.Location{Latitude: "40.7128", Longitude: "-74.0060"},
		}

		expectedForecast := &third_party.ForecastServiceResponse{
			Code:        123,
			Description: "Sunny",
		}

		mockForecastService.EXPECT().FetchForecastByLocation(
			requestData.Location.Longitude, requestData.Location.Latitude, "2",
		).Return(expectedForecast, nil).Times(1)

		mockRepo.EXPECT().GetNotificationCodes(gomock.Any()).Return([]string{"123"}, nil).Times(1)
		mockRepo.EXPECT().SaveNotification(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		mockSender.EXPECT().Send(requestData.Email, gomock.Any()).Times(1)

		response, err := usecases.SendNotification(requestData, mockForecastService, mockRepo, mockSender)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, float64(123), response.ForecastCode)
		assert.Equal(t, "Sunny", response.ForecastDescription)
	})

	t.Run("ForecastError", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		mockSender := mocks.NewMockNotificationSender(ctrl)
		mockForecastService := mocks.NewMockForecastService(ctrl)

		requestData := &usecases.RequestDataNotification{
			Email:    "test@example.com",
			Location: usecases.Location{Latitude: "40.7128", Longitude: "-74.0060"},
		}

		mockForecastService.EXPECT().FetchForecastByLocation(
			requestData.Location.Longitude, requestData.Location.Latitude, "2",
		).Return(nil, errors.New("failed to fetch forecast")).Times(1)

		response, err := usecases.SendNotification(requestData, mockForecastService, mockRepo, mockSender)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "failed to fetch forecast")

	})

	t.Run("SaveNotificationError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		mockSender := mocks.NewMockNotificationSender(ctrl)
		mockForecastService := mocks.NewMockForecastService(ctrl)

		requestData := &usecases.RequestDataNotification{
			Email:    "test@example.com",
			Location: usecases.Location{Latitude: "40.7128", Longitude: "-74.0060"},
		}

		expectedForecast := &third_party.ForecastServiceResponse{
			Code:        123,
			Description: "Sunny",
		}

		mockForecastService.EXPECT().FetchForecastByLocation(
			requestData.Location.Longitude, requestData.Location.Latitude, "2",
		).Return(expectedForecast, nil).Times(1)

		mockRepo.EXPECT().GetNotificationCodes(gomock.Any()).Return([]string{"123"}, nil).Times(1)
		mockRepo.EXPECT().SaveNotification(gomock.Any(), gomock.Any()).Return(errors.New("failed to save notification")).Times(1)
		mockSender.EXPECT().Send(requestData.Email, gomock.Any()).Times(1)

		response, err := usecases.SendNotification(requestData, mockForecastService, mockRepo, mockSender)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.EqualError(t, err, "failed to save notification")
	})
}

func TestGetBuyerNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("Success", func(t *testing.T) {

		mockRepo := mocks.NewMockNotificationRepository(ctrl)

		// requestGetNotification := &usecases.RequestGetNotification{
		// 	Email: "test@example.com",
		// }
		email := "test@example.com"
		notifications := []domain.Notification{
			{
				Email: "test@example.com",
				DeliveryLocation: domain.DeliveryLocation{
					Longitude: "-15.6345",
					Latitude:  "-15.6345",
				},
				ForecastCode: 12345,
			},
			{
				Email: "Notification 2",
				DeliveryLocation: domain.DeliveryLocation{
					Longitude: "-15.6345",
					Latitude:  "-10.6345",
				},
				ForecastCode: 15450,
			},
		}

		mockRepo.EXPECT().
			GetNotifications(gomock.Any(), email).
			Return(notifications, nil).
			Times(1)

		result, err := usecases.GetBuyerNotification(email, mockRepo)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		assert.Equal(t, len(notifications), len(result.History))
		assert.Equal(t, float64(12345), result.History[0].ForecastCode)
		assert.Equal(t, float64(15450), result.History[1].ForecastCode)
	})
	t.Run("Error", func(t *testing.T) {
		mockRepo := mocks.NewMockNotificationRepository(ctrl)
		email := "test@example.com"

		mockRepo.EXPECT().
			GetNotifications(gomock.Any(), email).
			Return(nil, errors.New("error getting notifications")).
			Times(1)

		result, err := usecases.GetBuyerNotification(email, mockRepo)

		assert.Nil(t, result)
		assert.EqualError(t, err, "error getting notifications")
	})

}

func TestNotificationEntityToDTO_Success(t *testing.T) {
	sampleNotification := domain.Notification{
		Email: "test@example.com",
		DeliveryLocation: domain.DeliveryLocation{
			Longitude: "40.7128",
			Latitude:  "74.0060",
		},
		ForecastCode: 12345,
		Created_at:   time.Now(),
	}

	result := usecases.NotificationEntityToDTO(sampleNotification)

	assert.Equal(t, sampleNotification.Created_at, result.NotificationSendAt, "Expected the Created_at value to be correctly mapped")
	assert.Equal(t, sampleNotification.DeliveryLocation.Latitude, result.Location.Latitude, "Expected the Latitude to be correctly mapped")
	assert.Equal(t, sampleNotification.DeliveryLocation.Longitude, result.Location.Longitude, "Expected the Longitude to be correctly mapped")
	assert.Equal(t, sampleNotification.ForecastCode, result.ForecastCode, "Expected the Codigo to be correctly mapped to ForecastCode")
}
