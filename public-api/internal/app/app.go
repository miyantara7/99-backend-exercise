package app

import (
	"flag"
	"fmt"
	"log"
	"time"

	"public-api/internal/adapter"
	"public-api/internal/bootstrap"
)

func Run() {
	port := flag.String("port", "7300", "port to listen on")
	listingServiceURL := flag.String("listing-service-url", "http://127.0.0.1:6300", "listing service base URL")
	userServiceURL := flag.String("user-service-url", "http://127.0.0.1:6301", "user service base URL")
	flag.Parse()

	router := bootstrap.NewRouter(adapter.Config{
		ListingServiceURL: *listingServiceURL,
		UserServiceURL:    *userServiceURL,
		UserTimeout:       30 * time.Second,
		ListingTimeout:    10 * time.Second,
	})

	log.Printf("public api listening on :%s", *port)
	log.Fatal(router.Run(fmt.Sprintf(":%s", *port)))
}
