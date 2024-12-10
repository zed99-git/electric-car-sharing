package billingService

import (
	"database/sql"
	"fmt"
	"time"
)

// Membership tiers with base rates
const (
	BasicRate   = 10.0 // Base hourly rate for Basic members
	PremiumRate = 7.0  // Base hourly rate for Premium members
	VIPRate     = 5.0  // Base hourly rate for VIP members
)

// Calculate billing based on user membership and time
func CalculateBilling(db *sql.DB, userID int, reservationStart, reservationEnd time.Time) (float64, error) {
	// Get the user's membership tier from the database
	var membershipTier string
	err := db.QueryRow("SELECT MembershipTier FROM users WHERE UserID = ?", userID).Scan(&membershipTier)
	if err != nil {
		return 0, fmt.Errorf("error fetching user membership tier: %v", err)
	}

	// Calculate the duration of reservation
	duration := reservationEnd.Sub(reservationStart).Hours()
	if duration <= 0 {
		return 0, fmt.Errorf("invalid reservation duration")
	}

	// Select appropriate rate based on membership tier
	var rate float64
	switch membershipTier {
	case "Basic":
		rate = BasicRate
	case "Premium":
		rate = PremiumRate
	case "VIP":
		rate = VIPRate
	default:
		return 0, fmt.Errorf("unknown membership tier")
	}

	// Calculate total cost
	totalCost := duration * rate

	// Apply a discount for certain membership tiers if applicable
	if membershipTier == "VIP" {
		// VIP members get a 20% discount
		totalCost *= 0.8
	}

	return totalCost, nil
}

// RealTimeBilling computes the current estimated cost for a user in real-time
func RealTimeBilling(db *sql.DB, userID int, reservationStart time.Time) (float64, error) {
	// Fetch user's membership tier from the database
	var membershipTier string
	err := db.QueryRow("SELECT MembershipTier FROM users WHERE UserID = ?", userID).Scan(&membershipTier)
	if err != nil {
		return 0, fmt.Errorf("error fetching user membership tier: %v", err)
	}

	// Calculate elapsed time
	currentTime := time.Now()
	duration := currentTime.Sub(reservationStart).Hours()
	if duration < 0 {
		return 0, fmt.Errorf("invalid start time")
	}

	// Assign rate based on membership tier
	var rate float64
	switch membershipTier {
	case "Basic":
		rate = BasicRate
	case "Premium":
		rate = PremiumRate
	case "VIP":
		rate = VIPRate
	default:
		return 0, fmt.Errorf("unknown membership tier")
	}

	// Calculate total cost dynamically based on duration
	estimatedCost := duration * rate

	// Apply discount if VIP
	if membershipTier == "VIP" {
		estimatedCost *= 0.8
	}

	return estimatedCost, nil
}

// GenerateInvoice generates an invoice based on a reservation and saves it to the database
func GenerateInvoice(db *sql.DB, reservationID int, userID int, amount float64) error {
	// Insert invoice into the database
	_, err := db.Exec(
		"INSERT INTO invoices (ReservationID, UserID, Amount, PaymentStatus) VALUES (?, ?, ?, 'Pending')",
		reservationID, userID, amount,
	)
	if err != nil {
		return fmt.Errorf("error inserting invoice into database: %v", err)
	}
	return nil
}
