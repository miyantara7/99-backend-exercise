package route

import (
	"github.com/gin-gonic/gin"
)

func Register(router gin.IRouter) {
	handlers := BuildHandlers()

	router.POST("/public-api/users", handlers.Users.CreateUser)
	router.GET("/public-api/listings", handlers.Listing.ListListings)
	router.POST("/public-api/listings", handlers.Listing.CreateListing)
}
