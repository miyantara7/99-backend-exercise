package app

import (
	"flag"
	"fmt"
	"log"

	"user-service/internal/bootstrap"
)

func Run() {
	port := flag.String("port", "6301", "port to listen on")
	flag.Parse()

	router, err := bootstrap.NewRouter()
	if err != nil {
		log.Fatalf("failed to bootstrap user service: %v", err)
	}

	log.Printf("user service listening on :%s", *port)
	log.Fatal(router.Run(fmt.Sprintf(":%s", *port)))
}
