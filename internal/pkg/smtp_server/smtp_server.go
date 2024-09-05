package smtpserver

import (
	"cmd/app/main.go/internal/config"
	"net/smtp"
)

type SMTPServer struct {
	cfg config.SMTPConfig
}

func NewSMTPServer(cfg config.SMTPConfig) *SMTPServer {
	s := SMTPServer{
		cfg: cfg,
	}
	return &s
}

func (s *SMTPServer) SendEmail(
	receiverEmail string,
	subject string,
	body string,
) error {
	auth := smtp.PlainAuth(
		"",
		s.cfg.SenderEmail,
		s.cfg.SenderPassword,
		s.cfg.SmtpHost,
	)

	message := []byte("From: " + s.cfg.SenderEmail + "\n" +
		"Subject: " + subject + "\n\n" +
		body,
	)
	err := smtp.SendMail(
		s.cfg.SmtpHost+":"+s.cfg.SmtpPort,
		auth,
		s.cfg.SenderEmail,
		[]string{receiverEmail},
		message,
	)
	return err
}
