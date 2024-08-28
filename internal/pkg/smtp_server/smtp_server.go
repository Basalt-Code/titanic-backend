package smtp_server

import (
	"cmd/app/main.go/internal/config"
	"net/smtp"
)

type SMTPServer struct {
	host     string
	port     string
	email    string
	password string
	auth     smtp.Auth
}

func NewSMTPServer(cfg *config.Config) *SMTPServer {
	auth := smtp.PlainAuth(
		"",
		cfg.SMTPConfig.SenderEmail,
		cfg.SMTPConfig.SenderPassword,
		cfg.SMTPConfig.SmtpHost,
	)
	s := SMTPServer{
		cfg.SMTPConfig.SmtpHost,
		cfg.SMTPConfig.SmtpPort,
		cfg.SMTPConfig.SenderEmail,
		cfg.SMTPConfig.SenderPassword,
		auth,
	}
	return &s
}

func (s *SMTPServer) GetEmail() string {
	return s.email
}

func (s *SMTPServer) SendEmail(receiver_email string, from_email string, subject string, body string) error {
	message := []byte("From: " + s.email + "\n" +
		"Subject: " + subject + "\n\n" +
		body,
	)
	err := smtp.SendMail(
		s.host+":"+s.port,
		s.auth,
		s.email,
		[]string{receiver_email},
		message,
	)
	return err
}
