package middlewares

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Arga-12/SimpleNotesSharingApp/app/backend/models"
)

// Global DB instance for logging
var LogDB *sql.DB

// SetLogDB sets the database instance for logging
func SetLogDB(db *sql.DB) {
	LogDB = db
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture request body
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // Restore body for handlers
		}

		// Capture request headers (mask sensitive data)
		headers := make(map[string]string)
		for name, values := range r.Header {
			if strings.ToLower(name) == "authorization" {
				// Mask Authorization header
				if len(values) > 0 && len(values[0]) > 10 {
					headers[name] = values[0][:10] + "***MASKED***"
				} else {
					headers[name] = "***MASKED***"
				}
			} else if strings.ToLower(name) == "cookie" {
				// Mask cookie values
				headers[name] = "***MASKED***"
			} else {
				headers[name] = strings.Join(values, ", ")
			}
		}
		headersJSON, _ := json.Marshal(headers)

		// Custom response writer to capture response
		lrw := &logResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           &bytes.Buffer{},
		}

		// Get user ID from context if exists
		var userID sql.NullInt64
		if uid, err := GetUserIDFromCookie(r); err == nil {
			userID = sql.NullInt64{Int64: int64(uid), Valid: true}
		}

		// Process request
		next.ServeHTTP(lrw, r)

		// Calculate duration
		duration := time.Since(start)
		durationMs := int(duration.Milliseconds())

		// Log to console
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, duration)

		// Save to database asynchronously
		if LogDB != nil {
			go saveLogToDB(
				r.Method,
				r.URL.Path,
				string(headersJSON),
				string(requestBody),
				lrw.body.String(),
				lrw.statusCode,
				durationMs,
				userID,
			)
		}
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (l *logResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *logResponseWriter) Write(b []byte) (int, error) {
	// Capture response body
	l.body.Write(b)
	return l.ResponseWriter.Write(b)
}

func saveLogToDB(method, endpoint, headers, payload, responseBody string, status, durationMs int, userID sql.NullInt64) {
	if LogDB == nil {
		return
	}

	// Create log entry using model
	logEntry := &models.Log{
		Datetime:       time.Now(),
		Method:         method,
		Endpoint:       endpoint,
		RequestHeaders: headers,
		RequestPayload: payload,
		ResponseBody:   responseBody,
		ResponseStatus: status,
		DurationMs:     durationMs,
		UserID:         userID,
	}

	// Save to database using model method
	err := models.SaveLogToDB(LogDB, logEntry)
	if err != nil {
		log.Printf("❌ Failed to save log to DB: %v", err)
	}
}
