package route

import (
	"github.com/gin-gonic/gin"
)

func Register(router gin.IRouter) {
	handlers := BuildHandlers()

	router.GET("/users", handlers.Users.ListUsers)
	router.POST("/users", handlers.Users.CreateUser)
	router.GET("/users/:id", handlers.Users.GetUserByID)
}
