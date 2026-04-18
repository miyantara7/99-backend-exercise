package route

import (
	"public-api/internal/adapter"
	handlerListing "public-api/internal/publicapi/handler/listing"
	handlerUsers "public-api/internal/publicapi/handler/users"
	listingUc "public-api/internal/publicapi/usecase/listing"
	usersUc "public-api/internal/publicapi/usecase/users"
)

type Handlers struct {
	Users   *handlerUsers.Handler
	Listing *handlerListing.Handler
}

func BuildHandlers() Handlers {
	usersUsecase := usersUc.NewUsers(adapter.UserClient)
	listingUsecase := listingUc.New(usersUsecase, adapter.ListingClient)

	return Handlers{
		Users:   handlerUsers.NewUsers(usersUsecase),
		Listing: handlerListing.NewListing(listingUsecase),
	}
}
