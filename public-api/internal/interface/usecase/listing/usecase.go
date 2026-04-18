package iface

import (
	"public-api/internal/domain/dto"
)

type PublicAPIUsecase interface {
	ListListings(rawQuery string) (dto.ListingViewsResult, error)
	CreateListing(userID int, listingType string, price int) (dto.ListingResult, error)
}
