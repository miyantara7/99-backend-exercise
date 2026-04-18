package dto

import "public-api/internal/domain/model"

type Result struct {
	Status int
	Errors []string
}

func (r Result) HasErrors() bool {
	return len(r.Errors) > 0
}

type UserResult struct {
	Result
	User model.User
}

type ListingResult struct {
	Result
	Listing model.Listing
}

type ListingsResult struct {
	Result
	Listings []model.Listing
}

type ListingViewsResult struct {
	Result
	Listings []ListingView
}
