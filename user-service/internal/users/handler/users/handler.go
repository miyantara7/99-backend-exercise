package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"user-service/internal/domain/dto"
	iface "user-service/internal/interface/usecase/users"
	"user-service/internal/platform"
)

type Handler struct {
	usecase iface.Usecase
}

func NewUsers(uc iface.Usecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		platform.WriteJSON(c, 400, dto.ErrorResponse{Result: false, Errors: []string{"invalid user_id"}})
		return
	}

	result := h.usecase.GetUser(id)
	if platform.WriteResultError(c, result.Result, nil) {
		return
	}

	platform.WriteJSON(c, 200, dto.GetUserResponse{Result: true, User: result.User})
}

func (h *Handler) ListUsers(c *gin.Context) {
	pageNum, pageSize, errors := platform.ParsePageParams(c)
	if len(errors) > 0 {
		platform.WriteJSON(c, 400, dto.ErrorResponse{Result: false, Errors: errors})
		return
	}

	result := h.usecase.ListUsers(pageNum, pageSize)
	platform.WriteJSON(c, 200, dto.ListUsersResponse{Result: true, Users: result.Users})
}

func (h *Handler) CreateUser(c *gin.Context) {
	result, err := h.usecase.CreateUser(c.PostForm("name"))
	if platform.WriteResultError(c, result.Result, err) {
		return
	}

	platform.WriteJSON(c, 200, dto.CreateUserResponse{Result: true, User: result.User})
}
