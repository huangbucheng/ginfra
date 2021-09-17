package utils

import (
	"net/smtp"
)

func SendMail(from, password string, smtpHost, smtpPort string, message []byte, to []string) error {
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}
