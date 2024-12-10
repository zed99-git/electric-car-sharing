package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zed99-git/electric-car-sharing/middleware"
	"github.com/zed99-git/electric-car-sharing/routes"
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
	router.HandleFunc("/api/v1/protected", middleware.AuthMiddleware(routes.ProtectedRouteHandler)).Methods("GET") //testing middleware token
	router.HandleFunc("/api/v1/users/register", routes.RegisterUserHandler).Methods("POST")                        //register
	router.HandleFunc("/api/v1/users/login", routes.LoginUserHandler).Methods("POST")                              //login

	router.HandleFunc("/api/v1/vehicles", routes.GetAllVehiclesHandler).Methods("GET")
	router.HandleFunc("/api/v1/vehicles/book", routes.BookVehicleHandler).Methods("POST")
	router.HandleFunc("/api/v1/vehicles/cancel", routes.CancelReservationHandler).Methods("POST")

	router.HandleFunc("/api/v1/billing/calculate", routes.CalculateBillingHandler).Methods("POST")
	router.HandleFunc("/api/v1/billing/realtime", routes.RealTimeBillingHandler).Methods("POST")
	router.HandleFunc("/api/v1/billing/generate-invoice", routes.GenerateInvoiceHandler).Methods("POST")

	http.Handle("/", enableCORS(router))

	fmt.Println("Listening on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", router))
}
