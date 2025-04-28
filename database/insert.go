package database

import "database/sql"

func InsertEntry(db *sql.DB, entry Entry) error {
	query := `
		INSERT INTO entries (requestID, timestamp, name, email, subject, comment, iphash, city, country, complete, has_error, error_source, error_message, response)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, entry.RequestID, entry.Timestamp, entry.Name, entry.Email, entry.Subject, entry.Comment, entry.IPHash, entry.City, entry.Country, entry.Complete, entry.HasError, entry.ErrorSource, entry.ErrorMessage, entry.Response)
	return err
}
