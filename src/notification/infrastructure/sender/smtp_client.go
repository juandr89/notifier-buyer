package sender

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/juandr89/delivery-notifier-buyer/server"
)

type SmtpClient struct {
	configSMTP server.SMTPConfig
}

func NewSmtpClient(configSMTP server.SMTPConfig) *SmtpClient {
	log.Printf("Loading smtp config host %s port %d", configSMTP.Host, configSMTP.Port)
	return &SmtpClient{
		configSMTP: configSMTP,
	}
}

func (smtpClient *SmtpClient) Send(email string, text string) error {

	from := smtpClient.configSMTP.Username
	password := smtpClient.configSMTP.Password

	to := []string{email}

	smtpHost := smtpClient.configSMTP.Host
	smtpPort := fmt.Sprint(smtpClient.configSMTP.Port)

	log.Printf("Send email host %s  port %s", smtpHost, smtpPort)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	message := []byte("To: " + to[0] + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + "Entrega retrasada por clima" + "\r\n" +
		"\r\n" + text)
	err_sending := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err_sending != nil {
		log.Println(err_sending)
		return err_sending
	}
	log.Println("Email Sent Successfully!")
	return nil
}
