package iface

import "user-service/internal/domain/model"

type Repository interface {
	List(pageNum int, pageSize int) []model.User
	GetByID(id int) (model.User, bool)
	Create(name string) (model.User, error)
}
