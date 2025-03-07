package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"soft-hsm/internal/config"
	"strconv"
)

type Mailer struct {
	mailerConfig *config.MailerConfig
}

func NewMailer(mailerConfig *config.MailerConfig) *Mailer {
	return &Mailer{mailerConfig: mailerConfig}
}

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

	auth := smtp.PlainAuth("", m.mailerConfig.From, m.mailerConfig.Password, m.mailerConfig.SMTPHost)

	return smtp.SendMail(
		m.mailerConfig.SMTPHost+":"+strconv.Itoa(m.mailerConfig.SMTPPort),
		auth,
		m.mailerConfig.From,
		[]string{to},
		msg,
	)
}
