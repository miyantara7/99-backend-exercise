package usecase

import (
	"net/http"

	"public-api/internal/domain/dto"
	ifaceAdapter "public-api/internal/interface/adapter/users"
	ifaceUsecase "public-api/internal/interface/usecase/users"
)

type Usecase struct {
	userClient ifaceAdapter.UserServiceClient
}

func NewUsers(userClient ifaceAdapter.UserServiceClient) ifaceUsecase.PublicAPIUsecase {
	return &Usecase{
		userClient: userClient,
	}
}

func (u *Usecase) CreateUser(name string) (dto.UserResult, error) {
	result, err := u.userClient.CreateUser(name)
	if err != nil {
		return dto.UserResult{
			Result: dto.Result{Status: http.StatusBadGateway, Errors: []string{"failed to reach user service"}},
		}, err
	}
	return result, nil
}

func (u *Usecase) GetUser(id int) (dto.UserResult, error) {
	result, err := u.userClient.GetUser(id)
	if err != nil {
		return dto.UserResult{
			Result: dto.Result{Status: http.StatusBadGateway, Errors: []string{"failed to reach user service"}},
		}, err
	}
	return result, nil
}
