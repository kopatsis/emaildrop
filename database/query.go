package database

import (
	"database/sql"
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

func GetEntriesInLast24Hours(db *sql.DB) ([]Entry, error) {
	now := time.Now()
	oneDayAgo := now.Add(-24 * time.Hour)

	rows, err := db.Query("SELECT timestamp, requestID, name, email, subject, comment, has_error, error_source, error_message, response FROM entries WHERE timestamp >= ?", oneDayAgo)
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
