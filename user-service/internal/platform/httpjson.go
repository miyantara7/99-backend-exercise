package platform

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"user-service/internal/domain/dto"
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

func WriteResultError(c *gin.Context, result dto.Result, err error) bool {
	if err != nil || result.HasErrors() {
		WriteJSON(c, result.Status, dto.ErrorResponse{Result: false, Errors: result.Errors})
		return true
	}

	return false
}

func ParsePositiveInt(raw string, field string, allowZero bool) (int, string) {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, "invalid " + field
	}
	if (!allowZero && value < 1) || (allowZero && value < 0) {
		return 0, "invalid " + field
	}
	return value, ""
}

func ParsePageParams(c *gin.Context) (pageNum int, pageSize int, errors []string) {
	pageNum = 1
	pageSize = 10

	if raw := c.Query("page_num"); raw != "" {
		value, msg := ParsePositiveInt(raw, "page_num", false)
		if msg != "" {
			errors = append(errors, msg)
		} else {
			pageNum = value
		}
	}

	if raw := c.Query("page_size"); raw != "" {
		value, msg := ParsePositiveInt(raw, "page_size", false)
		if msg != "" {
			errors = append(errors, msg)
		} else {
			pageSize = value
		}
	}

	return pageNum, pageSize, errors
}
