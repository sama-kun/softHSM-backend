package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type Mailer struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
}

func NewMailer(from, password, host, port string) *Mailer {
	return &Mailer{From: from, Password: password, SMTPHost: host, SMTPPort: port}
}

// SendActivationEmail отправляет красивое письмо с активацией
func (m *Mailer) SendActivationEmail(to, token string) error {
	tmpl, err := template.ParseFiles("templates/activation_email.html")
	if err != nil {
		return fmt.Errorf("ошибка загрузки шаблона: %w", err)
	}

	var body bytes.Buffer
	data := map[string]string{
		"ActivationLink": fmt.Sprintf("http://localhost:8080/activate?token=%s", token),
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("ошибка генерации HTML: %w", err)
	}

	msg := []byte(
		"Subject: Подтверждение регистрации\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
			body.String(),
	)

	auth := smtp.PlainAuth("", m.From, m.Password, m.SMTPHost)
	return smtp.SendMail(m.SMTPHost+":"+m.SMTPPort, auth, m.From, []string{to}, msg)
}
