package dto

import "public-api/internal/domain/model"

type CreateListingResponse struct {
	Listing model.Listing `json:"listing"`
}
