package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zed99-git/electric-car-sharing/database"
	billingService "github.com/zed99-git/electric-car-sharing/services/billing-service"
)

// Handler for calculating billing
func CalculateBillingHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Decode request body
	var req struct {
		UserID    int    `json:"user_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse times
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

	// Call billing service to compute cost
	totalCost, err := billingService.CalculateBilling(db, req.UserID, startTime, endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with calculated billing
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"total_amount": totalCost,
	})
}

// Handler for real-time billing
func RealTimeBillingHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req struct {
		UserID           int    `json:"user_id"`
		ReservationStart string `json:"reservation_start"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse the reservation start time
	startTime, err := time.Parse(time.RFC3339, req.ReservationStart)
	if err != nil {
		http.Error(w, "Invalid time format", http.StatusBadRequest)
		return
	}

	// Call the real-time service to calculate the estimated cost
	estimatedCost, err := billingService.RealTimeBilling(db, req.UserID, startTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"estimated_total": estimatedCost,
	})
}

// Handler for generating invoice
func GenerateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req struct {
		ReservationID int     `json:"reservation_id"`
		UserID        int     `json:"user_id"`
		Amount        float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service to generate invoice
	err = billingService.GenerateInvoice(db, req.ReservationID, req.UserID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Invoice generated successfully"}`))
}
