// backend/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yourname/fullstack-auth-backend/handlers"
	"github.com/yourname/fullstack-auth-backend/middlewares"
)

var (
	DBHost = getenv("DB_HOST", "localhost")
	DBPort = getenv("DB_PORT", "5432")
	DBUser = getenv("DB_USER", "pgandul")
	DBPass = getenv("DB_PASSWORD", "postgreroot")
	DBName = getenv("DB_NAME", "authdb")
	Port   = getenv("PORT", "8080")
)

func main() {
	db := InitDB()
	defer db.Close()

	// Initialize handlers
	authHandler := &handlers.AuthHandler{DB: db}
	notesHandler := &handlers.NotesHandler{DB: db}

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.Handle("/api/register", middlewares.Logging(http.HandlerFunc(authHandler.HandleRegister)))
	mux.Handle("/api/login", middlewares.Logging(http.HandlerFunc(authHandler.HandleLogin)))
	mux.Handle("/api/me", middlewares.Logging(http.HandlerFunc(authHandler.HandleMe)))
	mux.Handle("/api/logout", middlewares.Logging(http.HandlerFunc(authHandler.HandleLogout)))

	// Notes routes
	mux.Handle("/api/notes", middlewares.Logging(http.HandlerFunc(notesHandler.HandleNotes)))
	mux.Handle("/api/notes/", middlewares.Logging(http.HandlerFunc(notesHandler.HandleNoteByID)))

	addr := ":" + Port
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, middlewares.AllowLocalhostCookies(mux)); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}