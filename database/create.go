package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	ID           int64     `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	RequestID    string    `json:"requestID"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Subject      string    `json:"subject"`
	Comment      string    `json:"comment"`
	Complete     bool      `json:"complete"`
	IPHash       string    `json:"iphash"`
	City         string    `json:"city"`
	Country      string    `json:"country"`
	HasError     bool      `json:"has_error"`
	ErrorSource  string    `json:"error_source"`
	ErrorMessage string    `json:"error_message"`
	Response     string    `json:"response"`
}

func SetupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		requestID TEXT,
		name TEXT,
		email TEXT,
		subject TEXT,
		comment TEXT,
		complete BOOLEAN,
		iphash TEXT,
		city TEXT,
		country TEXT,
		has_error BOOLEAN,
		error_source TEXT,
		error_message TEXT,
		response TEXT
	);`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}
