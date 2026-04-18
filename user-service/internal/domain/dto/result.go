package dto

import "user-service/internal/domain/model"

type Result struct {
	Status int
	Errors []string
}

func (r Result) HasErrors() bool {
	return len(r.Errors) > 0
}

type UserResult struct {
	Result
	User model.User
}

type UsersResult struct {
	Result
	Users []model.User
}
