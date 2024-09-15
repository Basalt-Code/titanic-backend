package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error any `json:"error"`
}

type Error struct {
	Message string `json:"message"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func WithJSONError(c *gin.Context, err error) {
	var jSyntaxErr *json.SyntaxError
	if errors.As(err, &jSyntaxErr) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: Error{
				Message: jSyntaxErr.Error(),
			},
		})

		return
	}

	var jErr *json.UnmarshalTypeError
	if errors.As(err, &jErr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": FieldError{
				Field: jErr.Field,
				Message: fmt.Sprintf(
					"invalid data type for field '%s': expected %s",
					jErr.Field,
					jErr.Type.String(),
				),
			},
		})

		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error{
			Message: err.Error(),
		},
	})

	return
}

func WithUnauthorizedError(c *gin.Context, errMsg string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: Error{
			Message: errMsg,
		},
	})
}

func WithInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: Error{
			Message: "что-то пошло не так, попробуйте позже.",
		},
	})
}

func WithBadRequestError(c *gin.Context, errMsg string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: Error{
			Message: errMsg,
		},
	})
}
