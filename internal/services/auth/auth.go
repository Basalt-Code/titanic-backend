package auth

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"cmd/app/main.go/internal/config"
	"cmd/app/main.go/internal/model"
	"cmd/app/main.go/internal/pkg/logger"
	"cmd/app/main.go/internal/pkg/smtp_server"
)

type authRepo interface {
	Create(ctx context.Context, u model.User) error
}

type Service struct {
	cfg         config.ServerConfig
	smtp_server *smtp_server.SMTPServer
	logger      logger.Logger
	authRepo    authRepo
}

func New(
	cfg config.ServerConfig,
	smtp_server *smtp_server.SMTPServer,
	logger logger.Logger,
	authRepo authRepo,
) *Service {
	return &Service{
		cfg:         cfg,
		smtp_server: smtp_server,
		logger:      logger,
		authRepo:    authRepo,
	}
}

func (s *Service) Register(ctx context.Context, credentials model.RegistrationCredentials) error {
	err := s.authRepo.Create(ctx, model.User{
		ID:       uuid.NewString(),
		Nickname: &credentials.Nickname,
		Email:    &credentials.Email,
		Password: lo.ToPtr(s.hash(credentials.Password)),
		Role:     lo.ToPtr("user"),
	})
	if err != nil {
		return err
	}

	go func() {
		subject := "Вы зарегистрированы в Titanic!"
		body := fmt.Sprintf(
			"Ваш логин: %s\nВаш пароль: %s",
			credentials.Email,
			credentials.Password,
		)
		err := s.smtp_server.SendEmail(
			credentials.Email,
			s.smtp_server.GetEmail(),
			subject,
			body,
		)
		if err != nil {
			s.logger.Err(
				fmt.Printf(
					"Failed to send welcome email to %s: %v",
					credentials.Nickname,
					err,
				),
			)
		}
	}()

	return nil
}

func (s *Service) hash(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(s.cfg.SecretKey))

	return fmt.Sprintf("%x", pwd.Sum(nil))
}
