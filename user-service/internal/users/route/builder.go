package route

import (
	"flag"
	"path/filepath"
	dbUser "user-service/internal/users/db"
	handlerUsers "user-service/internal/users/handler/users"
	"user-service/internal/users/usecase"
)

type Handlers struct {
	Users *handlerUsers.Handler
}

func BuildHandlers() Handlers {
	dataFile := flag.String("data-file", filepath.Join("data", "users.json"), "path to users data file")
	flag.Parse()

	repo, err := dbUser.NewFileRepository(*dataFile)
	if err != nil {
		panic(err)
	}

	userUsecase := usecase.New(repo)

	return Handlers{
		Users: handlerUsers.NewUsers(userUsecase),
	}
}
