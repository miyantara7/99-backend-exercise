package bootstrap

import (
	"github.com/gin-gonic/gin"

	"public-api/internal/adapter"
	"public-api/internal/publicapi/route"
)

func NewRouter(cfg adapter.Config) *gin.Engine {
	adapter.InitConnections(cfg)
	router := gin.Default()
	route.Register(router)
	return router
}
