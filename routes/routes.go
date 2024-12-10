package routes

import (
	"encoding/json"
	"net/http"

	"github.com/zed99-git/electric-car-sharing/database"
	"github.com/zed99-git/electric-car-sharing/middleware"
	userservice "github.com/zed99-git/electric-car-sharing/services/user-service"
)

// Handler for user registration
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Decode user registration payload
	var req struct {
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Email          string `json:"email"`
		PhoneNumber    string `json:"phone_number"`
		Password       string `json:"password"`
		MembershipTier string `json:"membership_tier"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Default value for MembershipTier if it's missing
	if req.MembershipTier == "" {
		req.MembershipTier = "Basic"
	}

	user := userservice.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		Password:       req.Password,
		MembershipTier: req.MembershipTier,
	}

	// Register user with database
	err = userservice.RegisterUser(db, user)
	if err != nil {
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}

// Handler for user login/authentication
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = userservice.AuthenticateUser(db, credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// Example of a protected route
func ProtectedRouteHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You have access to this protected route!"))
}

// Register routes with authentication
func RegisterRoutes() {
	http.HandleFunc("/api/v1/protected", middleware.AuthMiddleware(ProtectedRouteHandler))
}
