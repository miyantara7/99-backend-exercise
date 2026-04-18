package listingservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"public-api/internal/domain/dto"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *Client) ListListings(rawQuery string) (dto.ListingsResult, error) {
	endpoint := c.baseURL + "/listings"
	if rawQuery != "" {
		endpoint += "?" + rawQuery
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return dto.ListingsResult{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return dto.ListingsResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.ListingsResult{}, err
	}

	var payload dto.ClientResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return dto.ListingsResult{}, fmt.Errorf("invalid listing service response: %w", err)
	}

	return dto.ListingsResult{
		Result:   dto.Result{Status: resp.StatusCode, Errors: payload.Errors},
		Listings: payload.Listings,
	}, nil
}

func (c *Client) CreateListing(userID int, listingType string, price int) (dto.ListingResult, error) {
	form := url.Values{}
	form.Set("user_id", strconv.Itoa(userID))
	form.Set("listing_type", listingType)
	form.Set("price", strconv.Itoa(price))

	resp, err := c.httpClient.Post(
		c.baseURL+"/listings",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return dto.ListingResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.ListingResult{}, err
	}

	var payload dto.ClientResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return dto.ListingResult{}, fmt.Errorf("invalid listing service response: %w", err)
	}

	return dto.ListingResult{
		Result:  dto.Result{Status: resp.StatusCode, Errors: payload.Errors},
		Listing: payload.Listing,
	}, nil
}
