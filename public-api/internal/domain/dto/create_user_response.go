package dto

import "public-api/internal/domain/model"

type CreateUserResponse struct {
	User model.User `json:"user"`
}
