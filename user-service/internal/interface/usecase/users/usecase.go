package iface

import "user-service/internal/domain/dto"

type Usecase interface {
	ListUsers(pageNum int, pageSize int) dto.UsersResult
	GetUser(id int) dto.UserResult
	CreateUser(name string) (dto.UserResult, error)
}
