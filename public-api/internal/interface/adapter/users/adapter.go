package iface

import "public-api/internal/domain/dto"

type UserServiceClient interface {
	CreateUser(name string) (dto.UserResult, error)
	GetUser(userID int) (dto.UserResult, error)
}
