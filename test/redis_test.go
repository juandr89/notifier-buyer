package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	repository "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/repository"
	"github.com/stretchr/testify/assert"
)

func TestSaveNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ConnectionTlsDisable", func(t *testing.T) {
		redisConfig := server.RedisConfig{Host: "localhost", Port: 6379}
		_, mock := redismock.NewClientMock()

		mock.ExpectPing().SetErr(errors.New("connection error"))

		repo := repository.NewNotificationRepository(redisConfig)

		assert.NotNil(t, repo)
	})
	t.Run("ConnectionTlsEnable", func(t *testing.T) {
		redisConfig := server.RedisConfig{Host: "localhost", Port: 6379, TlsEnable: true}
		_, mock := redismock.NewClientMock()

		mock.ExpectPing().SetErr(errors.New("connection error"))

		repo := repository.NewNotificationRepository(redisConfig)

		assert.NotNil(t, repo)
	})

	t.Run("Success", func(t *testing.T) {
		redisMock, mock := redismock.NewClientMock()
		repo := &repository.RedisRepository{Client: redisMock}

		ctx := context.Background()
		notification := domain.Notification{Email: "test@example.com"}

		data, _ := json.Marshal(notification)
		key := fmt.Sprintf("notifications:%s", notification.Email)

		mock.ExpectRPush(key, data).SetVal(1)

		err := repo.SaveNotification(ctx, notification)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error", func(t *testing.T) {
		redisMock, mock := redismock.NewClientMock()
		repo := &repository.RedisRepository{Client: redisMock}

		ctx := context.Background()
		notification := domain.Notification{Email: "test@example.com"}

		data, _ := json.Marshal(notification)
		key := fmt.Sprintf("notifications:%s", notification.Email)

		mock.ExpectRPush(key, data).SetErr(errors.New("redis error"))

		err := repo.SaveNotification(ctx, notification)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetNotifications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("Success", func(t *testing.T) {

		redisMock, mock := redismock.NewClientMock()
		repo := &repository.RedisRepository{Client: redisMock}

		ctx := context.Background()
		email := "test@example.com"
		notification := domain.Notification{Email: email}

		data, _ := json.Marshal(notification)
		key := fmt.Sprintf("notifications:%s", email)

		mock.ExpectLRange(key, 0, -1).SetVal([]string{string(data)})

		notifications, err := repo.GetNotifications(ctx, email)

		assert.NoError(t, err)
		assert.Len(t, notifications, 1)
		assert.Equal(t, notification, notifications[0])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		redisMock, mock := redismock.NewClientMock()
		repo := &repository.RedisRepository{Client: redisMock}

		ctx := context.Background()
		email := "nonexistent@example.com"

		key := fmt.Sprintf("notifications:%s", email)
		mock.ExpectLRange(key, 0, -1).SetVal([]string{})

		notifications, err := repo.GetNotifications(ctx, email)

		assert.Error(t, err)
		assert.Nil(t, notifications)
		assert.IsType(t, &domain.NotFoundError{}, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetNotificationCodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		redisMock, mock := redismock.NewClientMock()
		repo := &repository.RedisRepository{Client: redisMock}

		ctx := context.Background()

		mock.ExpectLRange("notification:codes", 0, -1).SetVal([]string{"code1", "code2"})

		codes, err := repo.GetNotificationCodes(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []string{"code1", "code2"}, codes)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		redisMock, mock := redismock.NewClientMock()
		repo := &repository.RedisRepository{Client: redisMock}

		ctx := context.Background()

		mock.ExpectLRange("notification:codes", 0, -1).SetVal([]string{})

		codes, err := repo.GetNotificationCodes(ctx)

		assert.Error(t, err)
		assert.Nil(t, codes)
		assert.IsType(t, &domain.NotFoundError{}, err)
		assert.NoError(t, mock.ExpectationsWereMet())

	})
}
