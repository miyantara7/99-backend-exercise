package dto

import "public-api/internal/domain/model"

type ClientResponse struct {
	Result   bool            `json:"result"`
	User     model.User      `json:"user,omitempty"`
	Listing  model.Listing   `json:"listing,omitempty"`
	Listings []model.Listing `json:"listings,omitempty"`
	Errors   []string        `json:"errors,omitempty"`
}
