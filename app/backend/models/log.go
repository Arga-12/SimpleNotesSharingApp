package models

import (
	"database/sql"
	"time"
)

type Log struct {
	ID             int           `json:"id"`
	Datetime       time.Time     `json:"datetime"`
	Method         string        `json:"method"`
	Endpoint       string        `json:"endpoint"`
	RequestHeaders string        `json:"requestHeaders"` // JSON string
	RequestPayload string        `json:"requestPayload"`
	ResponseBody   string        `json:"responseBody"`
	ResponseStatus int           `json:"responseStatus"`
	DurationMs     int           `json:"durationMs"`
	UserID         sql.NullInt64 `json:"userId,omitempty"`
	CreatedAt      time.Time     `json:"createdAt"`
}

// SaveLogToDB saves log entry to database
func SaveLogToDB(db *sql.DB, l *Log) error {
	// Limit sizes to prevent huge logs
	if len(l.ResponseBody) > 10000 {
		l.ResponseBody = l.ResponseBody[:10000] + "... [TRUNCATED]"
	}
	if len(l.RequestPayload) > 10000 {
		l.RequestPayload = l.RequestPayload[:10000] + "... [TRUNCATED]"
	}

	query := `
		INSERT INTO logs (datetime, method, endpoint, request_headers, request_payload, response_body, response_status, duration_ms, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	err := db.QueryRow(
		query,
		l.Datetime,
		l.Method,
		l.Endpoint,
		l.RequestHeaders,
		l.RequestPayload,
		l.ResponseBody,
		l.ResponseStatus,
		l.DurationMs,
		l.UserID,
	).Scan(&l.ID, &l.CreatedAt)

	return err
}
