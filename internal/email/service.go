package email

import (
	"gopkg.in/gomail.v2"
)

func SendWelcomeEmail(to, username string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your_email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "欢迎注册我们的服务！")
	m.SetBody("text/html", "亲爱的 <b>"+username+"</b>，欢迎加入！")

	// 配置 SMTP 客户端信息
	d := gomail.NewDialer("smtp.gmail.com", 587, "lindesong666@gmail.com", "pafdqhupcsprtlue")

	// 发送邮件
	return d.DialAndSend(m)
}
