package main

import (
	"fmt"
	"log"

	"github.com/zed99-git/electric-car-sharing/middleware"
)

func main() {
	// Generate token for testing purposes
	token, err := middleware.GenerateToken("user@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generated Token for Testing:", token)
}
