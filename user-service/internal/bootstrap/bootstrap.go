package bootstrap

import (
	"github.com/gin-gonic/gin"

	"user-service/internal/users/route"
)

func NewRouter() (*gin.Engine, error) {
	router := gin.Default()
	route.Register(router)
	return router, nil
}
