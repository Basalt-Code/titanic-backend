package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apimodels "cmd/app/main.go/internal/api/models"
	"cmd/app/main.go/internal/model/domain"
	response "cmd/app/main.go/internal/pkg/http_resp"
	"cmd/app/main.go/internal/services/auth"
)

func (api *API) RefreshToken(c *gin.Context) {
	var req apimodels.RefreshTokenReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.WithJSONError(c, err)
		return
	}

	authHeader := strings.Split(c.GetHeader("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		response.WithUnauthorizedError(c, headerAuthFormatErrMsg)
		return
	}

	tokens, err := api.authService.RefreshToken(c, domain.Tokens{
		Access:  authHeader[1],
		Refresh: req.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidTokenPair) {
			response.WithUnauthorizedError(c, "invalid token pair")
			return
		}

		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			response.WithUnauthorizedError(c, "invalid refresh token")
			return
		}

		response.WithInternalServerError(c)
		api.logger.Err(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": apimodels.TokensResp(tokens),
	})
}
