package usecase

import (
	"net/http"
	"strings"

	"user-service/internal/domain/dto"
	iface "user-service/internal/interface/db/users"
	ifaceUsecase "user-service/internal/interface/usecase/users"
)

type Usecase struct {
	repo iface.Repository
}

func New(repo iface.Repository) ifaceUsecase.Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) ListUsers(pageNum int, pageSize int) dto.UsersResult {
	return dto.UsersResult{
		Result: dto.Result{Status: http.StatusOK},
		Users:  u.repo.List(pageNum, pageSize),
	}
}

func (u *Usecase) GetUser(id int) dto.UserResult {
	user, found := u.repo.GetByID(id)
	if !found {
		return dto.UserResult{
			Result: dto.Result{Status: http.StatusNotFound, Errors: []string{"user not found"}},
		}
	}

	return dto.UserResult{
		Result: dto.Result{Status: http.StatusOK},
		User:   user,
	}
}

func (u *Usecase) CreateUser(name string) (dto.UserResult, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return dto.UserResult{
			Result: dto.Result{Status: http.StatusBadRequest, Errors: []string{"invalid name"}},
		}, nil
	}

	user, err := u.repo.Create(name)
	if err != nil {
		return dto.UserResult{
			Result: dto.Result{Status: http.StatusInternalServerError, Errors: []string{"failed to create user"}},
		}, err
	}

	return dto.UserResult{
		Result: dto.Result{Status: http.StatusOK},
		User:   user,
	}, nil
}
