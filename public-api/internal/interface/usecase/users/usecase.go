package iface

import "public-api/internal/domain/dto"

type PublicAPIUsecase interface {
	CreateUser(name string) (dto.UserResult, error)
	GetUser(id int) (dto.UserResult, error)
}
