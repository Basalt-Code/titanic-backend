package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	response "cmd/app/main.go/internal/pkg/http_resp"
)

const (
	userIDParam    = "userID"
	sessionIDParam = "sessionID"

	headerAuthFormatErrMsg = "header: \"Authorization\": must be in the format: Bearer <token>"
)

func (api *API) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.Split(c.GetHeader("Authorization"), " ")
		if len(authHeader) != 2 {
			response.WithUnauthorizedError(c, headerAuthFormatErrMsg)
			c.Abort()

			return
		}

		token, claims, err := api.authService.ParseTokenWitClaims(authHeader[1])
		if err != nil {
			response.WithInternalServerError(c)
			api.logger.Err(fmt.Errorf("token parsing error: %w", err))
			c.Abort()

			return
		}

		if !token.Valid {
			response.WithUnauthorizedError(c, "token is invalid")
			c.Abort()

			return
		}

		c.Set(userIDParam, claims.UserID)
		c.Set(sessionIDParam, claims.SessionID)

		c.Next()
	}
}
