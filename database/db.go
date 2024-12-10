package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Initialize and return database connection
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:#17436592a@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Database ping failed: %v", err)
		return nil, err
	}

	return db, nil
}
