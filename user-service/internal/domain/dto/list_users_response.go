package dto

import "user-service/internal/domain/model"

type ListUsersResponse struct {
	Result bool         `json:"result"`
	Users  []model.User `json:"users"`
}
