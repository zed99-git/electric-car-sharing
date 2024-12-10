package vehicleService

import (
	"database/sql"
	"fmt"
	"time"
)

// VehicleCleanliness enum
type VehicleCleanliness string

const (
	Clean   VehicleCleanliness = "Clean"
	Unclean VehicleCleanliness = "Needs Cleaning"
)

type VehicleStatus string

const (
	Available   VehicleStatus = "Available"
	Reserved    VehicleStatus = "Reserved"
	Maintenance VehicleStatus = "Maintenance"
)

type Vehicle struct {
	VehicleID    int                `json:"vehicle_id"`
	Make         string             `json:"make"`
	Model        string             `json:"model"`
	Year         int                `json:"year"`
	LicensePlate string             `json:"license_plate"`
	Status       VehicleStatus      `json:"status"`
	Location     string             `json:"location"`
	ChargeLevel  float64            `json:"charge_level"`
	Cleanliness  VehicleCleanliness `json:"cleanliness"`
}

// GetAllVehicles fetches all vehicle entries from the database
func GetAllVehicles(db *sql.DB) ([]Vehicle, error) {
	rows, err := db.Query("SELECT VehicleID, Make, Model, Year, LicensePlate, Status, Location, ChargeLevel, Cleanliness FROM vehicles WHERE Status = 'Available'")
	if err != nil {
		return nil, fmt.Errorf("error fetching vehicles: %v", err)
	}
	defer rows.Close()

	// Parse rows into a slice
	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		err := rows.Scan(
			&v.VehicleID,
			&v.Make,
			&v.Model,
			&v.Year,
			&v.LicensePlate,
			&v.Status,
			&v.Location,
			&v.ChargeLevel,
			&v.Cleanliness,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

// BookVehicle attempts to book a vehicle for the provided time range
func BookVehicle(db *sql.DB, userID, vehicleID int, startTime, endTime time.Time) error {
	// Check if the vehicle is available in the desired time range
	var count int
	query := `
		SELECT COUNT(*) 
		FROM Reservations 
		WHERE VehicleID = ? AND Status = 'Active' 
		AND ((StartTime BETWEEN ? AND ?) OR (EndTime BETWEEN ? AND ?))
	`
	err := db.QueryRow(query, vehicleID, startTime, endTime, startTime, endTime).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking availability: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("vehicle is already reserved for the selected time range")
	}

	// Insert reservation
	insertQuery := `
		INSERT INTO Reservations (UserID, VehicleID, StartTime, EndTime, Status)
		VALUES (?, ?, ?, ?, 'Active')
	`
	_, err = db.Exec(insertQuery, userID, vehicleID, startTime, endTime)
	if err != nil {
		return fmt.Errorf("error booking the vehicle: %v", err)
	}

	return nil
}

// CancelReservation cancels an active reservation
func CancelReservation(db *sql.DB, reservationID int) error {
	// Ensure the reservation exists
	var status string
	query := `SELECT Status FROM Reservations WHERE ReservationID = ?`
	err := db.QueryRow(query, reservationID).Scan(&status)
	if err != nil {
		return fmt.Errorf("error finding reservation: %v", err)
	}

	if status != "Active" {
		return fmt.Errorf("reservation cannot be cancelled as it is not active")
	}

	// Cancel the reservation
	updateQuery := `UPDATE Reservations SET Status = 'Cancelled' WHERE ReservationID = ?`
	_, err = db.Exec(updateQuery, reservationID)
	if err != nil {
		return fmt.Errorf("error cancelling the reservation: %v", err)
	}

	return nil
}
