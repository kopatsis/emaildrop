package database

import (
	"database/sql"
	"errors"
	"time"
)

func GetCountOfSuccessfulEntries(db *sql.DB) (int, error) {
	now := time.Now()
	oneDayAgo := now.Add(-24 * time.Hour)

	var count int
	query := `SELECT COUNT(*) FROM entries WHERE timestamp >= ? AND complete = 1`
	err := db.QueryRow(query, oneDayAgo).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetEntriesSince(db *sql.DB, since time.Time) ([]Entry, error) {
	if since.After(time.Now()) || since.Before(time.Date(2025, 4, 27, 0, 0, 0, 0, time.Local)) {
		return nil, errors.New("incorrect date provided: " + since.String())
	}

	rows, err := db.Query("SELECT timestamp, requestID, name, email, subject, comment, has_error, error_source, error_message, response FROM entries WHERE timestamp >= ?", since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		err := rows.Scan(&entry.Timestamp, &entry.RequestID, &entry.Name, &entry.Email, &entry.Subject, &entry.Comment, &entry.HasError, &entry.ErrorSource, &entry.ErrorMessage, &entry.Response)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func GetEntry(db *sql.DB, reqID string) (Entry, error) {

	var entry Entry
	if reqID == "" {
		return entry, errors.New("no param provided")
	}

	err := db.QueryRow("SELECT timestamp, requestID, name, email, subject, comment, has_error, error_source, error_message, response FROM entries WHERE requestID = ?", reqID).
		Scan(&entry.Timestamp, &entry.RequestID, &entry.Name, &entry.Email, &entry.Subject, &entry.Comment, &entry.HasError, &entry.ErrorSource, &entry.ErrorMessage, &entry.Response)
	return entry, err
}
