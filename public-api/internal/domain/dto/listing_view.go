package dto

import "public-api/internal/domain/model"

type ListingView struct {
	ID          int        `json:"id"`
	ListingType string     `json:"listing_type"`
	Price       int        `json:"price"`
	CreatedAt   int64      `json:"created_at"`
	UpdatedAt   int64      `json:"updated_at"`
	User        model.User `json:"user"`
}
