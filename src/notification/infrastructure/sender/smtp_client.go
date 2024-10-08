package sender

import (
	"fmt"
	"net/smtp"

	"github.com/juandr89/delivery-notifier-buyer/server"
)

type SmtpClient struct {
	configSMTP server.SMTPConfig
}

func NewSmtpClient(configSMTP server.SMTPConfig) *SmtpClient {
	return &SmtpClient{
		configSMTP: configSMTP,
	}
}

func (smtpClient *SmtpClient) Send(email string, text string) error {

	from := smtpClient.configSMTP.Username
	password := smtpClient.configSMTP.Password

	to := []string{"juandruiz101@gmail.com"}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Subject: Test Email\n\nThis is the email body.")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}
