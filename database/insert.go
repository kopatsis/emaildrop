package database

import "database/sql"

func InsertEntry(db *sql.DB, entry Entry) error {
	query := `
		INSERT INTO entries (timestamp, requestID, name, email, subject, comment, iphash, city, country, has_error, error_source, error_message, response)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, entry.Timestamp, entry.RequestID, entry.Name, entry.Email, entry.Subject, entry.Comment, entry.IPHash, entry.City, entry.Country, entry.HasError, entry.ErrorSource, entry.ErrorMessage, entry.Response)
	if err != nil {
		return err
	}

	return nil
}
