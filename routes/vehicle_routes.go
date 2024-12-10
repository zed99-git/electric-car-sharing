package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zed99-git/electric-car-sharing/database"
	vehicleService "github.com/zed99-git/electric-car-sharing/services/vehicle-service"
)

// Handler for getting all vehicles
func GetAllVehiclesHandler(w http.ResponseWriter, r *http.Request) {
	// Establish database connection
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Call the service layer function to fetch vehicles
	vehicles, err := vehicleService.GetAllVehicles(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Send the response
	if err := json.NewEncoder(w).Encode(vehicles); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Handler for booking a vehicle
func BookVehicleHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req struct {
		UserID    int    `json:"user_id"`
		VehicleID int    `json:"vehicle_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	// Call service to handle booking
	err = vehicleService.BookVehicle(db, req.UserID, req.VehicleID, startTime, endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Vehicle booked successfully"))
}

// Handler for cancelling a reservation
func CancelReservationHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req struct {
		ReservationID int `json:"reservation_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = vehicleService.CancelReservation(db, req.ReservationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reservation cancelled successfully"))
}
