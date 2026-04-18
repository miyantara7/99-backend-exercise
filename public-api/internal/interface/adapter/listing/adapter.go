package iface

import "public-api/internal/domain/dto"

type ListingServiceClient interface {
	ListListings(rawQuery string) (dto.ListingsResult, error)
	CreateListing(userID int, listingType string, price int) (dto.ListingResult, error)
}
