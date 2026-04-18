package dto

type ListingsResponse struct {
	Result   bool          `json:"result"`
	Listings []ListingView `json:"listings"`
}
