package adapter

import (
	"net/http"
	"time"

	"public-api/internal/adapter/listingservice"
	"public-api/internal/adapter/userservice"
	ifaceListing "public-api/internal/interface/adapter/listing"
	ifaceUsers "public-api/internal/interface/adapter/users"
)

var (
	Connection    *http.Client
	UserClient    ifaceUsers.UserServiceClient
	ListingClient ifaceListing.ListingServiceClient
)

type Config struct {
	ListingServiceURL string
	UserServiceURL    string
	UserTimeout       time.Duration
	ListingTimeout    time.Duration
}

func InitConnections(cfg Config) {
	UserClient = userservice.NewClient(cfg.UserServiceURL, &http.Client{
		Timeout: cfg.UserTimeout,
	})

	ListingClient = listingservice.NewClient(cfg.ListingServiceURL, &http.Client{
		Timeout: cfg.ListingTimeout,
	})
}
