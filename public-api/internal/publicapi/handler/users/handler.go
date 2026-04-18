package handler

import (
	"public-api/internal/domain/dto"
	iface "public-api/internal/interface/usecase/users"
	"public-api/internal/platform"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase iface.PublicAPIUsecase
}

func NewUsers(uc iface.PublicAPIUsecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var payload *dto.CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		platform.WriteBindError(c)
		return
	}

	result, err := h.usecase.CreateUser(payload.Name)
	if platform.WriteResultError(c, result.Result, err) {
		return
	}

	platform.WriteJSON(c, 200, dto.CreateUserResponse{User: result.User})
}
