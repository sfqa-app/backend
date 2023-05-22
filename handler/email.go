package handler

import (
	"errors"
	"log"
	"net/smtp"
	"os"
)

var from, pass, host, port string

func init() {
	from = os.Getenv("SMTP_EMAIL")
	pass = os.Getenv("SMTP_PASSWORD")
	host = os.Getenv("SMTP_HOST")
	port = os.Getenv("SMTP_PORT")
}

func SendEmail(msg string, to string) error {
	auth := smtp.PlainAuth("", from, pass, host)

	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Println(err)
		return errors.New("error sending email")
	}
	return nil
}
