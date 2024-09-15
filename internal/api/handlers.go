package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"cmd/app/main.go/internal/model/domain"
	"cmd/app/main.go/internal/model/dto"
)

type authService interface {
	Register(ctx context.Context, credentials dto.RegistrationCredentials) error
	Login(ctx context.Context, credentials dto.Credentials) (domain.Tokens, error)
	Logout(ctx context.Context, userID, sessionID string) error
	ParseTokenWitClaims(access string) (*jwt.Token, *domain.Claims, error)
	RefreshToken(ctx context.Context, tokens domain.Tokens) (domain.Tokens, error)
}

type logger interface {
	Err(text ...any)
	Warn(text ...any)
	Info(text ...any)
}

type API struct {
	logger      logger
	authService authService
}

func New(logger logger, authService authService) *gin.Engine {
	h := API{logger, authService}

	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api")
	api.POST("/auth/register", h.Register)
	api.POST("/auth/login", h.Login)
	api.POST("/auth/refresh", h.RefreshToken)

	api.Use(h.Auth())
	api.POST("/auth/logout", h.Logout)

	return r
}
