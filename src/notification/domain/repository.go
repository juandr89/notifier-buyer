package domain

import (
	"context"
)

type NotificationRepository interface {
	SaveNotification(ctx context.Context, notification Notification) error
	GetNotifications(ctx context.Context, email string) ([]Notification, error)
	GetNotificationCodes(ctx context.Context) ([]string, error)
}
