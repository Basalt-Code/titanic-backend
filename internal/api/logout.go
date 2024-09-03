package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	response "cmd/app/main.go/internal/pkg/http_resp"
)

func (api *API) Logout(c *gin.Context) {
	if err := api.authService.Logout(c, c.GetString(userIDParam), c.GetString(sessionIDParam)); err != nil {
		response.WithInternalServerError(c)
		api.logger.Err(err)

		return
	}

	c.Status(http.StatusOK)
}
