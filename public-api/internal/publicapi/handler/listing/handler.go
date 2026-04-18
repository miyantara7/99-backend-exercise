package handler

import (
	"public-api/internal/domain/dto"
	iface "public-api/internal/interface/usecase/listing"
	"public-api/internal/platform"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase iface.PublicAPIUsecase
}

func NewListing(uc iface.PublicAPIUsecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) ListListings(c *gin.Context) {
	result, err := h.usecase.ListListings(c.Request.URL.RawQuery)
	if platform.WriteResultError(c, result.Result, err) {
		return
	}

	platform.WriteJSON(c, 200, dto.ListingsResponse{Result: true, Listings: result.Listings})
}

func (h *Handler) CreateListing(c *gin.Context) {
	var payload *dto.CreateListingRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		platform.WriteBindError(c)
		return
	}

	result, err := h.usecase.CreateListing(payload.UserID, payload.ListingType, payload.Price)
	if platform.WriteResultError(c, result.Result, err) {
		return
	}

	platform.WriteJSON(c, 200, dto.CreateListingResponse{Listing: result.Listing})
}
