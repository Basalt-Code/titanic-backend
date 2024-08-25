package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"cmd/app/main.go/internal/api/models"
	"cmd/app/main.go/internal/model"
	"cmd/app/main.go/internal/pkg/http_resp"
	repo "cmd/app/main.go/internal/repository"
)

func (api *API) Register(c *gin.Context) {
	var req apimodels.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WithJSONError(c, err)
		return
	}

	if err := api.service.Register(c, model.RegistrationCredentials{
		Nickname: strings.ToLower(req.Nickname),
		Email:    strings.ToLower(req.Email),
		Password: req.Password,
	}); err != nil {
		var errDupl repo.ErrDuplicateField
		if errors.As(err, &errDupl) {
			response.WithBadRequestError(c, err.Error())
			return
		}

		response.WithInternalServerError(c)
		api.logger.Err(err)

		return
	}

	c.Status(http.StatusOK)
}
