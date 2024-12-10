package userservice

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FirstName      string
	LastName       string
	Email          string
	PhoneNumber    string
	Password       string
	MembershipTier string
}

func RegisterUser(db *sql.DB, user User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	// Insert user data into the database
	_, err = db.Exec(
		"INSERT INTO Users (FirstName, LastName, Email, PhoneNumber, PasswordHash, MembershipTier) VALUES (?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.PhoneNumber, string(hashedPassword), user.MembershipTier,
	)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}

	return nil
}

// AuthenticateUser handles user login logic
func AuthenticateUser(db *sql.DB, email, password string) error {
	// Query user from database
	var hashedPassword string
	err := db.QueryRow("SELECT PasswordHash FROM Users WHERE Email = ?", email).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return errors.New("invalid email or password")
	}

	// Compare hashed password with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Printf("Invalid password")
		return errors.New("invalid email or password")
	}

	return nil
}
