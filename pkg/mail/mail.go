package mail

import (
	"fmt"
	"jumyste-app-backend/config"
	"log"
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	smtpConfig := config.AppConfig.SMTP

	log.Printf("SMTP_HOST: %s, SMTP_PORT: %s, SMTP_SENDER: %s", smtpConfig.Host, smtpConfig.Port, smtpConfig.Sender)

	addr := fmt.Sprintf("%s:%s", smtpConfig.Host, smtpConfig.Port)

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	msg := []byte("From: " + smtpConfig.Sender + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		body)

	err := smtp.SendMail(addr, auth, smtpConfig.Sender, []string{to}, msg)
	if err != nil {
		log.Printf("Ошибка при отправке email: %v", err)
		return err
	}

	log.Printf("Email отправлен на %s", to)
	return nil
}
