package utils

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

var e EmailConfig

func InitEmailConfig() {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	e = EmailConfig{
		Host:     os.Getenv("MAIL_SERVER"),
		Port:     port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
	}
}

func sendEmail(to, body, subject string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)
	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}

func SendResetLink(to, resetToken string) error {

	resetLink := fmt.Sprintf("http://192.168.1.10:8080/auth/reset/%s", resetToken)
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>密码重置啦</h2>
				<p>请复制下面的链接到浏览器满足你重置的愿望：</p>
				<p>%s</p>
				<p>如果这不是你发起的请求，请忽略此邮件</p>
				<p>此链接将在15分钟后失效。</p>
			</body>
		</html>
	`, resetLink)

	return sendEmail(to, body, "密码重置邮件")
}

func SendCaptcha(to, captcha string) error {

	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>6位验证码<u>%s</u></h2>
				<p>【gin-test】您的验证码是<strong>%s</strong>, 在20分钟内有效，如非本人操作请忽略此邮件。</p>
			</body>
		</html>
	`, captcha, captcha)

	return sendEmail(to, body, "身份校验邮件")
}
