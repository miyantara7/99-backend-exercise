package dto

import "user-service/internal/domain/model"

type CreateUserResponse struct {
	Result bool       `json:"result"`
	User   model.User `json:"user"`
}
