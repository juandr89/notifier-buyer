package domain

type NotificationSender interface {
	Send(email string, text string) error
}
