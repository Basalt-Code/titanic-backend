package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apimodels "cmd/app/main.go/internal/api/models"
	"cmd/app/main.go/internal/model/dto"
	response "cmd/app/main.go/internal/pkg/http_resp"
	"cmd/app/main.go/internal/services/auth"
)

const invalidCredentials = "Неправильный логин или пароль"

func (api *API) Login(c *gin.Context) {
	var req apimodels.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WithJSONError(c, err)
		return
	}

	tokens, err := api.authService.Login(c, dto.Credentials{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			response.WithUnauthorizedError(c, invalidCredentials)
			return
		}

		response.WithInternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
