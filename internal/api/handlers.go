package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"cmd/app/main.go/internal/model"
)

type service interface {
	Register(ctx context.Context, credentials model.RegistrationCredentials) error
}

type logger interface {
	Err(text ...any)
	Warn(text ...any)
	Info(text ...any)
}

type API struct {
	logger  logger
	service service
}

func New(logger logger, service service) *gin.Engine {
	h := API{logger, service}

	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api")
	api.POST("/auth/register", h.Register)

	return r
}
