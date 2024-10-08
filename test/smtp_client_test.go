package service_test

import (
	"errors"
	"net/smtp"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/sender"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/usecases"
	mocks "github.com/juandr89/delivery-notifier-buyer/test/mocks_test"
	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("Success", func(t *testing.T) {
		smtpClient := sender.NewSmtpClient(server.SMTPConfig{
			Username: "testuser",
			Password: "testpass",
		})

		email := "recipient@example.com"
		text := "This is a test email."

		monkey.Patch(smtp.SendMail, func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			return nil
		})

		defer monkey.Unpatch(usecases.SendNotification)
		err := smtpClient.Send(email, text)

		assert.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockSmtpClient := mocks.NewMockSmtpClient(ctrl)

		email := "recipient@example.com"
		text := "This is a test email."

		mockSmtpClient.EXPECT().Send(email, text).Return(errors.New("failed to authenticate"))

		defer monkey.Unpatch(usecases.SendNotification)
		err := mockSmtpClient.Send(email, text)

		assert.EqualError(t, err, "failed to authenticate")
	})

	t.Run("NotFoundError", func(t *testing.T) {

		expected_message := "Not Found: Resource not found"
		err := &domain.NotFoundError{Message: "Resource not found"}

		if err.Error() != expected_message {
			t.Errorf("expected %q, got %q", expected_message, err.Error())
		}
	})
}
