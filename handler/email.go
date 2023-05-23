package handler

import (
	"errors"
	"net/smtp"
	"os"
)

func SendEmail(msg string, to string) error {
  from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, pass, host)

	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return errors.New("error sending email")
	}
	return nil
}
