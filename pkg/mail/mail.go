package mail

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/spf13/viper"
)

// SendEmail отправляет email через MailHog
func SendEmail(to, subject, body string) error {
	smtpHost := viper.GetString("smtp.host")
	smtpPort := viper.GetString("smtp.port")
	smtpUser := viper.GetString("smtp.username")
	smtpPass := viper.GetString("smtp.password")
	sender := viper.GetString("smtp.sender") // Исправлено: раньше использовался smtp.from

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Формируем email
	msg := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"Subject: " + subject + "\r\n" +
		"From: " + sender + "\r\n" +
		"To: " + to + "\r\n\r\n" +
		body)

	// Отправляем email
	err := smtp.SendMail(addr, auth, sender, []string{to}, msg)
	if err != nil {
		log.Printf("Ошибка при отправке email: %v", err)
		return err
	}

	log.Printf("Email отправлен на %s", to)
	return nil
}
