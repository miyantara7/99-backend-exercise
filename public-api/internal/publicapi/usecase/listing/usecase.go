package usecase

import (
	"net/http"

	"public-api/internal/domain/dto"
	ifaceListing "public-api/internal/interface/adapter/listing"
	ifaceUsecase "public-api/internal/interface/usecase/listing"
	ifaceUsersUsecase "public-api/internal/interface/usecase/users"
)

type Usecase struct {
	usersUsecase  ifaceUsersUsecase.PublicAPIUsecase
	listingClient ifaceListing.ListingServiceClient
}

func New(
	usersUsecase ifaceUsersUsecase.PublicAPIUsecase,
	listingClient ifaceListing.ListingServiceClient,
) ifaceUsecase.PublicAPIUsecase {
	return &Usecase{
		usersUsecase:  usersUsecase,
		listingClient: listingClient,
	}
}

func (u *Usecase) ListListings(rawQuery string) (dto.ListingViewsResult, error) {
	listingResult, err := u.listingClient.ListListings(rawQuery)
	if err != nil {
		return dto.ListingViewsResult{
			Result: dto.Result{Status: http.StatusBadGateway, Errors: []string{"failed to reach listing service"}},
		}, err
	}
	if listingResult.Status != http.StatusOK {
		return dto.ListingViewsResult{
			Result: dto.Result{Status: listingResult.Status, Errors: []string{"listing service returned an error"}},
		}, nil
	}

	usersByID := make(map[int]dto.UserResult, len(listingResult.Listings))
	result := make([]dto.ListingView, 0, len(listingResult.Listings))
	for _, listing := range listingResult.Listings {
		user, ok := usersByID[listing.UserID]
		if !ok {
			user, err = u.usersUsecase.GetUser(listing.UserID)
			if err != nil {
				return dto.ListingViewsResult{
					Result: dto.Result{Status: http.StatusBadGateway, Errors: []string{"failed to reach user service"}},
				}, err
			}
			if user.Status != http.StatusOK {
				if len(user.Errors) == 0 {
					user.Errors = []string{"failed to fetch user from user service"}
				}
				return dto.ListingViewsResult{
					Result: dto.Result{Status: user.Status, Errors: user.Errors},
				}, nil
			}
			usersByID[listing.UserID] = user
		}

		result = append(result, dto.ListingView{
			ID:          listing.ID,
			ListingType: listing.ListingType,
			Price:       listing.Price,
			CreatedAt:   listing.CreatedAt,
			UpdatedAt:   listing.UpdatedAt,
			User:        user.User,
		})
	}

	return dto.ListingViewsResult{
		Result:   dto.Result{Status: http.StatusOK},
		Listings: result,
	}, nil
}

func (u *Usecase) CreateListing(userID int, listingType string, price int) (dto.ListingResult, error) {
	userResult, err := u.usersUsecase.GetUser(userID)
	if err != nil {
		return dto.ListingResult{
			Result: dto.Result{Status: http.StatusBadGateway, Errors: []string{"failed to reach user service"}},
		}, err
	}
	if userResult.Status != http.StatusOK {
		if userResult.Status == http.StatusNotFound {
			return dto.ListingResult{
				Result: dto.Result{Status: http.StatusBadRequest, Errors: []string{"invalid user_id"}},
			}, nil
		}
		if len(userResult.Errors) == 0 {
			userResult.Errors = []string{"user service returned an error"}
		}
		return dto.ListingResult{
			Result: dto.Result{Status: userResult.Status, Errors: userResult.Errors},
		}, nil
	}

	listingResult, err := u.listingClient.CreateListing(userID, listingType, price)
	if err != nil {
		return dto.ListingResult{
			Result: dto.Result{Status: http.StatusBadGateway, Errors: []string{"failed to reach listing service"}},
		}, err
	}
	if listingResult.Status != http.StatusOK {
		return dto.ListingResult{
			Result: dto.Result{Status: listingResult.Status, Errors: []string{"listing service returned an error"}},
		}, nil
	}

	return dto.ListingResult{
		Result:  dto.Result{Status: http.StatusOK},
		Listing: listingResult.Listing,
	}, nil
}
