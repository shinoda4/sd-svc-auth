package email

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendWelcomeEmail(to, username string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your_email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Account verified!")
	m.SetBody("text/html", "Dear <b>"+username+"</b>, you are already verified! Welcome to our system!")

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

	// 配置 SMTP 客户端信息
	d := gomail.NewDialer("smtp.gmail.com", 587, emailAddress, emailPassword)

	// 发送邮件
	return d.DialAndSend(m)
}

func SendVerifyEmail(to, username, token, verifyLink string) error {
	// 使用传入的 verifyLink 拼接 token
	fullLink := fmt.Sprintf("%s?token=%s", verifyLink, token)

	m := gomail.NewMessage()
	m.SetHeader("From", "your_email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify your email!")
	m.SetBody("text/html", fmt.Sprintf(
		"Dear <b>%s</b>, please finish your account validation by clicking the following link: <a href='%s'>Verify Email</a>",
		username, fullLink,
	))
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
