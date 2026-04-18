package platform

import (
	"public-api/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

func WriteJSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

func WriteError(c *gin.Context, status int, errors ...string) {
	WriteJSON(c, status, map[string]any{
		"result": false,
		"errors": errors,
	})
}

func WriteBindError(c *gin.Context) {
	WriteJSON(c, 400, dto.ErrorResponse{Result: false, Errors: []string{"invalid json body"}})
}

func WriteResultError(c *gin.Context, result dto.Result, err error) bool {
	if err != nil || result.HasErrors() {
		WriteJSON(c, result.Status, dto.ErrorResponse{Result: false, Errors: result.Errors})
		return true
	}

	return false
}
