package main

import (
	"fmt"
	"log"
	"net/http"

	"electric-car-sharing/database/routes"

	"github.com/gorilla/mux"
)

// Enable CORS middleware
func enableCORS(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	}
}

func main() {
	// Initialize router
	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/api/v1/users/register", routes.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/api/v1/users/login", routes.LoginUserHandler).Methods("POST")

	http.Handle("/", enableCORS(router))

	fmt.Println("Listening on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", router))
}
