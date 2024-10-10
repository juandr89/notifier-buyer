package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Client *redis.ClusterClient
}

func NewNotificationRepository(redisConfig server.RedisConfig) *RedisRepository {
	var redisClusterOptions *redis.ClusterOptions
	if redisConfig.TlsEnable {
		redisClusterOptions = &redis.ClusterOptions{
			Addrs: []string{
				fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			},
			Password: redisConfig.Password,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: redisConfig.TlsEnable,
			},
		}
	} else {
		redisClusterOptions = &redis.ClusterOptions{
			Addrs: []string{
				fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			},
			Password: redisConfig.Password,
		}
	}
	redisClient := redis.NewClusterClient(redisClusterOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Could not connect to Redis: %v", err)
		return nil
	}
	fmt.Printf("Successfully connected to Redis")

	return &RedisRepository{
		Client: redisClient,
	}
}

func (r *RedisRepository) SaveNotification(ctx context.Context, notification domain.Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error when try to map Notification it JSON: %w", err)
	}

	key := fmt.Sprintf("notifications:%s", notification.Email)
	err = r.Client.RPush(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("error while saving Notification in Redis: %w", err)
	}

	return nil
}

func (r *RedisRepository) GetNotifications(ctx context.Context, email string) ([]domain.Notification, error) {
	emailKey := fmt.Sprintf("notifications:%s", email)
	values, err := r.Client.LRange(ctx, emailKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("error al obtener el historial de notificaciones de Redis: %w", err)
	}

	if len(values) < 1 {
		return nil, &domain.NotFoundError{Message: fmt.Sprintf("Notifications with email %s not found", email)}
	}

	notifications := make([]domain.Notification, len(values))
	for i, value := range values {
		var notification domain.Notification
		if err := json.Unmarshal([]byte(value), &notification); err != nil {
			return nil, fmt.Errorf("error decoding notification: %w", err)
		}
		notifications[i] = notification
	}
	return notifications, nil
}

func (r *RedisRepository) GetNotificationCodes(ctx context.Context) ([]string, error) {
	values, err := r.Client.LRange(ctx, "notification:codes", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting notification history from Redis: %w", err)
	}

	if len(values) < 1 {
		return nil, &domain.NotFoundError{Message: "Notification codes not found"}
	}

	return values, nil
}
