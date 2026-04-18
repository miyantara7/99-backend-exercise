package dto

import "user-service/internal/domain/model"

type GetUserResponse struct {
	Result bool       `json:"result"`
	User   model.User `json:"user"`
}
