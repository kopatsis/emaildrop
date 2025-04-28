package database

import "database/sql"

func InsertEntry(db *sql.DB, entry Entry) (int64, error) {
	query := `
		INSERT INTO entries (timestamp, requestID, name, email, subject, comment, iphash, city, country, complete, has_error, error_source, error_message, response)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := db.Exec(query, entry.Timestamp, entry.RequestID, entry.Name, entry.Email, entry.Subject, entry.Comment, entry.Complete, entry.IPHash, entry.City, entry.Country, entry.HasError, entry.ErrorSource, entry.ErrorMessage, entry.Response)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func UpdateEntry(db *sql.DB, entry Entry) error {
	query := `
		UPDATE entries
		SET timestamp = ?, requestID = ?, name = ?, email = ?, subject = ?, comment = ?, iphash = ?, city = ?, country = ?, complete = ?, has_error = ?, error_source = ?, error_message = ?, response = ?
		WHERE id = ?
	`

	_, err := db.Exec(query, entry.Timestamp, entry.RequestID, entry.Name, entry.Email, entry.Subject, entry.Comment, entry.IPHash, entry.City, entry.Country, entry.Complete, entry.HasError, entry.ErrorSource, entry.ErrorMessage, entry.Response, entry.ID)
	return err
}
