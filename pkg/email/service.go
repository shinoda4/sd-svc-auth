package email

import (
	"errors"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(from string, to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	if emailPassword == "" {
		err := errors.New("EMAIL_PASSWORD environment variable not set")
		return err
	}
	emailAddress := os.Getenv("EMAIL_ADDRESS")
	if emailAddress == "" {
		err := errors.New("EMAIL_ADDRESS environment variable not set")
		return err
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, emailAddress, emailPassword)
	return d.DialAndSend(m)
}
